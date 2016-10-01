package jira

import (
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/data"
	"bytes"
	"fmt"
	"strings"
)

type SearchProvider interface {
	SearchRequest() *jiradata.SearchRequest
}

type SearchOptions struct {
	Assignee string
	Query string
	QueryFields string
	Project string
	Component string
	IssueType string
	Watcher string
	Reporter string
	Sort string
	MaxResults int
}

func (o *SearchOptions) SearchRequest() *jiradata.SearchRequest {
	req := &jiradata.SearchRequest{}

	if o.Query == "" {
		qbuff := bytes.NewBufferString("resolution = unresolved")
		if o.Project == "" {
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
