package test

import (
	"os/exec"
)

func session(jira string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"session",
	)
	return cmd
}

func createIssue(jira, project, summary, description string) *exec.Cmd {
	sum := "summary=" + summary
	desc := "description=" + description
	cmd := exec.Command(
		jira,
		"create",
		"--project", project,
		"-o", sum,
		"-o", desc,
		"--noedit",
	)
	return cmd
}

func viewIssue(jira, issue string) *exec.Cmd {
	return exec.Command(
		jira,
		"view",
		issue,
	)
}

func listIssues(jira, project, query, rawquery, template, limit string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"ls",
		"--project",
		project,
	)
	if query != "" {
		cmd.Args = append(cmd.Args, "-n", query)
	}
	if rawquery != "" {
		cmd.Args = append(cmd.Args, "-q", rawquery)
	}
	if template != "" {
		cmd.Args = append(cmd.Args, "--template", template)
	}
	if limit != "" {
		cmd.Args = append(cmd.Args, "--limit", limit)
	}
	return cmd
}

func editIssue(jira, issue, message, override, query string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"edit",
		issue,
		"-m",
		message,
		"--override",
		override,
		"--noedit",
	)
	if query != "" {
		cmd.Args = append(cmd.Args, "--query", query)
	}
	return cmd
}

func closeIssue(jira, issue string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"close",
		issue,
	)
	return cmd
}

func doneIssue(jira, issue string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"done",
		issue,
	)
	return cmd
}

func dupIssue(jira, issue, duplicate string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"dup",
		duplicate,
		issue,
	)
	return cmd
}

func blockIssue(jira, blocker, issue string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"block",
		blocker,
		issue,
	)
	return cmd
}

func voteIssue(jira, issue string, down bool) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"vote",
		issue,
	)
	if down {
		cmd.Args = append(cmd.Args, "--down")
	}
	return cmd
}

func watchIssue(jira, issue string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"watch",
		issue,
	)
	return cmd
}

func transIssue(jira, trans, issue string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"trans",
		trans,
		issue,
		"--noedit",
	)
	return cmd
}

func todoIssue(jira, issue string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"todo",
		issue,
	)
	return cmd
}

func progIssue(jira, issue string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"prog",
		issue,
	)
	return cmd
}

func commentIssue(jira, issue, comment string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"comment",
		issue,
		"--noedit",
		"-m",
		comment,
	)
	return cmd
}

func addLabelsIssue(jira, issue string, labels ...string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"labels",
		"add",
		issue,
	)
	cmd.Args = append(cmd.Args, labels...)
	return cmd
}

func removeLabelsIssue(jira, issue string, labels ...string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"labels",
		"remove",
		issue,
	)
	cmd.Args = append(cmd.Args, labels...)
	return cmd
}

func setLabelsIssue(jira, issue string, labels ...string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"labels",
		"set",
		issue,
	)
	cmd.Args = append(cmd.Args, labels...)
	return cmd
}

func takeIssue(jira, issue string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"take",
		issue,
	)
	return cmd
}

func giveIssue(jira, issue, taker string) *exec.Cmd {
	cmd := exec.Command(
		jira,
		"give",
		issue,
		taker,
	)
	return cmd
}
