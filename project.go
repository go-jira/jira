package jira

import "github.com/arenstar/go-jira/jiradata"

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

// https://developer.atlassian.com/cloud/jira/platform/rest/v2#api-api-2-project-projectIdOrKey-versions-get
func (j *Jira) GetProjectVersions(project string) (*jiradata.ProjectVersions, error) {
	return GetProjectVersions(j.UA, j.Endpoint, project)
}

func GetProjectVersions(ua HttpClient, endpoint string, project string) (*jiradata.ProjectVersions, error) {
	uri := URLJoin(endpoint, "rest/api/2/project", project, "versions")
	resp, err := ua.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := jiradata.ProjectVersions{}
		return &results, readJSON(resp.Body, &results)
	}
	return nil, responseError(resp)
}
