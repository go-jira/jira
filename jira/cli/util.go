package cli

import (
	"os"
	"fmt"
	"errors"
	"strings"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"text/template"
	"io"
	"bufio"
	"bytes"
	"github.com/mgutz/ansi"
)

func FindParentPaths(fileName string) []string {
	cwd, _ := os.Getwd()

	paths := make([]string,0)

	// special case if homedir is not in current path then check there anyway
	homedir := os.Getenv("HOME")
	if ! strings.HasPrefix(cwd, homedir) {
		file := fmt.Sprintf("%s/%s", homedir, fileName)
		if _, err := os.Stat(file); err == nil {
			paths = append(paths, file)
		}
	}

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

func runTemplate(templateContent string, data interface{}, out io.Writer) error {

	if out == nil {
		out = os.Stdout
	}

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
		"split": func(sep string, content string) []string {
			return strings.Split(content, sep)
		},
	}
	if tmpl, err := template.New("template").Funcs(funcs).Parse(templateContent); err != nil {
		log.Error("Failed to parse template: %s", err)
		return err
	} else {
		if err := tmpl.Execute(out, data); err != nil {
			log.Error("Failed to execute template: %s", err)
			return err
		}
	}
	return nil
}

func responseToJson(resp *http.Response, err error) (interface{}, error) {
	if err != nil {
		return nil, err
	} else {
		return jsonDecode(resp.Body), nil
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

func jsonEncode(data interface{}) (string, error) {
	buffer := bytes.NewBuffer(make([]byte, 0))
	enc := json.NewEncoder(buffer)
	
	err := enc.Encode(data); if err != nil {
		log.Error("Failed to encode data %s: %s", data, err)
		return "", err
	}
	return buffer.String(), nil
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

func promptYN(prompt string, yes bool) bool {
	reader := bufio.NewReader(os.Stdin)
	if !yes {
		prompt = fmt.Sprintf("%s [y/N]: ", prompt)
	} else {
		prompt = fmt.Sprintf("%s [Y/n]: ", prompt)
	}
	
	fmt.Printf("%s", prompt)
	text, _ := reader.ReadString('\n')
	ans := strings.ToLower(strings.TrimRight(text, "\n"))
	if ans == "" {
		return yes
	}
	if( strings.HasPrefix(ans, "y") ) {
		return true
	}
	return false
}

func yamlFixup( data interface{} ) (interface{}, error) {
	switch d := data.(type) {
	case map[interface{}]interface{}:
		// need to copy this map into a string map so json can encode it
		copy := make(map[string]interface{})
		for key, val := range d {
			switch k := key.(type) {
			case string:
				if fixed, err := yamlFixup(val); err != nil {
					return nil, err
				} else {
					copy[k] = fixed
				}
			default:
				err := fmt.Errorf("YAML: key %s is type '%T', require 'string'", key, k)
				log.Error("%s", err)
				return nil, err
			}
		}
		return copy, nil
	case map[string]interface{}: 
		for k, v := range d {
			if fixed, err := yamlFixup(v); err != nil {
				return nil, err
			} else {
				d[k] = fixed
			}
		}
		return d, nil
	case []interface{}:
		for i, val := range d {
			if fixed, err := yamlFixup(val); err != nil {
				return nil, err
			} else {
				d[i] = fixed
			}
		}
		return data, nil
	default: return d, nil
	}
}

