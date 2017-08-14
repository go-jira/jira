package jiradata

type RankRequest struct {
	Issues          []string `json:"issues,omitempty" yaml:"issues,omitempty"`
	RankBeforeIssue string   `json:"rankBeforeIssue,omitempty" yaml:"rankBeforeIssue,omitempty"`
	RankAfterIssue  string   `json:"rankAfterIssue,omitempty" yaml:"rankAfterIssue,omitempty"`
}
