package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdTakeRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := AssignOptions{}

	return &jiracli.CommandRegistryEntry{
		"Assign issue to yourself",
		func() error {
			opts.Assignee = opts.User.Value
			return CmdAssign(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdAssignUsage(cmd, &opts)
		},
	}
}

func CmdTakeUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue to assign").Required().StringVar(&opts.Issue)
	return nil
}
