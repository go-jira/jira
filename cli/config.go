package jiracli

import (
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/lib"
)

type JiraOptions struct {
	jira.SearchOptions `yaml:",inline"`
}
