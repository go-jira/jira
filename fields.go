package jira

import (
	"encoding/json"

	"github.com/go-jira/jira/jiradata"
)

// https://docs.atlassian.com/jira/REST/cloud/#api/3/field-getFields
func (j *Jira) GetFields() ([]jiradata.Field, error) {
	return GetFields(j.UA, j.Endpoint)
}

func GetFields(ua HttpClient, endpoint string) ([]jiradata.Field, error) {
	uri := URLJoin(endpoint, "rest/api/3/field")
	resp, err := ua.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		results := []jiradata.Field{}
		return results, json.NewDecoder(resp.Body).Decode(&results)
	}
	return nil, responseError(resp)
}
