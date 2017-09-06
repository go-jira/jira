package figtree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsExecConfigD3(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1/d2/d3")
	defer os.Chdir("../../..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"exec.yml", true, "d3arr1val1"})
	arr1 = append(arr1, StringOption{"exec.yml", true, "d3arr1val2"})
	arr1 = append(arr1, StringOption{"../exec.yml", true, "d2arr1val1"})
	arr1 = append(arr1, StringOption{"../exec.yml", true, "d2arr1val2"})
	arr1 = append(arr1, StringOption{"../../exec.yml", true, "d1arr1val1"})
	arr1 = append(arr1, StringOption{"../../exec.yml", true, "d1arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"exec.yml", true, "d3str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key0": StringOption{"../../exec.yml", true, "d1map1val0"},
			"key1": StringOption{"../exec.yml", true, "d2map1val1"},
			"key2": StringOption{"exec.yml", true, "d3map1val2"},
			"key3": StringOption{"exec.yml", true, "d3map1val3"},
		},
		Int1:   IntOption{"exec.yml", true, 333},
		Float1: Float32Option{"exec.yml", true, 3.33},
		Bool1:  BoolOption{"exec.yml", true, true},
	}

	err := LoadAllConfigs("exec.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestOptionsExecConfigD2(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1/d2")
	defer os.Chdir("../..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"exec.yml", true, "d2arr1val1"})
	arr1 = append(arr1, StringOption{"exec.yml", true, "d2arr1val2"})
	arr1 = append(arr1, StringOption{"../exec.yml", true, "d1arr1val1"})
	arr1 = append(arr1, StringOption{"../exec.yml", true, "d1arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"exec.yml", true, "d2str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key0": StringOption{"../exec.yml", true, "d1map1val0"},
			"key1": StringOption{"exec.yml", true, "d2map1val1"},
			"key2": StringOption{"exec.yml", true, "d2map1val2"},
		},
		Int1:   IntOption{"exec.yml", true, 222},
		Float1: Float32Option{"exec.yml", true, 2.22},
		Bool1:  BoolOption{"exec.yml", true, false},
	}

	err := LoadAllConfigs("exec.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestOptionsExecConfigD1(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1")
	defer os.Chdir("..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"exec.yml", true, "d1arr1val1"})
	arr1 = append(arr1, StringOption{"exec.yml", true, "d1arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"exec.yml", true, "d1str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key0": StringOption{"exec.yml", true, "d1map1val0"},
			"key1": StringOption{"exec.yml", true, "d1map1val1"},
		},
		Int1:   IntOption{"exec.yml", true, 111},
		Float1: Float32Option{"exec.yml", true, 1.11},
		Bool1:  BoolOption{"exec.yml", true, true},
	}

	err := LoadAllConfigs("exec.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinExecConfigD3(t *testing.T) {
	opts := TestBuiltin{}
	os.Chdir("d1/d2/d3")
	defer os.Chdir("../../..")

	arr1 := []string{}
	arr1 = append(arr1, "d3arr1val1")
	arr1 = append(arr1, "d3arr1val2")
	arr1 = append(arr1, "d2arr1val1")
	arr1 = append(arr1, "d2arr1val2")
	arr1 = append(arr1, "d1arr1val1")
	arr1 = append(arr1, "d1arr1val2")

	expected := TestBuiltin{
		String1:    "d3str1val1",
		LeaveEmpty: "",
		Array1:     arr1,
		Map1: map[string]string{
			"key0": "d1map1val0",
			"key1": "d2map1val1",
			"key2": "d3map1val2",
			"key3": "d3map1val3",
		},
		Int1:   333,
		Float1: 3.33,
		Bool1:  true,
	}

	err := LoadAllConfigs("exec.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinExecConfigD2(t *testing.T) {
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
			"key1": "d2map1val1",
			"key2": "d2map1val2",
		},
		Int1:   222,
		Float1: 2.22,
		// note this will be true from d1/exec.yml since the
		// d1/d2/exec.yml set it to false which is a zero value
		Bool1: true,
	}

	err := LoadAllConfigs("exec.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinExecConfigD1(t *testing.T) {
	opts := TestBuiltin{}
	os.Chdir("d1")
	defer os.Chdir("..")

	arr1 := []string{}
	arr1 = append(arr1, "d1arr1val1")
	arr1 = append(arr1, "d1arr1val2")

	expected := TestBuiltin{
		String1:    "d1str1val1",
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

	err := LoadAllConfigs("exec.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}
