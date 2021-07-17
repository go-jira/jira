package jira

import (
	"github.com/coryb/oreo"
)

// replace by ldflags
var VERSION = "development"

type Jira struct {
	Endpoint string     `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	UA       HttpClient `json:"-" yaml:"-"`
}

func NewJira(endpoint string) *Jira {
	return &Jira{
		Endpoint: endpoint,
		UA:       oreo.New(),
	}
}
