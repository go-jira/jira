package jira

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

// https://docs.atlassian.com/jira/REST/cloud/#api/2/field-getFields
func (j *Jira) GetFields() ([]jiradata.Field, error) {
	return GetFields(j.UA, j.Endpoint)
}

func GetFields(ua HttpClient, endpoint string) ([]jiradata.Field, error) {
	uri := fmt.Sprintf("%s/rest/api/2/field", endpoint)
	resp, err := ua.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		results := []jiradata.Field{}
		return results, readJSON(resp.Body, &results)
	}
	return nil, responseError(resp)
}
