package jiracli

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/coryb/figtree"
)

func homedir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func findClosestParentPath(fileName string) (string, error) {
	paths := figtree.FindParentPaths(fileName)
	if len(paths) > 0 {
		return paths[len(paths)-1], nil
	}
	return "", errors.New(fmt.Sprintf("%s not found in parent directory hierarchy", fileName))
}

func (jc *JiraCli) tmpYml(tmpFilePrefix string) (*os.File, error) {
	tmpdir := filepath.Join(homedir(), jc.ConfigDir, "tmp")
	if err := os.MkdirAll(tmpdir, 0755); err != nil {
		return nil, err
	}

	fh, err := ioutil.TempFile(tmpdir, tmpFilePrefix)
	if err != nil {
		return nil, err
	}
	// now we need to rename the file since we dont control the file extensions
	// ... it has to be `.yml` so that vim/emacs etc know what edit mode to apply
	// for easier editing
	oldFileName := fh.Name()
	newFileName := oldFileName + ".yml"

	// close tmpfile so we can rename on windows
	fh.Close()

	if err := os.Rename(oldFileName, newFileName); err != nil {
		return nil, err
	}

	return os.OpenFile(newFileName, os.O_RDWR|os.O_EXCL, 0600)
}

func flagValue(ctx *kingpin.ParseContext, name string) string {
	for _, elem := range ctx.Elements {
		if flag, ok := elem.Clause.(*kingpin.FlagClause); ok {
			if flag.Model().Name == name {
				return *elem.Value
			}
		}
	}
	return ""
}

// func readFile(file string) string {
// 	var bytes []byte
// 	var err error
// 	log.Debugf("readFile: reading %q", file)
// 	if bytes, err = ioutil.ReadFile(file); err != nil {
// 		log.Errorf("Failed to read file %s: %s", file, err)
// 		os.Exit(1)
// 	}
// 	return string(bytes)
// }

func copyFile(src, dst string) (err error) {
	var s, d *os.File
	if s, err = os.Open(src); err == nil {
		defer s.Close()
		if d, err = os.Create(dst); err == nil {
			if _, err = io.Copy(d, s); err != nil {
				d.Close()
				return
			}
			return d.Close()
		}
	}
	return
}

func fuzzyAge(start string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05.000-0700", start)
	if err != nil {
		return "", err
	}
	delta := time.Now().Sub(t)
	if delta.Minutes() < 2 {
		return "a minute", nil
	} else if dm := delta.Minutes(); dm < 45 {
		return fmt.Sprintf("%d minutes", int(dm)), nil
	} else if dm := delta.Minutes(); dm < 90 {
		return "an hour", nil
	} else if dh := delta.Hours(); dh < 24 {
		return fmt.Sprintf("%d hours", int(dh)), nil
	} else if dh := delta.Hours(); dh < 48 {
		return "a day", nil
	}
	return fmt.Sprintf("%d days", int(delta.Hours()/24)), nil
}

func dateFormat(format string, content string) (string, error) {
	t, err := time.Parse("2006-01-02T15:04:05.000-0700", content)
	if err != nil {
		return "", err
	}
	return t.Format(format), nil
}

// // RunTemplate will run the give templateContent as a golang text/template
// // and pass the provided data to the template execution.  It will write
// // the output to the provided "out" writer.
// func RunTemplate(templateContent string, data interface{}, out io.Writer) error {
// 	return runTemplate(templateContent, data, out)
// }

// func responseToJSON(resp *http.Response, err error) (interface{}, error) {
// 	if err != nil {
// 		return nil, err
// 	}

// 	data := jsonDecode(resp.Body)
// 	if resp.StatusCode == 400 {
// 		if val, ok := data.(map[string]interface{})["errorMessages"]; ok {
// 			for _, errMsg := range val.([]interface{}) {
// 				log.Errorf("%s", errMsg)
// 			}
// 		}
// 	}

// 	return data, nil
// }

// func jsonDecode(io io.Reader) interface{} {
// 	content, err := ioutil.ReadAll(io)
// 	var data interface{}
// 	err = json.Unmarshal(content, &data)
// 	if err != nil {
// 		log.Errorf("JSON Parse Error: %s from %s", err, content)
// 	}
// 	return data
// }

