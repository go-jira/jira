package slice

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnsureSlice(t *testing.T) {

	myType := new(MyType)
	slice := EnsureSlice(myType)
	if assert.NotNil(t, slice) {
		assert.Equal(t, slice[0], myType)
	}

	slice = EnsureSlice(slice)
	log.Printf("%#v", slice[0])
	if assert.NotNil(t, slice) {
		assert.Equal(t, slice[0], myType)
	}

}
