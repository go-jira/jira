package figtree

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	StringifyValue = false
}

type TestOptions struct {
	String1    StringOption     `json:"str1,omitempty" yaml:"str1,omitempty"`
	LeaveEmpty StringOption     `json:"leave-empty,omitempty" yaml:"leave-empty,omitempty"`
	Array1     ListStringOption `json:"arr1,omitempty" yaml:"arr1,omitempty"`
	Map1       MapStringOption  `json:"map1,omitempty" yaml:"map1,omitempty"`
	Int1       IntOption        `json:"int1,omitempty" yaml:"int1,omitempty"`
	Float1     Float32Option    `json:"float1,omitempty" yaml:"float1,omitempty"`
	Bool1      BoolOption       `json:"bool1,omitempty" yaml:"bool1,omitempty"`
}

type TestBuiltin struct {
	String1    string            `yaml:"str1,omitempty"`
	LeaveEmpty string            `yaml:"leave-empty,omitempty"`
	Array1     []string          `yaml:"arr1,omitempty"`
	Map1       map[string]string `yaml:"map1,omitempty"`
	Int1       int               `yaml:"int1,omitempty"`
	Float1     float32           `yaml:"float1,omitempty"`
	Bool1      bool              `yaml:"bool1,omitempty"`
}

func TestOptionsLoadConfigD3(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1/d2/d3")
	defer os.Chdir("../../..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"figtree.yml", true, "d3arr1val1"})
	arr1 = append(arr1, StringOption{"figtree.yml", true, "d3arr1val2"})
	arr1 = append(arr1, StringOption{"../figtree.yml", true, "d2arr1val1"})
	arr1 = append(arr1, StringOption{"../figtree.yml", true, "d2arr1val2"})
	arr1 = append(arr1, StringOption{"../../figtree.yml", true, "d1arr1val1"})
	arr1 = append(arr1, StringOption{"../../figtree.yml", true, "d1arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"figtree.yml", true, "d3str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key0": StringOption{"../../figtree.yml", true, "d1map1val0"},
			"key1": StringOption{"../figtree.yml", true, "d2map1val1"},
			"key2": StringOption{"figtree.yml", true, "d3map1val2"},
			"key3": StringOption{"figtree.yml", true, "d3map1val3"},
		},
		Int1:   IntOption{"figtree.yml", true, 333},
		Float1: Float32Option{"figtree.yml", true, 3.33},
		Bool1:  BoolOption{"figtree.yml", true, true},
	}

	err := LoadAllConfigs("figtree.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestOptionsLoadConfigD2(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1/d2")
	defer os.Chdir("../..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"figtree.yml", true, "d2arr1val1"})
	arr1 = append(arr1, StringOption{"figtree.yml", true, "d2arr1val2"})
	arr1 = append(arr1, StringOption{"../figtree.yml", true, "d1arr1val1"})
	arr1 = append(arr1, StringOption{"../figtree.yml", true, "d1arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"figtree.yml", true, "d2str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key0": StringOption{"../figtree.yml", true, "d1map1val0"},
			"key1": StringOption{"figtree.yml", true, "d2map1val1"},
			"key2": StringOption{"figtree.yml", true, "d2map1val2"},
		},
		Int1:   IntOption{"figtree.yml", true, 222},
		Float1: Float32Option{"figtree.yml", true, 2.22},
		Bool1:  BoolOption{"figtree.yml", true, false},
	}

	err := LoadAllConfigs("figtree.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestOptionsLoadConfigD1(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1")
	defer os.Chdir("..")

	arr1 := []StringOption{}
	arr1 = append(arr1, StringOption{"figtree.yml", true, "d1arr1val1"})
	arr1 = append(arr1, StringOption{"figtree.yml", true, "d1arr1val2"})

	expected := TestOptions{
		String1:    StringOption{"figtree.yml", true, "d1str1val1"},
		LeaveEmpty: StringOption{},
		Array1:     arr1,
		Map1: map[string]StringOption{
			"key0": StringOption{"figtree.yml", true, "d1map1val0"},
			"key1": StringOption{"figtree.yml", true, "d1map1val1"},
		},
		Int1:   IntOption{"figtree.yml", true, 111},
		Float1: Float32Option{"figtree.yml", true, 1.11},
		Bool1:  BoolOption{"figtree.yml", true, true},
	}

	err := LoadAllConfigs("figtree.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinLoadConfigD3(t *testing.T) {
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

	err := LoadAllConfigs("figtree.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinLoadConfigD2(t *testing.T) {
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
		// note this will be true from d1/figtree.yml since the
		// d1/d2/figtree.yml set it to false which is a zero value
		Bool1: true,
	}

	err := LoadAllConfigs("figtree.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}

func TestBuiltinLoadConfigD1(t *testing.T) {
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

	err := LoadAllConfigs("figtree.yml", &opts)
	assert.Nil(t, err)
	assert.Exactly(t, expected, opts)
}
