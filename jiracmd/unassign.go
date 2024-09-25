package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdUnassignRegistry() *jiracli.CommandRegistryEntry {
	opts := AssignOptions{}

	return &jiracli.CommandRegistryEntry{
		"Unassign an issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdAssignUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdAssign(o, globals, &opts)
		},
	}
}

func CmdUnassignUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "issue to unassign").Required().StringVar(&opts.Issue)
	return nil
}