// func jsonEncode(data interface{}) (string, error) {
// 	buffer := bytes.NewBuffer(make([]byte, 0))
// 	enc := json.NewEncoder(buffer)

// 	err := enc.Encode(data)
// 	if err != nil {
// 		log.Errorf("Failed to encode data %s: %s", data, err)
// 		return "", err
// 	}
// 	return buffer.String(), nil
// }

// func jsonWrite(file string, data interface{}) {
// 	fh, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
// 	defer fh.Close()
// 	if err != nil {
// 		log.Errorf("Failed to open %s: %s", file, err)
// 		os.Exit(1)
// 	}
// 	enc := json.NewEncoder(fh)
// 	enc.Encode(data)
// }

// func yamlWrite(file string, data interface{}) {
// 	fh, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
// 	defer fh.Close()
// 	if err != nil {
// 		log.Errorf("Failed to open %s: %s", file, err)
// 		os.Exit(1)
// 	}
// 	if out, err := yaml.Marshal(data); err != nil {
// 		log.Errorf("Failed to marshal yaml %v: %s", data, err)
// 		os.Exit(1)
// 	} else {
// 		fh.Write(out)
// 	}
// }

// func promptYN(prompt string, yes bool) bool {
// 	reader := bufio.NewReader(os.Stdin)
// 	if !yes {
// 		prompt = fmt.Sprintf("%s [y/N]: ", prompt)
// 	} else {
// 		prompt = fmt.Sprintf("%s [Y/n]: ", prompt)
// 	}

// 	fmt.Printf("%s", prompt)
// 	text, _ := reader.ReadString('\n')
// 	ans := strings.ToLower(strings.TrimRight(text, "\n"))
// 	if ans == "" {
// 		return yes
// 	}
// 	if strings.HasPrefix(ans, "y") {
// 		return true
// 	}
// 	return false
// }

// this is a HACK to make yaml parsed documents to be serializable
// to json, so prevent this:
// json: unsupported type: map[interface {}]interface {}
// Also we want to clean up common input errors for the edit
// templates, like dangling "\n"
func yamlFixup(data interface{}) (interface{}, error) {
	switch d := data.(type) {
	case map[interface{}]interface{}:
		// need to copy this map into a string map so json can encode it
		copy := make(map[string]interface{})
		for key, val := range d {
			switch k := key.(type) {
			case string:
				if fixed, err := yamlFixup(val); err != nil {
					return nil, err
				} else if fixed != nil {
					copy[k] = fixed
				}
			default:
				err := fmt.Errorf("YAML: key %s is type '%T', require 'string'", key, k)
				log.Errorf("%s", err)
				return nil, err
			}
		}
		if len(copy) == 0 {
			return nil, nil
		}
		return copy, nil
	case map[string]interface{}:
		copy := make(map[string]interface{})
		for k, v := range d {
			if fixed, err := yamlFixup(v); err != nil {
				return nil, err
			} else if fixed != nil {
				copy[k] = fixed
			}
		}
		if len(copy) == 0 {
			return nil, nil
		}
		return copy, nil
	case []interface{}:
		copy := make([]interface{}, 0, len(d))
		for _, val := range d {
			if fixed, err := yamlFixup(val); err != nil {
				return nil, err
			} else if fixed != nil {
				copy = append(copy, fixed)
			}
		}
		if len(copy) == 0 {
			return nil, nil
		}
		return copy, nil
	case *interface{}:
		if fixed, err := yamlFixup(*d); err != nil {
			return nil, err
		} else if fixed != nil {
			*d = fixed
		}
		return d, nil
	case string:
		if d == "" || d == "\n" {
			return nil, nil
		}
		return d, nil
	default:
		return d, nil
	}
}

// func mkdir(dir string) error {
// 	if stat, err := os.Stat(dir); err != nil && !os.IsNotExist(err) {
// 		log.Errorf("Failed to stat %s: %s", dir, err)
// 		return err
// 	} else if err == nil && !stat.IsDir() {
// 		err := fmt.Errorf("%s exists and is not a directory", dir)
// 		log.Errorf("%s", err)
// 		return err
// 	} else {
// 		// dir does not exist, so try to create it
// 		if err := os.MkdirAll(dir, 0755); err != nil {
// 			log.Errorf("Failed to mkdir -p %s: %s", dir, err)
// 			return err
// 		}
// 	}
// 	return nil
// }
