package jira

import (
	"encoding/json"
	"net/http"

	"gopkg.in/jira.v1/jiradata"
)

func responseError(resp *http.Response) error {
	results := &jiradata.ErrorCollection{}
	if err := json.NewDecoder(resp.Body).Decode(results); err != nil {
		results.Status = resp.StatusCode
		results.ErrorMessages = append(results.ErrorMessages, err.Error())
	}
	if len(results.ErrorMessages) == 0 && len(results.Errors) == 0 {
		results.Status = resp.StatusCode
		results.ErrorMessages = append(results.ErrorMessages, resp.Status)
	}
	return results
}
