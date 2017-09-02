package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WorklogAddOptions struct {
	jiracli.GlobalOptions    `yaml:",inline" json:",inline" figtree:",inline"`
	jiradata.Worklog `yaml:",inline" json:",inline" figtree:",inline"`
	Issue            string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdWorklogAddRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := WorklogAddOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("worklog"),
		},
	}
	return &jiracli.CommandRegistryEntry{
		"Add a worklog to an issue",
		func() error {
			return CmdWorklogAdd(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdWorklogAddUsage(cmd, &opts)
		},
	}
}

func CmdWorklogAddUsage(cmd *kingpin.CmdClause, opts *WorklogAddOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	jiracli.EditorUsage(cmd, &opts.GlobalOptions)
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
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
	err := jiracli.EditLoop(&opts.GlobalOptions, &opts.Worklog, &opts.Worklog, func() error {
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
