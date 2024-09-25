package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdTakeRegistry() *jiracli.CommandRegistryEntry {
	opts := AssignOptions{}

	return &jiracli.CommandRegistryEntry{
		"Assign issue to yourself",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdAssignUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			if opts.Assignee == "" {
				opts.Assignee = globals.Login.Value
			}
			return CmdAssign(o, globals, &opts)
		},
	}
}

func CmdTakeUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "issue to assign").Required().StringVar(&opts.Issue)
	return nil
}
