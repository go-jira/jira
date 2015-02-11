package cli

import (
	"os"
	"fmt"
	"errors"
	"strings"

	"encoding/json"
	"io/ioutil"
	"text/template"
	"io"
	"github.com/mgutz/ansi"
)

func FindParentPaths(fileName string) []string {
	cwd, _ := os.Getwd()

	paths := make([]string,0)

	var dir string
	for _, part := range strings.Split(cwd, string(os.PathSeparator)) {
		if dir == "/" {
			dir = fmt.Sprintf("/%s", part)
		} else {
			dir = fmt.Sprintf("%s/%s", dir, part)
		}
		file := fmt.Sprintf("%s/%s", dir, fileName)
		if _, err := os.Stat(file); err == nil {
			paths = append(paths, file)
		}
	}
	return paths
}

func FindClosestParentPath(fileName string) (string,error) {
	paths := FindParentPaths(fileName)
	if len(paths) > 0 { 
		return paths[len(paths)-1], nil
	}
	return "", errors.New(fmt.Sprintf("%s not found in parent directory hierarchy", fileName))
}

func readFile(file string) string {
	var bytes []byte
	var err error
	if bytes, err = ioutil.ReadFile(file); err != nil {
		log.Error("Failed to read file %s: %s", file, err)
		os.Exit(1)
	}
	return string(bytes)
}

func runTemplate(text string, data interface{}) {
	funcs := map[string]interface{}{
		"toJson": func(content interface{}) (string, error) {
			if bytes, err := json.MarshalIndent(content, "", "    "); err != nil {
				return "", err
			} else {
				return string(bytes), nil
			}
		},
		"append": func(more string, content interface{}) (string, error) {
			switch value := content.(type) {
			case string: return string(append([]byte(content.(string)), []byte(more)...)), nil
			case []byte: return string(append(content.([]byte), []byte(more)...)), nil
			default: return "", errors.New(fmt.Sprintf("Unknown type: %s", value))
			}
		},
		"indent": func(spaces int, content string) string {
			indent  := make([]byte, spaces + 1, spaces +1)
			indent[0] = '\n'
			for i := 1; i < spaces + 1; i += 1 {
				indent[i] = ' '
			}
			return strings.Replace(content, "\n", string(indent), -1)
		},
		"color": func(color string) string {
			return ansi.ColorCode(color)
		},
	}
	if tmpl, err := template.New("template").Funcs(funcs).Parse(text); err != nil {
		log.Error("Failed to parse template: %s", err)
		os.Exit(1)
	} else {
		if err := tmpl.Execute(os.Stdout, data); err != nil {
			log.Error("Failed to execute template: %s", err)
			os.Exit(1)
		}
	}
}

func jsonDecode(io io.Reader) interface{} {
	content, err := ioutil.ReadAll(io)
	var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		log.Error("JSON Parse Error: %s from %s", err, content)
	}
	return data
}

func jsonWrite(file string, data interface{}) {
	fh, err := os.OpenFile(file, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
	defer fh.Close()
	if err != nil {
		log.Error("Failed to open %s: %s", file, err)
		os.Exit(1)
	}
	enc := json.NewEncoder(fh)
	enc.Encode(data)
}

