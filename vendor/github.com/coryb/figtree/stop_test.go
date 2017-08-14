package figtree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsStopConfigD3(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1/d2/d3")
	defer os.Chdir("../../..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"stop.yml", true, "d3arr1val1"})
	arr1 = append(arr1, StringOption{"stop.yml", true, "d3arr1val2"})
	arr1 = append(arr1, StringOption{"../stop.yml", true, "d2arr1val1"})
	arr1 = append(arr1, StringOption{"../stop.yml", true, "d2arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"stop.yml", true, "d3str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key1": StringOption{"../stop.yml", true, "d2map1val1"},
			"key2": StringOption{"stop.yml", true, "d3map1val2"},
			"key3": StringOption{"stop.yml", true, "d3map1val3"},
		},
		Int1:   IntOption{"stop.yml", true, 333},
		Float1: Float32Option{"stop.yml", true, 3.33},
		Bool1:  BoolOption{"stop.yml", true, true},
	}

	err := LoadAllConfigs("stop.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestOptionsStopConfigD2(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1/d2")
	defer os.Chdir("../..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"stop.yml", true, "d2arr1val1"})
	arr1 = append(arr1, StringOption{"stop.yml", true, "d2arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"stop.yml", true, "d2str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key1": StringOption{"stop.yml", true, "d2map1val1"},
			"key2": StringOption{"stop.yml", true, "d2map1val2"},
		},
		Int1:   IntOption{"stop.yml", true, 222},
		Float1: Float32Option{"stop.yml", true, 2.22},
		Bool1:  BoolOption{"stop.yml", true, false},
	}

	err := LoadAllConfigs("stop.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinStopConfigD3(t *testing.T) {
	opts := TestBuiltin{}
	os.Chdir("d1/d2/d3")
	defer os.Chdir("../../..")

	arr1 := []string{}
	arr1 = append(arr1, "d3arr1val1")
	arr1 = append(arr1, "d3arr1val2")
	arr1 = append(arr1, "d2arr1val1")
	arr1 = append(arr1, "d2arr1val2")

	expected := TestBuiltin{
		String1:    "d3str1val1",
		LeaveEmpty: "",
		Array1:     arr1,
		Map1: map[string]string{
			"key1": "d2map1val1",
			"key2": "d3map1val2",
			"key3": "d3map1val3",
		},
		Int1:   333,
		Float1: 3.33,
		Bool1:  true,
	}

	err := LoadAllConfigs("stop.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinStopConfigD2(t *testing.T) {
	opts := TestBuiltin{}
	os.Chdir("d1/d2")
	defer os.Chdir("../..")

	arr1 := []string{}
	arr1 = append(arr1, "d2arr1val1")
	arr1 = append(arr1, "d2arr1val2")

	expected := TestBuiltin{
		String1:    "d2str1val1",
		LeaveEmpty: "",
		Array1:     arr1,
		Map1: map[string]string{
			"key1": "d2map1val1",
			"key2": "d2map1val2",
		},
		Int1:   222,
		Float1: 2.22,
		Bool1:  false,
	}

	err := LoadAllConfigs("stop.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}
