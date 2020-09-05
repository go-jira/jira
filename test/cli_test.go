package test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	endpoint       = "https://go-jira.atlassian.net"
	goJiraApiToken = "Rw1cPlKI40TJeEl1Pj88A5ED"
	goJiraLogin    = "gojira@corybennett.org"
	mothraApiToken = "UNXrI9gq5p0LWUtblAxDA7A6"
	mothraLogin    = "mothra@corybennett.org"
)

var jira string = "../dist/github.com/go-jira/jira-linux-amd64"

func Test_CLI(t *testing.T) {
	// setup the jira cli environment
	jira, err := filepath.Abs(jira)
	if err != nil {
		t.Fatal(err)
	}
	if !filepath.IsAbs(jira) {
		t.Fatalf("could not obtain absolute path to jira binary")
	}

	if _, err := os.Stat(jira); err != nil {
		t.Fatalf("could not stat %v: %v", jira, err)
	}

	os.Setenv("COLUMNS", "149")
	os.Setenv("JIRA_LOG_FORMAT", "%{level:-5s} %{message}")
	os.Setenv("ENDPOINT", endpoint)
	os.Setenv("JIRACLOUD", "1")

	t.Run("basic", test_Basic)
	t.Run("pagination", test_Pagination)
}

// test_Basic will test the basic functionality required in the cli
func test_Basic(t *testing.T) {
	// we'll reassign these often, just create
	// them here.
	var buf bytes.Buffer
	var err error

	// Create an issue
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		createIssue(
			jira,
			"BASIC",
			"summary",
			"description",
		),
	)
	if err != nil {
		t.Fatalf("cmd failed. stdout: %v err: %v", buf.String(), err)
	}
	issue := checkCreateIssue(t, buf, endpoint)

	// View the issue
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		viewIssue(
			jira,
			issue,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
priority: Medium
votes: 0
description: |
  description
`, issue)

	// confirm new issue shows in project list
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		listIssues(
			jira,
			"BASIC",
			"", // empty string means do not use a named query
			"", // empty string means do not use raw query
			"", // empty string means do not use a template
			"", // empty string means do not limit response
		),
	)
	checkIssueInOutput(t, buf, issue)

	// confirm issue appears with named query
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		listIssues(
			jira,
			"BASIC",
			"todo",
			"",
			"",
			"",
		),
	)
	checkIssueInOutput(t, buf, issue)

	// confirm issue appears with table template
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		listIssues(
			jira,
			"BASIC",
			"",
			"",
			"table",
			"",
		),
	)
	checkIssueInOutput(t, buf, issue)

	// edit an issue
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		editIssue(
			jira,
			issue,
			"edit comment",
			"priority=High",
			"",
		),
	)
	checkEditIssue(t, buf, issue, endpoint)

	// edit multiple issues with query and check comments updated
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		editIssue(
			jira,
			issue,
			"bulk edit comment",
			"priority=High",
			`resolution = unresolved AND project = 'BASIC' AND status = 'To Do'`,
		),
	)
	checkEditIssue(t, buf, issue, endpoint)
	// view the issue
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		viewIssue(
			jira,
			issue,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
priority: High
votes: 0
description: |
  description

comments:
  - | # GoJira, a minute ago
    edit comment
  - | # GoJira, a minute ago
    bulk edit comment

`, issue)

	// try invalid close of issue
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		closeIssue(
			jira,
			issue,
		),
	)
	checkDiff(t, buf, `ERROR Invalid Transition "close" from "To Do", Available: To Do, In Progress, In Review, Done
`)

	// put issue in done state
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		doneIssue(
			jira,
			issue,
		),
	)
	checkEditIssue(t, buf, issue, endpoint)

	// make sure our resolved issue is not present in the project
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		listIssues(
			jira,
			"BASIC",
			"",
			"",
			"",
			"",
		),
	)
	checkIssueNotInOutput(t, buf, issue)

	// create two new issues to test duplicating
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		createIssue(
			jira,
			"BASIC",
			"summary",
			"description",
		),
	)
	issue = checkCreateIssue(t, buf, endpoint)
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		createIssue(
			jira,
			"BASIC",
			"dup",
			"dup",
		),
	)
	dup := checkCreateIssue(t, buf, endpoint)

	// mark issue as duplicate
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		dupIssue(
			jira,
			issue,
			dup,
		),
	)
	checkDupIssue(t, buf, issue, dup, endpoint)
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		viewIssue(
			jira,
			issue,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: %s[Done]
priority: Medium
votes: 0
description: |
  description
`, issue, dup)

	// check dup is resolved and not in listed issue
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		listIssues(jira,
			"BASIC",
			"",
			"",
			"",
			""),
	)
	checkIssueNotInOutput(t, buf, dup)

	// create blocker issue
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		createIssue(
			jira,
			"BASIC",
			"blocks",
			"blocks",
		),
	)
	blocker := checkCreateIssue(t, buf, endpoint)

	// set blocker
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		blockIssue(
			jira,
			blocker,
			issue,
		),
	)
	checkBlockIssue(t, buf, blocker, issue, endpoint)

	// confirm blocker shows up when viewing issue
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		viewIssue(
			jira,
			issue,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: %s[To Do]
depends: %s[Done]
priority: Medium
votes: 0
description: |
  description
`, issue, blocker, dup)

	// confirm both issues are unresolved
	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		listIssues(
			jira,
			"BASIC",
			"",
			"",
			"",
			"",
		),
	)
	checkIssueInOutput(t, buf, issue)
	checkIssueInOutput(t, buf, blocker)

	//                          //
	// begin using mothra user  //
	//                          //

	// use mothra to vote for main issue
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		voteIssue(
			jira,
			issue,
			false,
		),
	)
	checkEditIssue(t, buf, issue, endpoint)

	// view issue to confirm vote
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		viewIssue(
			jira,
			issue,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: %s[To Do]
depends: %s[Done]
priority: Medium
votes: 1
description: |
  description
`, issue, blocker, dup)

	// down vote and confirm
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		voteIssue(
			jira,
			issue,
			true, // down vote true
		),
	)
	checkEditIssue(t, buf, issue, endpoint)

	// view issue to confirm vote
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		viewIssue(
			jira,
			issue,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: %s[To Do]
depends: %s[Done]
priority: Medium
votes: 0
description: |
  description
`, issue, blocker, dup)

	// TODO(louis): skipping watcher test for now until a
	// "watchers" command is implemented.

	// set blocker to "In Progress"
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		transIssue(
			jira,
			"In Progress",
			blocker,
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	// set it back to "To Do"
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		todoIssue(
			jira,
			blocker,
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	// set it to "In Review"
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		transIssue(
			jira,
			"review",
			blocker,
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	// set it back to "To Do"
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		todoIssue(
			jira,
			blocker,
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	// set it to in progress
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		progIssue(
			jira,
			blocker,
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	// set it to in done
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		doneIssue(
			jira,
			blocker,
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	// confirm blocker is done
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		viewIssue(
			jira,
			issue,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: %s[Done]
depends: %s[Done]
priority: Medium
votes: 0
description: |
  description
`, issue, blocker, dup)

	// verify we can add comment
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		commentIssue(
			jira,
			issue,
			"Yo, Comment",
		),
	)
	checkEditIssue(t, buf, issue, endpoint)

	// verify we can see comment
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		viewIssue(
			jira,
			issue,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: %s[Done]
depends: %s[Done]
priority: Medium
votes: 0
description: |
  description

comments:
  - | # Mothra, a minute ago
    Yo, Comment

`, issue, blocker, dup)

	// verify we can add labels
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		addLabelsIssue(
			jira,
			blocker,
			"test-label",
			"another-label",
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		viewIssue(
			jira,
			blocker,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: Done
summary: blocks
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: %s[To Do]
priority: Medium
votes: 0
labels: another-label, test-label
description: |
  blocks
`, blocker, issue)

	// verify we can remove labels
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		removeLabelsIssue(
			jira,
			blocker,
			"another-label",
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		viewIssue(
			jira,
			blocker,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: Done
summary: blocks
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: %s[To Do]
priority: Medium
votes: 0
labels: test-label
description: |
  blocks
`, blocker, issue)

	// verify we can replace labels
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		setLabelsIssue(
			jira,
			blocker,
			"more-label",
			"better-label",
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		viewIssue(
			jira,
			blocker,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: Done
summary: blocks
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: %s[To Do]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
`, blocker, issue)

	// verify mothra can take an issue
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		takeIssue(
			jira,
			blocker,
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		viewIssue(
			jira,
			blocker,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: Done
summary: blocks
project: BASIC
issuetype: Bug
assignee: Mothra
reporter: GoJira
blockers: 
depends: %s[To Do]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
`, blocker, issue)

	// verify martha can give the issue back
	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		giveIssue(
			jira,
			blocker,
			"gojira",
		),
	)
	checkEditIssue(t, buf, blocker, endpoint)

	buf, err = withApiLogin(
		mothraLogin,
		mothraApiToken,
		viewIssue(
			jira,
			blocker,
		),
	)
	checkDiff(t, buf, `issue: %s
created: a minute ago
status: Done
summary: blocks
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: %s[To Do]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
`, blocker, issue)
}

func test_Pagination(t *testing.T) {
	var buf bytes.Buffer
	var err error

	// note:
	// we test limit+1 to handle extra newline split

	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		listIssues(
			jira,
			"BASIC",
			"",
			"project = 'BASIC' AND status = 'Done'", // query
			"",
			"102",
		),
	)
	if err != nil {
		t.Fatalf("failed to list issues. stderr:%v err: %v", buf.String(), err)
	}
	if len(strings.Split(buf.String(), "\n")) != 103 {
		t.Fatalf("got: %v want: %v", len(strings.Split(buf.String(), "\n")), 103)
	}

	buf, err = withApiLogin(
		goJiraLogin,
		goJiraApiToken,
		listIssues(
			jira,
			"BASIC",
			"",
			"project = 'BASIC' AND status = 'Done'", // query
			"",                                      // empty string means do not use a template
			"1",
		),
	)
	if err != nil {
		t.Fatalf("failed to list issues. stderr:%v err: %v", buf.String(), err)
	}
	if len(strings.Split(buf.String(), "\n")) != 2 {
		t.Fatalf("got: %v want: %v", len(strings.Split(buf.String(), "\n")), 2)
	}
}
