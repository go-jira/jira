package jira

import (
	"encoding/json"

	"github.com/sosheskaz/jira/jiradata"
)

func ServerInfo(ua HttpClient, endpoint string) (*jiradata.ServerInfo, error) {
	uri := URLJoin(endpoint, "rest/api/2/serverInfo")
	resp, err := ua.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := jiradata.ServerInfo{}
		return &results, json.NewDecoder(resp.Body).Decode(&results)
	}
	return nil, responseError(resp)
}
