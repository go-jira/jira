package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func readJSON(input io.Reader, data interface{}) error {
	content, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return nil
	}
	err = json.Unmarshal(content, data)
	if err != nil {
		return fmt.Errorf("JSON Parse Error: %s from %q", err, content)
	}
	return nil
}
