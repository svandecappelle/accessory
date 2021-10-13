package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/afero"

	"github.com/masaushi/accessory/internal/types"
)

type generator struct {
	buf *bytes.Buffer
}

// Generate generates a file and accessor methods in it.
func Generate(fs afero.Fs, pkg *types.Package, typeName, output, receiverName string) error {
	g := generator{buf: new(bytes.Buffer)}

	setterGen := g.setterGenerator(receiverName)
	getterGen := g.getterGenerator(receiverName)

	accessors := make([]string, 0)
	imports := make(map[string]string)

	for _, file := range pkg.Files {
		for _, st := range file.Structs {
			if st.Name != typeName {
				continue
			}

			for _, field := range st.Fields {
				if field.Tag == nil {
					continue
				}

				typePkg := strings.Split(strings.TrimPrefix(field.DataType, "*"), ".")[0]
				if _, ok := imports[typePkg]; !ok {
					for _, imp := range file.Imports {
						if imp.Name == typePkg {
							imports[imp.Name] = imp.PkgPath
							break
						}
					}
				}

				if field.Tag.Getter != nil {
					getter, err := getterGen(st.Name, field)
					if err != nil {
						return err
					}
					accessors = append(accessors, getter)
				}
				if field.Tag.Setter != nil {
					setter, err := setterGen(st.Name, field)
					if err != nil {
						return err
					}
					accessors = append(accessors, setter)
				}
			}
		}
	}

	g.write(pkg.Name, imports, accessors)

	content, err := g.format()
	if err != nil {
		return err
	}

	outputFile := g.outputFile(output, typeName, pkg.Dir)
	return afero.WriteFile(fs, outputFile, content, 0644)
}

func (g *generator) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

func (g *generator) write(pkgName string, importMap map[string]string, accessors []string) {
	g.printf("// Code generated by accessory; DO NOT EDIT.\n")
	g.printf("\n")
	g.printf("package %s\n", pkgName)
	g.printf("\n")

	if len(importMap) > 0 {
		// Ensure imports are same order as previous if there are no declaration changes.
		importNames := make([]string, 0, len(importMap))
		for name := range importMap {
			importNames = append(importNames, name)
		}
		sort.Strings(importNames)

		g.printf("import (\n")
		for _, name := range importNames {
			path := importMap[name]
			if name == filepath.Base(path) {
				g.printf("\t\"%s\"\n", path)
			} else {
				g.printf("\t%s \"%s\"\n", name, path)
			}
		}
		g.printf(")\n")
	}

	for i := range accessors {
		g.printf("%s\n", accessors[i])
	}
}

func (g *generator) setterGenerator(
	receiverName string,
) func(structName string, field *types.Field) (string, error) {
	const tpl = `
func ({{.Receiver}} *{{.Struct}}) {{.MethodName}}(val {{.Type}}) {
	{{.Receiver}}.{{.Field}} = val
}`
	t := template.Must(template.New("setter").Parse(tpl))

	return func(structName string, field *types.Field) (string, error) {
		methodName := *field.Tag.Setter
		if methodName == "" {
			methodName = fmt.Sprintf("Set%s", strings.Title(field.Name))
		}

		buf := new(bytes.Buffer)
		err := t.Execute(buf, map[string]string{
			"Receiver":   g.receiverName(receiverName, structName),
			"Struct":     structName,
			"MethodName": methodName,
			"Field":      field.Name,
			"Type":       field.DataType,
		})
		if err != nil {
			return "", err
		}

		return buf.String(), nil
	}
}

func (g *generator) getterGenerator(
	receiverName string,
) func(structName string, field *types.Field) (string, error) {
	const tpl = `
func ({{.Receiver}} *{{.Struct}}) {{.MethodName}}() {{.Type}} {
	return {{.Receiver}}.{{.Field}}
}`
	t := template.Must(template.New("getter").Parse(tpl))

	return func(structName string, field *types.Field) (string, error) {
		methodName := *field.Tag.Getter
		if methodName == "" {
			methodName = strings.Title(field.Name)
		}

		buf := new(bytes.Buffer)
		err := t.Execute(buf, map[string]string{
			"Receiver":   g.receiverName(receiverName, structName),
			"Struct":     structName,
			"MethodName": methodName,
			"Field":      field.Name,
			"Type":       field.DataType,
		})
		if err != nil {
			return "", err
		}

		return buf.String(), nil
	}
}

func (g *generator) receiverName(userInput string, structName string) string {
	if userInput != "" {
		// Do nothing if receiver name specified in args.
		return userInput
	}

	// Use the first letter of struct as receiver if receiver name is not specified.
	return strings.ToLower(string(structName[0]))
}

func (g *generator) outputFile(output, typeName, dir string) string {
	if output == "" {
		// Use snake_case name of type as output file if output file is not specified.
		// type TestStruct will be test_struct_accessor.go
		var firstCapMatcher = regexp.MustCompile("(.)([A-Z][a-z]+)")
		var articleCapMatcher = regexp.MustCompile("([a-z0-9])([A-Z])")

		name := firstCapMatcher.ReplaceAllString(typeName, "${1}_${2}")
		name = articleCapMatcher.ReplaceAllString(name, "${1}_${2}")
		output = strings.ToLower(fmt.Sprintf("%s_accessor.go", name))
	}

	return filepath.Join(dir, output)
}

func (g *generator) format() ([]byte, error) {
	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		return g.buf.Bytes(), err
	}
	return src, nil
}
