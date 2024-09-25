package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WorklogAddOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	jiradata.Worklog      `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdWorklogAddRegistry() *jiracli.CommandRegistryEntry {
	opts := WorklogAddOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("worklog"),
		},
	}
	return &jiracli.CommandRegistryEntry{
		"Add a worklog to an issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdWorklogAddUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdWorklogAdd(o, globals, &opts)
		},
	}
}

func CmdWorklogAddUsage(cmd *kingpin.CmdClause, opts *WorklogAddOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("comment", "Comment message for worklog").Short('m').StringVar(&opts.Comment)
	cmd.Flag("time-spent", "Time spent working on issue").Short('T').StringVar(&opts.TimeSpent)
	cmd.Flag("started", "Time you started work").Short('S').StringVar(&opts.Started)
	cmd.Arg("ISSUE", "issue id to fetch worklogs").Required().StringVar(&opts.Issue)
	return nil
}

// CmdWorklogAdd will attempt to add (action=add) a worklog to the given issue.
// It will spawn the editor (unless --noedit isused) and post edited YAML
// content as JSON to the worklog endpoint
func CmdWorklogAdd(o *oreo.Client, globals *jiracli.GlobalOptions, opts *WorklogAddOptions) error {
	err := jiracli.EditLoop(&opts.CommonOptions, &opts.Worklog, &opts.Worklog, func() error {
		_, err := jira.AddIssueWorklog(o, globals.Endpoint.Value, opts.Issue, opts)
		return err
	})
	if err != nil {
		return err
	}
	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", opts.Issue, jira.URLJoin(globals.Endpoint.Value, "browse", opts.Issue))
	}
	if opts.Browse.Value {
		return CmdBrowse(globals, opts.Issue)
	}
	return nil
}
