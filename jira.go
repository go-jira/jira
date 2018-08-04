package jira

import (
	"github.com/coryb/oreo"
	logging "gopkg.in/op/go-logging.v1"
)

var log = logging.MustGetLogger("jira")

const VERSION = "1.0.20"

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
