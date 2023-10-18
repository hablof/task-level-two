package pattern

import "testing"

func Test_builder(t *testing.T) {
	b := builder{}
	res := b.newProduct().
		addParam("param1").
		setEnum(firstOption).
		setStrField("strField").
		addParam("param2").
		setIntValue(42).
		result()

	ref := product{
		strField: "strField",
		params:   []string{"param1", "param2"},
		intVal:   42,
		oneOf:    firstOption,
	}
	if res.intVal != ref.intVal ||
		res.oneOf != ref.oneOf ||
		res.strField != ref.strField {
		t.Fail()
		return
	}

	if len(res.params) != len(ref.params) {
		t.Fail()
		return
	}

	for i, elem := range ref.params {
		if res.params[i] != elem {
			t.Fail()
			return
		}
	}
}
