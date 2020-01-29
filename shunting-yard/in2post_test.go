package main

import (
	"testing"
)

func byte_array_equal(a, b []byte) bool {

	if len(a) != len(b) {
		return false
	}

	for i, a_e := range a {
		if b[i] != a_e {
			return false
		}
	}
	return true
}

func TestInfix2Post(t *testing.T) {

	output := In2Post("s+.a+")
	expect := []byte("s+a+.")
	if !byte_array_equal(output, expect) {
		t.Errorf("In2Post expected:%s, got:%s", expect, output)
	}

	output = In2Post("s+.a+.b+")
	expect = []byte("s+a+.b+.")
	if !byte_array_equal(output, expect) {
		t.Errorf("In2Post expected:%s, got:%s", expect, output)
	}

}
