package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path"
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

func URLJoin(endpoint string, paths ...string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(fmt.Errorf("Unable to parse endpoint: %s", endpoint))
	}
	paths = append([]string{u.Path}, paths...)
	u.Path = path.Join(paths...)
	return u.String()
}
