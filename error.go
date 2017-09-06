package jira

import (
	"fmt"
	"net/http"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

func responseError(resp *http.Response) error {
	results := &jiradata.ErrorCollection{}
	if err := readJSON(resp.Body, results); err != nil {
		return err
	}
	if len(results.ErrorMessages) == 0 && len(results.Errors) == 0 {
		return fmt.Errorf(resp.Status)
	}
	return results
}
