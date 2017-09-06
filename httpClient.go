package jira

import (
	"io"
	"net/http"
)

type HttpClient interface {
	Delete(url string) (*http.Response, error)
	Do(*http.Request) (*http.Response, error)
	GetJSON(url string) (*http.Response, error)
	Post(url, bodyType string, body io.Reader) (*http.Response, error)
	Put(url, bodyType string, body io.Reader) (*http.Response, error)
}
