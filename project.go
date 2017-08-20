package jira

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

// https://docs.atlassian.com/jira/REST/cloud/#api/2/project-getProjectComponents
func (j *Jira) GetProjectComponents(project string) (*jiradata.Components, error) {
	uri := fmt.Sprintf("%s/rest/api/2/project/%s/components", j.Endpoint, project)
	resp, err := j.UA.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := jiradata.Components{}
		return &results, readJSON(resp.Body, &results)
	}
	return nil, responseError(resp)
}
