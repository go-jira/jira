package jira

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-jira/jira/jiradata"
)

type UserSearchOptions struct {
	Query      string `yaml:"query,omitempty" json:"query,omitempty"`
	Username   string `yaml:"username,omitempty" json:"username,omitempty"`
	AccountID  string `yaml:"accountId,omitempty" json:"accountId,omitempty"`
	StartAt    int    `yaml:"startAt,omitempty" json:"startAt,omitempty"`
	MaxResults int    `yaml:"max-results,omitempty" json:"max-results,omitempty"`
	Property   string `yaml:"property,omitempty" json:"property,omitempty"`
}

// https://developer.atlassian.com/cloud/jira/platform/rest/v2/#api-rest-api-2-user-search-get

func UserSearch(ua HttpClient, endpoint string, opts *UserSearchOptions) ([]*jiradata.User, error) {
	uri := URLJoin(endpoint, "rest/api/2/user/search")
	params := []string{}
	if opts.Query != "" {
		params = append(params, "query="+url.QueryEscape(opts.Query))
	}
	if opts.AccountID != "" {
		params = append(params, "accountId="+url.QueryEscape(opts.AccountID))
	}
	if opts.StartAt != 0 {
		params = append(params, fmt.Sprintf("startAt=%d", opts.StartAt))
	}
	if opts.MaxResults != 0 {
		params = append(params, fmt.Sprintf("maxResults=%d", opts.MaxResults))
	}
	if opts.Property != "" {
		params = append(params, "property="+url.QueryEscape(opts.Property))
	}
	if len(params) > 0 {
		uri += "?" + strings.Join(params, "&")
	}
	resp, err := ua.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := []*jiradata.User{}
		return results, json.NewDecoder(resp.Body).Decode(&results)
	}
	return nil, responseError(resp)
}
