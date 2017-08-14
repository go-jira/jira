package jiracli

import (
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WorklogAddOptions struct {
	GlobalOptions
	jiradata.Worklog
}

func (jc *JiraCli) CmdWorklogAddRegistry() *CommandRegistryEntry {
	issue := ""
	opts := WorklogAddOptions{
		GlobalOptions: GlobalOptions{
			Template: "worklog",
		},
	}
	return &CommandRegistryEntry{
		"Add a worklog to an issue",
		func() error {
			return jc.CmdWorklogAdd(issue, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdWorklogAddUsage(cmd, &issue, &opts)
		},
	}
}

func (jc *JiraCli) CmdWorklogAddUsage(cmd *kingpin.CmdClause, issue *string, opts *WorklogAddOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").BoolVar(&opts.SkipEditing)
	cmd.Flag("comment", "Comment message for worklog").Short('m').StringVar(&opts.Comment)
	cmd.Flag("time-spent", "Time spent working on issue").Short('T').StringVar(&opts.TimeSpent)
	cmd.Arg("ISSUE", "issue id to fetch worklogs").Required().StringVar(issue)
	return nil
}

// CmdWorklogAdd will attempt to add (action=add) a worklog to the given issue.
// It will spawn the editor (unless --noedit isused) and post edited YAML
// content as JSON to the worklog endpoint
func (jc *JiraCli) CmdWorklogAdd(issue string, opts *WorklogAddOptions) error {
	return jc.editLoop(&opts.GlobalOptions, &opts.Worklog, &opts.Worklog, func() error {
		_, err := jc.AddIssueWorklog(issue, opts)
		return err
	})
}
