package kingpeon

import (
	"io/ioutil"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

func TestRegisterDynamicCommands(t *testing.T) {
	data := struct {
		DynamicCommands []DynamicCommand `yaml:"dynamic-commands"`
	}{}

	config, err := ioutil.ReadFile("./sample.yml")
	assert.Nil(t, err)

	err = yaml.Unmarshal(config, &data)
	assert.Nil(t, err)

	tmpl := template.New("test")
	app := kingpin.New("kingpeon", "Testing Aliases")

	var expectedShell string
	run := func(bin string, cmd []string, env []string) error {
		assert.Equal(t, "/bin/sh", bin)
		assert.Equal(t, []string{"sh", "-c", expectedShell}, cmd)
		assert.NotEmpty(t, env)
		return nil
	}

	err = RegisterDynamicCommandsWithRunner(run, app, data.DynamicCommands, tmpl)
	assert.Nil(t, err)

	expectedShell = "echo hello world"
	_, err = app.Parse([]string{"echo"})
	assert.Nil(t, err)

	expectedShell = "echo -n hello world"
	_, err = app.Parse([]string{"echo", "--no-newline"})
	assert.Nil(t, err)

	expectedShell = "echo hello test"
	_, err = app.Parse([]string{"echo", "test", "--newline"})
	assert.Nil(t, err)

	expectedShell = "echo -n hello test"
	_, err = app.Parse([]string{"echo", "test", "--no-newline"})
	assert.Nil(t, err)

	expectedShell = "echo true"
	_, err = app.Parse([]string{"test", "bool", "arg", "true"})
	assert.Nil(t, err)

	expectedShell = "echo true"
	_, err = app.Parse([]string{"test", "bool", "opt", "--BOOL"})
	assert.Nil(t, err)

	expectedShell = "echo 2"
	_, err = app.Parse([]string{"test", "counter", "arg", "foo", "bar"})
	assert.Nil(t, err)

	expectedShell = "echo 2"
	_, err = app.Parse([]string{"test", "counter", "opt", "--COUNTER", "--COUNTER"})
	assert.Nil(t, err)

	expectedShell = "echo foo"
	_, err = app.Parse([]string{"test", "enum", "arg", "foo"})
	assert.Nil(t, err)

	expectedShell = "echo foo"
	_, err = app.Parse([]string{"test", "enum", "opt", "--ENUM", "foo"})
	assert.Nil(t, err)

	_, err = app.Parse([]string{"test", "enum", "opt", "--ENUM", "bogus"})
	assert.EqualError(t, err, "enum value must be one of foo,bar, got 'bogus'")

	expectedShell = "echo 1.23"
	_, err = app.Parse([]string{"test", "float32", "arg", "1.23"})
	assert.Nil(t, err)

	expectedShell = "echo 1.23"
	_, err = app.Parse([]string{"test", "float32", "opt", "--FLOAT32", "1.23"})
	assert.Nil(t, err)

	expectedShell = "echo 1.23"
	_, err = app.Parse([]string{"test", "float64", "arg", "1.23"})
	assert.Nil(t, err)

	expectedShell = "echo 1.23"
	_, err = app.Parse([]string{"test", "float64", "opt", "--FLOAT64", "1.23"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int8", "arg", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int8", "opt", "--INT8", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int8", "arg", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int8", "opt", "--INT8", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int16", "arg", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int16", "opt", "--INT16", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int32", "arg", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int32", "opt", "--INT32", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int64", "arg", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int64", "opt", "--INT64", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int", "arg", "127"})
	assert.Nil(t, err)

	expectedShell = "echo 127"
	_, err = app.Parse([]string{"test", "int", "opt", "--INT", "127"})
	assert.Nil(t, err)

	expectedShell = "echo hello"
	_, err = app.Parse([]string{"test", "string", "arg", "hello"})
	assert.Nil(t, err)

	expectedShell = "echo hello"
	_, err = app.Parse([]string{"test", "string", "opt", "--STRING", "hello"})
	assert.Nil(t, err)

	expectedShell = "echo [abc: def][foo: bar]"
	_, err = app.Parse([]string{"test", "stringmap", "arg", "foo=bar", "abc=def"})
	assert.Nil(t, err)

	expectedShell = "echo [abc: def][foo: bar]"
	_, err = app.Parse([]string{"test", "stringmap", "opt", "--STRINGMAP", "foo=bar", "--STRINGMAP", "abc=def"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint8", "arg", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint8", "opt", "--UINT8", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint8", "arg", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint8", "opt", "--UINT8", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint16", "arg", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint16", "opt", "--UINT16", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint32", "arg", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint32", "opt", "--UINT32", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint64", "arg", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint64", "opt", "--UINT64", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint", "arg", "255"})
	assert.Nil(t, err)

	expectedShell = "echo 255"
	_, err = app.Parse([]string{"test", "uint", "opt", "--UINT", "255"})
	assert.Nil(t, err)
}
