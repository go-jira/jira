package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

type IssueQueryProvider interface {
	ProvideIssueQueryString() string
}

type IssueOptions struct {
	Fields        []string `json:"fields,omitempty" yaml:"fields,omitempty"`
	Expand        []string `json:"expand,omitempty" yaml:"expand,omitempty"`
	Properties    []string `json:"properties,omitempty" yaml:"properties,omitempty"`
	FieldsByKeys  bool     `json:"fieldsByKeys,omitempty" yaml:"fieldsByKeys,omitempty"`
	UpdateHistory bool     `json:"updateHistory,omitempty" yaml:"updateHistory,omitempty"`
}

func (o *IssueOptions) ProvideIssueQueryString() string {
	params := []string{}
	if len(o.Fields) > 0 {
		params = append(params, "fields="+strings.Join(o.Fields, ","))
	}
	if len(o.Expand) > 0 {
		params = append(params, "expand="+strings.Join(o.Expand, ","))
	}
	if len(o.Properties) > 0 {
		params = append(params, "properties="+strings.Join(o.Properties, ","))
	}
	if o.FieldsByKeys {
		params = append(params, "fieldsByKeys=true")
	}
	if o.UpdateHistory {
		params = append(params, "updateHistory=true")
	}
	if len(params) > 0 {
		return "?" + strings.Join(params, "&")
	}
	return ""
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getIssue
func (j *Jira) GetIssue(issue string, iqg IssueQueryProvider) (*jiradata.Issue, error) {
	query := ""
	if iqg != nil {
		query = iqg.ProvideIssueQueryString()
	}
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s%s", j.Endpoint, issue, query)
	resp, err := j.UA.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.Issue{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue/{issueIdOrKey}/worklog-getIssueWorklog
func (j *Jira) GetIssueWorklog(issue string) (*jiradata.Worklogs, error) {
	startAt := 0
	total := 1
	maxResults := 100
	worklogs := jiradata.Worklogs{}
	for startAt < total {
		uri := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog?startAt=%d&maxResults=%d", j.Endpoint, issue, startAt, maxResults)
		resp, err := j.UA.GetJSON(uri)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			results := &jiradata.WorklogWithPagination{}
			err := readJSON(resp.Body, results)
			if err != nil {
				return nil, err
			}
			startAt = startAt + maxResults
			total = results.Total
			for _, worklog := range results.Worklogs {
				worklogs = append(worklogs, worklog)
			}
		} else {
			return nil, responseError(resp)
		}
	}
	return &worklogs, nil
}

type WorklogProvider interface {
	ProvideWorklog() *jiradata.Worklog
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue/{issueIdOrKey}/worklog-addWorklog
func (j *Jira) AddIssueWorklog(issue string, wp WorklogProvider) (*jiradata.Worklog, error) {
	req := wp.ProvideWorklog()
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog", j.Endpoint, issue)
	resp, err := j.UA.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		results := &jiradata.Worklog{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getEditIssueMeta
func (j *Jira) GetIssueEditMeta(issue string) (*jiradata.EditMeta, error) {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/editmeta", j.Endpoint, issue)
	resp, err := j.UA.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.EditMeta{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}

type IssueUpdateProvider interface {
	ProvideIssueUpdate() *jiradata.IssueUpdate
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-editIssue
func (j *Jira) EditIssue(issue string, iup IssueUpdateProvider) error {
	req := iup.ProvideIssueUpdate()
	encoded, err := json.Marshal(req)
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s", j.Endpoint, issue)
	resp, err := j.UA.Put(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-createIssue
func (j *Jira) CreateIssue(iup IssueUpdateProvider) (*jiradata.IssueCreateResponse, error) {
	req := iup.ProvideIssueUpdate()
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s/rest/api/2/issue", j.Endpoint)
	resp, err := j.UA.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		results := &jiradata.IssueCreateResponse{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getCreateIssueMeta
func (j *Jira) GetIssueCreateMetaProject(projectKey string) (*jiradata.CreateMetaProject, error) {
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&expand=projects.issuetypes.fields", j.Endpoint, projectKey)
	resp, err := j.UA.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.CreateMeta{}
		err = readJSON(resp.Body, results)
		if err != nil {
			return nil, err
		}
		for _, project := range results.Projects {
			if project.Key == projectKey {
				return project, nil
			}
		}
		return nil, fmt.Errorf("Project %s not found", projectKey)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getCreateIssueMeta
func (j *Jira) GetIssueCreateMetaIssueType(projectKey, issueTypeName string) (*jiradata.CreateMetaIssueType, error) {
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&issuetypeNames=%s&expand=projects.issuetypes.fields", j.Endpoint, projectKey, issueTypeName)
	resp, err := j.UA.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.CreateMeta{}
		err = readJSON(resp.Body, results)
		if err != nil {
			return nil, err
		}
		for _, project := range results.Projects {
			if project.Key == projectKey {
				for _, issueType := range project.Issuetypes {
					if issueType.Name == issueTypeName {
						return issueType, nil
					}
				}
			}
		}
		return nil, fmt.Errorf("Project %s and IssueType %s not found", projectKey, issueTypeName)
	}
	return nil, responseError(resp)
}

type LinkIssueProvider interface {
	ProvideLinkIssueRequest() *jiradata.LinkIssueRequest
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issueLink-linkIssues
func (j *Jira) LinkIssues(lip LinkIssueProvider) error {
	req := lip.ProvideLinkIssueRequest()
	encoded, err := json.Marshal(req)
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("%s/rest/api/2/issueLink", j.Endpoint)
	resp, err := j.UA.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		return nil
	}
	return responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-getTransitions
func (j *Jira) GetIssueTransitions(issue string) (*jiradata.TransitionsMeta, error) {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions?expand=transitions.fields", j.Endpoint, issue)
	resp, err := j.UA.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.TransitionsMeta{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-doTransition
func (j *Jira) TransitionIssue(issue string, iup IssueUpdateProvider) error {
	req := iup.ProvideIssueUpdate()
	encoded, err := json.Marshal(req)
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", j.Endpoint, issue)
	resp, err := j.UA.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issueLinkType-getIssueLinkTypes
func (j *Jira) GetIssueLinkTypes() (*jiradata.IssueLinkTypes, error) {
	uri := fmt.Sprintf("%s/rest/api/2/issueLinkType", j.Endpoint)
	resp, err := j.UA.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := struct {
			IssueLinkTypes jiradata.IssueLinkTypes
		}{
			IssueLinkTypes: jiradata.IssueLinkTypes{},
		}
		return &results.IssueLinkTypes, readJSON(resp.Body, &results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-addVote
func (j *Jira) IssueAddVote(issue string) error {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/votes", j.Endpoint, issue)
	resp, err := j.UA.Post(uri, "application/json", strings.NewReader("{}"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-removeVote
func (j *Jira) IssueRemoveVote(issue string) error {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/votes", j.Endpoint, issue)
	resp, err := j.UA.Delete(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}

type RankRequestProvider interface {
	ProvideRankRequest() *jiradata.RankRequest
}

// https://docs.atlassian.com/jira-software/REST/cloud/#agile/1.0/issue-rankIssues
func (j *Jira) RankIssues(rrp RankRequestProvider) error {
	req := rrp.ProvideRankRequest()
	encoded, err := json.Marshal(req)
	if err != nil {
		return err
	}
	uri := fmt.Sprintf("%s/rest/agile/1.0/issue/rank", j.Endpoint)
	resp, err := j.UA.Put(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-addWatcher
func (j *Jira) IssueAddWatcher(issue, user string) error {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/watchers", j.Endpoint, issue)
	resp, err := j.UA.Post(uri, "application/json", strings.NewReader(fmt.Sprintf("%q", user)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-addWatcher
func (j *Jira) IssueRemoveWatcher(issue, user string) error {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/watchers?username=%s", j.Endpoint, issue, user)
	resp, err := j.UA.Delete(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}

type CommentProvider interface {
	ProvideComment() *jiradata.Comment
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue/{issueIdOrKey}/comment-addComment
func (j *Jira) IssueAddComment(issue string, cp CommentProvider) (*jiradata.Comment, error) {
	req := cp.ProvideComment()
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/comment", j.Endpoint, issue)
	resp, err := j.UA.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		results := jiradata.Comment{}
		return &results, readJSON(resp.Body, &results)
	}
	return nil, responseError(resp)
}
