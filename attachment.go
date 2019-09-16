package jira

import (
	"encoding/json"

	"github.com/go-jira/jira/jiradata"
)

// https://docs.atlassian.com/jira/REST/cloud/#api/2/attachment-getAttachment
func (j *Jira) GetAttachment(id string) (*jiradata.Attachment, error) {
	return GetAttachment(j.UA, j.Endpoint, id)
}

func GetAttachment(ua HttpClient, endpoint string, id string) (*jiradata.Attachment, error) {
	uri := URLJoin(endpoint, "rest/api/2/attachment", id)
	resp, err := ua.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.Attachment{}
		return results, json.NewDecoder(resp.Body).Decode(results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/attachment-removeAttachment
func (j *Jira) RemoveAttachment(id string) error {
	return RemoveAttachment(j.UA, j.Endpoint, id)
}

func RemoveAttachment(ua HttpClient, endpoint string, id string) error {
	uri := URLJoin(endpoint, "rest/api/2/attachment", id)
	resp, err := ua.Delete(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}
	return responseError(resp)
}
