package jira

import (
	"bytes"
	"encoding/json"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

type AuthProvider interface {
	ProvideAuthParams() *jiradata.AuthParams
}

type AuthOptions struct {
	Username string
	Password string
}

func (a *AuthOptions) ProvideAuthParams() *jiradata.AuthParams {
	return &jiradata.AuthParams{
		Username: a.Username,
		Password: a.Password,
	}
}

// https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-login
func (j *Jira) NewSession(ap AuthProvider) (*jiradata.AuthSuccess, error) {
	return NewSession(j.UA, j.Endpoint, ap)
}

func NewSession(ua HttpClient, endpoint string, ap AuthProvider) (*jiradata.AuthSuccess, error) {
	req := ap.ProvideAuthParams()
	encoded, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	uri := URLJoin(endpoint, "rest/auth/1/session")
	resp, err := ua.Post(uri, "application/json", bytes.NewBuffer(encoded))
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
	return GetSession(j.UA, j.Endpoint)
}

func GetSession(ua HttpClient, endpoint string) (*jiradata.CurrentUser, error) {
	uri := URLJoin(endpoint, "rest/auth/1/session")
	resp, err := ua.GetJSON(uri)
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
	return DeleteSession(j.UA, j.Endpoint)
}

func DeleteSession(ua HttpClient, endpoint string) error {
	uri := URLJoin(endpoint, "rest/auth/1/session")
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
