package jira

import (
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

// https://docs.atlassian.com/jira/REST/cloud/#api/2/project-getProjectComponents
func (j *Jira) GetProjectComponents(project string) (*jiradata.Components, error) {
	return GetProjectComponents(j.UA, j.Endpoint, project)
}

func GetProjectComponents(ua HttpClient, endpoint string, project string) (*jiradata.Components, error) {
	uri := URLJoin(endpoint, "rest/api/2/project", project, "components")
	resp, err := ua.GetJSON(uri)
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
