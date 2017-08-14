package jira

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

type AuthProvider interface {
	ProvideAuthParams() *jiradata.AuthParams
}

type AuthOptions struct {
	Username string
	Password string
}

func (a *AuthOptions) AuthParams() *jiradata.AuthParams {
	return &jiradata.AuthParams{
		Username: a.Username,
		Password: a.Password,
	}
}

// https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-login
func (j *Jira) NewSession(ap AuthProvider) (*jiradata.AuthSuccess, error) {
	req := ap.ProvideAuthParams()
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("%s/rest/auth/1/session", j.Endpoint)
	resp, err := j.UA.Post(uri, "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.AuthSuccess{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-currentUser
func (j *Jira) GetSession() (*jiradata.CurrentUser, error) {
	uri := fmt.Sprintf("%s/rest/auth/1/session", j.Endpoint)
	resp, err := j.UA.GetJSON(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		results := &jiradata.CurrentUser{}
		return results, readJSON(resp.Body, results)
	}
	return nil, responseError(resp)
}

// https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-logout
func (j *Jira) DeleteSession() error {
	uri := fmt.Sprintf("%s/rest/auth/1/session", j.Endpoint)
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
