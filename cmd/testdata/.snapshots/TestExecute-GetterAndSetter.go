// Code generated by accessory; DO NOT EDIT.

package test

func (t *Tester) Field1() string {
	if t == nil {
		return ""
	}
	return t.field1
}

func (t *Tester) SetField1(val string) {
	if t == nil {
		return
	}
	t.field1 = val
}

func (t *Tester) GetSecondField() int32 {
	if t == nil {
		return 0
	}
	return t.field2
}

func (t *Tester) SetSecondField(val int32) {
	if t == nil {
		return
	}
	t.field2 = val
}

