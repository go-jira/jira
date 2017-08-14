package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsAlphaNumeric(t *testing.T) {

	for _, r := range []rune{'a', '1', '_', 'A', 'Z'} {
		assert.True(t, isAlphaNumeric(r))
	}
	for _, r := range []rune{' ', '[', ']', '!', '"'} {
		assert.False(t, isAlphaNumeric(r))
	}

}

func TestWordify(t *testing.T) {

	for word, wordified := range map[string]string{
		"int":         "Int",
		"*int":        "Int",
		"string":      "String",
		"*MyType":     "MyType",
		"*myType":     "MyType",
		"interface{}": "Interface",
		"pack.type":   "Packtype",
		"*pack.type":  "Packtype",
	} {
		assert.Equal(t, wordified, wordify(word, true))
	}

}
