package math_test

import (
	"testing"

	"github.com/cheekybits/genny/examples/davechaney"
)

func TestNumberTypeMax(t *testing.T) {

	var v math.NumberType
	v = math.NumberTypeMax(10, 20)
	if v != 20 {
		t.Errorf("Max of 10 and 20 is 20")
	}

	v = math.NumberTypeMax(20, 20)
	if v != 20 {
		t.Errorf("Max of 20 and 20 is 20")
	}

	v = math.NumberTypeMax(25, 20)
	if v != 25 {
		t.Errorf("Max of 25 and 20 is 25")
	}

}
