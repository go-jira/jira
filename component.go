package jira

import (
	"bytes"
	"encoding/json"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

type ComponentProvider interface {
	ProvideComponent() *jiradata.Component
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/component-createComponent
func (j *Jira) CreateComponent(cp ComponentProvider) (*jiradata.Component, error) {
	return CreateComponent(j.UA, j.Endpoint, cp)
}

func CreateComponent(ua HttpClient, endpoint string, cp ComponentProvider) (*jiradata.Component, error) {
	req := cp.ProvideComponent()
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	uri := URLJoin(endpoint, "rest/api/2/component")
	resp, err := ua.Post(uri, "application/json", bytes.NewBuffer(encoded))
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
