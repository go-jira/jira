package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

type SearchProvider interface {
	ProvideSearchRequest() *jiradata.SearchRequest
}

type SearchOptions struct {
	Assignee    string `yaml:"assignee,omitempty" json:"assignee,omitempty"`
	Query       string `yaml:"query,omitempty" json:"query,omitempty"`
	QueryFields string `yaml:"query-fields,omitempty" json:"query-fields,omitempty"`
	Project     string `yaml:"project,omitempty" json:"project,omitempty"`
	Component   string `yaml:"component,omitempty" json:"component,omitempty"`
	IssueType   string `yaml:"issue-type,omitempty" json:"issue-type,omitempty"`
	Watcher     string `yaml:"watcher,omitempty" json:"watcher,omitempty"`
	Reporter    string `yaml:"reporter,omitempty" json:"reporter,omitempty"`
	Status      string `yaml:"status,omitempty" json:"status,omitempty"`
	Sort        string `yaml:"sort,omitempty" json:"sort,omitempty"`
	MaxResults  int    `yaml:"max-results,omitempty" json:"max-results,omitempty"`
}

func (o *SearchOptions) ProvideSearchRequest() *jiradata.SearchRequest {
	req := &jiradata.SearchRequest{}

	if o.Query == "" {
		qbuff := bytes.NewBufferString("resolution = unresolved")
		if o.Project != "" {
			qbuff.WriteString(fmt.Sprintf(" AND project = '%s'", o.Project))
		}
		if o.Component != "" {
			qbuff.WriteString(fmt.Sprintf(" AND component = '%s'", o.Component))
		}
		if o.Assignee != "" {
			qbuff.WriteString(fmt.Sprintf(" AND assignee = '%s'", o.Assignee))
		}
		if o.IssueType != "" {
			qbuff.WriteString(fmt.Sprintf(" AND issuetype = '%s'", o.IssueType))
		}
		if o.Watcher != "" {
			qbuff.WriteString(fmt.Sprintf(" AND watcher = '%s'", o.Watcher))
		}
		if o.Reporter != "" {
			qbuff.WriteString(fmt.Sprintf(" AND reporter = '%s'", o.Reporter))
		}
		if o.Status != "" {
			qbuff.WriteString(fmt.Sprintf(" AND status = '%s'", o.Status))
		}
		if o.Sort != "" {
			qbuff.WriteString(fmt.Sprintf(" ORDER BY %s", o.Sort))
		}
		req.JQL = qbuff.String()
	} else {
		req.JQL = o.Query
	}

	req.Fields = append(req.Fields, "summary")
	if o.QueryFields != "" {
		fields := strings.Split(o.QueryFields, ",")
		req.Fields = append(req.Fields, fields...)
	}
	req.StartAt = 0
	req.MaxResults = o.MaxResults

	return req
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/search-searchUsingSearchRequest
func (j *Jira) Search(sp SearchProvider) (*jiradata.SearchResults, error) {
	return Search(j.UA, j.Endpoint, sp)
}

func Search(ua HttpClient, endpoint string, sp SearchProvider) (*jiradata.SearchResults, error) {
	req := sp.ProvideSearchRequest()
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	uri := URLJoin(endpoint, "rest/api/2/search")
	resp, err := ua.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.SearchResults{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}
