package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WorklogAddOptions struct {
	GlobalOptions    `yaml:",inline" figtree:",inline"`
	jiradata.Worklog `yaml:",inline" figtree:",inline"`
	Issue            string
}

func CmdWorklogAddRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := WorklogAddOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("worklog"),
		},
	}
	return &CommandRegistryEntry{
		"Add a worklog to an issue",
		func() error {
			return CmdWorklogAdd(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdWorklogAddUsage(cmd, &opts)
		},
	}
}

func CmdWorklogAddUsage(cmd *kingpin.CmdClause, opts *WorklogAddOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	EditorUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("comment", "Comment message for worklog").Short('m').StringVar(&opts.Comment)
	cmd.Flag("time-spent", "Time spent working on issue").Short('T').StringVar(&opts.TimeSpent)
	cmd.Arg("ISSUE", "issue id to fetch worklogs").Required().StringVar(&opts.Issue)
	return nil
}

// CmdWorklogAdd will attempt to add (action=add) a worklog to the given issue.
// It will spawn the editor (unless --noedit isused) and post edited YAML
// content as JSON to the worklog endpoint
func CmdWorklogAdd(o *oreo.Client, opts *WorklogAddOptions) error {
	err := editLoop(&opts.GlobalOptions, &opts.Worklog, &opts.Worklog, func() error {
		_, err := jira.AddIssueWorklog(o, opts.Endpoint.Value, opts.Issue, opts)
		return err
	})
	if err != nil {
		return err
	}
	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
