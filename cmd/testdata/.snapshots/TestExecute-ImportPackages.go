// Code generated by accessory; DO NOT EDIT.

package test

import (
	"github.com/masaushi/accessory/cmd/testdata/import_packages/sub1"
	"github.com/masaushi/accessory/cmd/testdata/import_packages/sub2"
	"github.com/masaushi/accessory/cmd/testdata/import_packages/sub3"
	"time"
)

func (t *Tester) Field1() time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.field1
}

func (t *Tester) SetField1(val time.Time) {
	if t == nil {
		return
	}
	t.field1 = val
}

func (t *Tester) Field2() *sub1.SubTester {
	if t == nil {
		return nil
	}
	return t.field2
}

func (t *Tester) SetField2(val *sub1.SubTester) {
	if t == nil {
		return
	}
	t.field2 = val
}

func (t *Tester) Field3() *sub2.SubTester {
	if t == nil {
		return nil
	}
	return t.field3
}

func (t *Tester) SetField3(val *sub2.SubTester) {
	if t == nil {
		return
	}
	t.field3 = val
}

func (t *Tester) Field4() *sub3.SubTester {
	if t == nil {
		return nil
	}
	return t.field4
}

func (t *Tester) SetField4(val *sub3.SubTester) {
	if t == nil {
		return
	}
	t.field4 = val
}

