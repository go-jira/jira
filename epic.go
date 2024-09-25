package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/eroshan/oreo"

	"github.com/go-jira/jira/jiradata"
)

// https://docs.atlassian.com/jira-software/REST/latest/#agile/1.0/epic-getIssuesForEpic
func (j *Jira) EpicSearch(epic string, sp SearchProvider) (*jiradata.SearchResults, error) {
	return EpicSearch(j.UA, j.Endpoint, epic, sp)
}

func EpicSearch(ua HttpClient, endpoint string, epic string, sp SearchProvider) (*jiradata.SearchResults, error) {
	req := sp.ProvideSearchRequest()
	// encoded, err := json.Marshal(req)
	// if err != nil {
	// 	return nil, err
	// }
	uri, err := url.Parse(URLJoin(endpoint, "rest/agile/1.0/epic", epic, "issue"))
	if err != nil {
		return nil, err
	}
	params := url.Values{}
	if len(req.Fields) > 0 {
		params.Add("fields", strings.Join(req.Fields, ","))
	}
	if req.JQL != "" {
		params.Add("jql", req.JQL)
	}
	if req.MaxResults != 0 {
		params.Add("maxResults", fmt.Sprintf("%d", req.MaxResults))
	}
	if req.StartAt != 0 {
		params.Add("startAt", fmt.Sprintf("%d", req.StartAt))
	}
	if req.ValidateQuery != "" {
		params.Add("validateQuery", req.ValidateQuery)
	}
	uri.RawQuery = params.Encode()

	resp, err := ua.Do(oreo.RequestBuilder(uri).WithHeader("Accept", "application/json").Build())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.SearchResults{}
		return results, json.NewDecoder(resp.Body).Decode(results)
	}
	return nil, responseError(resp)
}

type EpicIssuesProvider interface {
	ProvideEpicIssues() *jiradata.EpicIssues
}

// https://docs.atlassian.com/jira-software/REST/latest/#agile/1.0/epic-moveIssuesToEpic
func (j *Jira) EpicAddIssues(epic string, eip EpicIssuesProvider) error {
	return EpicAddIssues(j.UA, j.Endpoint, epic, eip)
}

func EpicAddIssues(ua HttpClient, endpoint string, epic string, eip EpicIssuesProvider) error {
	req := eip.ProvideEpicIssues()
	encoded, err := json.Marshal(req)
	if err != nil {
		return err
	}

	uri := URLJoin(endpoint, "rest/agile/1.0/epic", epic, "issue")
	resp, err := ua.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}

// https://docs.atlassian.com/jira-software/REST/latest/#agile/1.0/epic-removeIssuesFromEpic
func (j *Jira) EpicRemoveIssues(eip EpicIssuesProvider) error {
	return EpicRemoveIssues(j.UA, j.Endpoint, eip)
}

func EpicRemoveIssues(ua HttpClient, endpoint string, eip EpicIssuesProvider) error {
	req := eip.ProvideEpicIssues()
	encoded, err := json.Marshal(req)
	if err != nil {
		return err
	}

	uri := URLJoin(endpoint, "rest/agile/1.0/epic/none/issue")
	resp, err := ua.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}
