package jiradata

// these routines implement the various Provide interfaces that are required by the jira functions

func (i *IssueUpdate) ProvideIssueUpdate() *IssueUpdate {
	return i
}

func (w *Worklog) ProvideWorklog() *Worklog {
	return w
}

func (l *LinkIssueRequest) ProvideLinkIssueRequest() *LinkIssueRequest {
	return l
}

func (r *RankRequest) ProvideRankRequest() *RankRequest {
	return r
}

func (c *Comment) ProvideComment() *Comment {
	return c
}

func (c *Component) ProvideComponent() *Component {
	return c
}

func (e *EpicIssues) ProvideEpicIssues() *EpicIssues {
	return e
}
