package compare

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type TT1 struct {
	A1 string
	B1 string
	C1 string
	D1 string
}

func t1() TT1 {
	return TT1{
		A1: "A",
		B1: "B",
		C1: "C",
		D1: "D",
	}
}

type TT2 struct {
	A2 *string
	B2 *string
	C2 string
	D2 string
}

func t2() TT2 {
	return TT2{
		A2: strPtr("A"),
		B2: strPtr("B"),
		C2: "C",
		D2: "D",
	}
}

func strPtr(s string)*string{
	return &s
}

func testTranslateFunc(t1 TT1) TT2 {
	return t2()
}

func TestStruct(t *testing.T) {
	var tests = []struct {
		a        interface{}
		b        interface{}
		mappings []CustomMapping
		ignore   []string
	}{
		{t1(), t1(), nil, nil},
		{t2(), t2(), nil, nil},
		{t1(), t2(), nil, nil},
	}
	for _, tt := range tests {
		assert.NoError(t, Struct(tt.a, tt.b, tt.mappings, tt.ignore))
	}
}
