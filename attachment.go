package jira

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

// https://docs.atlassian.com/jira/REST/cloud/#api/2/attachment-getAttachment
func (j *Jira) GetAttachment(id string) (*jiradata.Attachment, error) {
	return GetAttachment(j.UA, j.Endpoint, id)
}

func GetAttachment(ua HttpClient, endpoint string, id string) (*jiradata.Attachment, error) {
	uri := fmt.Sprintf("%s/rest/api/2/attachment/%s", endpoint, id)
	resp, err := ua.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.Attachment{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#api/2/attachment-removeAttachment
func (j *Jira) RemoveAttachment(id string) error {
	return RemoveAttachment(j.UA, j.Endpoint, id)
}

func RemoveAttachment(ua HttpClient, endpoint string, id string) error {
	uri := fmt.Sprintf("%s/rest/api/2/attachment/%s", endpoint, id)
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
