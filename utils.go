package jira

import (
	"fmt"
	"net/url"
	"path"
)

func URLJoin(endpoint string, paths ...string) string {
	u, err := url.Parse(endpoint)
	if err != nil {
		panic(fmt.Errorf("unable to parse endpoint: %s", endpoint))
	}
	paths = append([]string{u.Path}, paths...)
	u.Path = path.Join(paths...)
	return u.String()
}
