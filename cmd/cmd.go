package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"

	"github.com/masaushi/accessory/internal/generator"
	"github.com/masaushi/accessory/internal/parser"
)

const Version = "dev"

// usage is a replacement for usage function in the flag package.
func usage() {
	fmt.Fprintf(os.Stderr, "Usage of accessory:\n")
	fmt.Fprintf(os.Stderr, "\taccessory [flags] [directory]\n")
	fmt.Fprintf(os.Stderr, "For more information, see:\n")
	fmt.Fprintf(os.Stderr, "\thttps://github.com/masaushi/accessory\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

// Execute executes a whole process of generating accessor codes.
func Execute(fs afero.Fs, args []string) {
	log.SetFlags(0)
	log.SetPrefix("accessory: ")

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.Usage = usage
	version := flags.Bool("version", false, "show the version of accessory")
	typeName := flags.String("type", "", "type name; must be set")
	receiver := flags.String("receiver", "", "receiver name; default first letter of type name")
	output := flags.String("output", "", "output file name; default <type_name>_accessor.go")

	if err := flags.Parse(args[1:]); err != nil {
		flags.Usage()
		os.Exit(1)
	}

	if *version {
		fmt.Fprintf(os.Stdout, "accessory version: %s\n", Version)
		os.Exit(0)
	}

	if typeName == nil || len(*typeName) == 0 {
		flags.Usage()
		os.Exit(1)
	}

	var dir string
	if cliArgs := flags.Args(); len(cliArgs) > 0 {
		dir = cliArgs[0]
	} else {
		// Default: process whole package in current directory.
		dir = "."
	}

	if !isDir(dir) {
		fmt.Fprintln(os.Stderr, "Specified argument is not a directory.")
		flags.Usage()
		os.Exit(1)
	}

	pkg, err := parser.ParsePackage(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		flags.Usage()
		os.Exit(1)
	}

	if err = generator.Generate(fs, pkg, *typeName, *output, *receiver); err != nil {
		log.Fatal(err)
	}
}

func isDir(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}
