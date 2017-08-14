package figtree

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsEnv(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1")
	defer os.Chdir("..")

	StringifyValue = true
	defer func() {
		StringifyValue = false
	}()

	os.Clearenv()
	err := LoadAllConfigs("figtree.yml", &opts)
	assert.Nil(t, err)

	got := []string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FIGTREE_") {
			got = append(got, env)
		}
	}

	sort.StringSlice(got).Sort()

	expected := []string{
		"FIGTREE_ARRAY_1=[\"d1arr1val1\",\"d1arr1val2\"]",
		"FIGTREE_BOOL_1=true",
		"FIGTREE_FLOAT_1=1.11",
		"FIGTREE_INT_1=111",
		"FIGTREE_MAP_1={\"key0\":\"d1map1val0\",\"key1\":\"d1map1val1\"}",
		"FIGTREE_STRING_1=d1str1val1",
	}

	assert.Equal(t, expected, got)
}

func TestOptionsNamedEnv(t *testing.T) {
	opts := TestOptions{}
	os.Chdir("d1")
	defer os.Chdir("..")

	StringifyValue = true
	defer func() {
		StringifyValue = false
	}()

	os.Clearenv()
	fig := NewFigTree()
	fig.EnvPrefix = "TEST"
	err := fig.LoadAllConfigs("figtree.yml", &opts)
	assert.Nil(t, err)

	got := []string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FIGTREE_") || strings.HasPrefix(env, "TEST_") {
			got = append(got, env)
		}
	}

	sort.StringSlice(got).Sort()

	expected := []string{
		"TEST_ARRAY_1=[\"d1arr1val1\",\"d1arr1val2\"]",
		"TEST_BOOL_1=true",
		"TEST_FLOAT_1=1.11",
		"TEST_INT_1=111",
		"TEST_MAP_1={\"key0\":\"d1map1val0\",\"key1\":\"d1map1val1\"}",
		"TEST_STRING_1=d1str1val1",
	}

	assert.Equal(t, expected, got)
}

func TestBuiltinEnv(t *testing.T) {
	opts := TestBuiltin{}
	os.Chdir("d1")
	defer os.Chdir("..")

	os.Clearenv()
	err := LoadAllConfigs("figtree.yml", &opts)
	assert.Nil(t, err)

	got := []string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FIGTREE_") {
			got = append(got, env)
		}
	}

	sort.StringSlice(got).Sort()

	expected := []string{
		"FIGTREE_ARRAY_1=[\"d1arr1val1\",\"d1arr1val2\"]",
		"FIGTREE_BOOL_1=true",
		"FIGTREE_FLOAT_1=1.11",
		"FIGTREE_INT_1=111",
		"FIGTREE_LEAVE_EMPTY=",
		"FIGTREE_MAP_1={\"key0\":\"d1map1val0\",\"key1\":\"d1map1val1\"}",
		"FIGTREE_STRING_1=d1str1val1",
	}

	assert.Equal(t, expected, got)
}
