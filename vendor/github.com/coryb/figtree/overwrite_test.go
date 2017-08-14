package figtree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	StringifyValue = false
}

func TestOptionsOverwriteConfigD3(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1/d2/d3")
	defer os.Chdir("../../..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"../overwrite.yml", true, "d2arr1val1"})
	arr1 = append(arr1, StringOption{"../overwrite.yml", true, "d2arr1val2"})
	arr1 = append(arr1, StringOption{"../../overwrite.yml", true, "d1arr1val1"})
	arr1 = append(arr1, StringOption{"../../overwrite.yml", true, "d1arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"../overwrite.yml", true, "d2str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key0": StringOption{"../../overwrite.yml", true, "d1map1val0"},
			"key1": StringOption{"../../overwrite.yml", true, "d1map1val1"},
		},
		Int1:   IntOption{"../../overwrite.yml", true, 111},
		Float1: Float32Option{"../../overwrite.yml", true, 1.11},
		Bool1:  BoolOption{"../overwrite.yml", true, false},
	}

	err := LoadAllConfigs("overwrite.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestOptionsOverwriteConfigD2(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1/d2")
	defer os.Chdir("../..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"overwrite.yml", true, "d2arr1val1"})
	arr1 = append(arr1, StringOption{"overwrite.yml", true, "d2arr1val2"})
	arr1 = append(arr1, StringOption{"../overwrite.yml", true, "d1arr1val1"})
	arr1 = append(arr1, StringOption{"../overwrite.yml", true, "d1arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"overwrite.yml", true, "d2str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key0": StringOption{"../overwrite.yml", true, "d1map1val0"},
			"key1": StringOption{"../overwrite.yml", true, "d1map1val1"},
		},
		Int1:   IntOption{"../overwrite.yml", true, 111},
		Float1: Float32Option{"../overwrite.yml", true, 1.11},
		Bool1:  BoolOption{"overwrite.yml", true, false},
	}

	err := LoadAllConfigs("overwrite.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinOverwriteConfigD3(t *testing.T) {
	opts := TestBuiltin{}
	os.Chdir("d1/d2/d3")
	defer os.Chdir("../../..")

	arr1 := []string{}
	arr1 = append(arr1, "d2arr1val1")
	arr1 = append(arr1, "d2arr1val2")
	arr1 = append(arr1, "d1arr1val1")
	arr1 = append(arr1, "d1arr1val2")

	expected := TestBuiltin{
		String1:    "d2str1val1",
		LeaveEmpty: "",
		Array1:     arr1,
		Map1: map[string]string{
			"key0": "d1map1val0",
			"key1": "d1map1val1",
		},
		Int1:   111,
		Float1: 1.11,
		Bool1:  true,
	}

	err := LoadAllConfigs("overwrite.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinOverwriteConfigD2(t *testing.T) {
	opts := TestBuiltin{}
	os.Chdir("d1/d2")
	defer os.Chdir("../..")

	arr1 := []string{}
	arr1 = append(arr1, "d2arr1val1")
	arr1 = append(arr1, "d2arr1val2")
	arr1 = append(arr1, "d1arr1val1")
	arr1 = append(arr1, "d1arr1val2")

	expected := TestBuiltin{
		String1:    "d2str1val1",
		LeaveEmpty: "",
		Array1:     arr1,
		Map1: map[string]string{
			"key0": "d1map1val0",
			"key1": "d1map1val1",
		},
		Int1:   111,
		Float1: 1.11,
		// note this will be true from d1/overwrite.yml since the
		// d1/d2/overwrite.yml set it to false which is a zero value
		Bool1: true,
	}

	err := LoadAllConfigs("overwrite.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}
