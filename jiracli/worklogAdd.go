package jiracli

import (
	"github.com/coryb/figtree"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WorklogAddOptions struct {
	GlobalOptions
	jiradata.Worklog
	Issue string
}

func (jc *JiraCli) CmdWorklogAddRegistry() *CommandRegistryEntry {
	opts := WorklogAddOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("worklog"),
		},
	}
	return &CommandRegistryEntry{
		"Add a worklog to an issue",
		func() error {
			return jc.CmdWorklogAdd(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdWorklogAddUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdWorklogAddUsage(cmd *kingpin.CmdClause, opts *WorklogAddOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("comment", "Comment message for worklog").Short('m').StringVar(&opts.Comment)
	cmd.Flag("time-spent", "Time spent working on issue").Short('T').StringVar(&opts.TimeSpent)
	cmd.Arg("ISSUE", "issue id to fetch worklogs").Required().StringVar(&opts.Issue)
	return nil
}

// CmdWorklogAdd will attempt to add (action=add) a worklog to the given issue.
// It will spawn the editor (unless --noedit isused) and post edited YAML
// content as JSON to the worklog endpoint
func (jc *JiraCli) CmdWorklogAdd(opts *WorklogAddOptions) error {
	err := jc.editLoop(&opts.GlobalOptions, &opts.Worklog, &opts.Worklog, func() error {
		_, err := jc.AddIssueWorklog(opts.Issue, opts)
		return err
	})
	if err != nil {
		return err
	}
	if opts.Browse.Value {
		return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
