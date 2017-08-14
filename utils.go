package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

func readJSON(input io.Reader, data interface{}) error {
	content, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, data)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("JSON Parse Error: %s from %s", err, content))
	}
	return nil
}
