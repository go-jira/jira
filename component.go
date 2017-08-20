package jira

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

type ComponentProvider interface {
	ProvideComponent() *jiradata.Component
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/component-createComponent
func (j *Jira) CreateComponent(cp ComponentProvider) (*jiradata.Component, error) {
	req := cp.ProvideComponent()
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s/rest/api/2/component", j.Endpoint)
	resp, err := j.UA.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		results := &jiradata.Component{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}
