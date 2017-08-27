package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdTakeRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := AssignOptions{}

	return &CommandRegistryEntry{
		"Assign issue to yourself",
		func() error {
			opts.Assignee = opts.User.Value
			return CmdAssign(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdAssignUsage(cmd, &opts)
		},
	}
}

func CmdTakeUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue to assign").Required().StringVar(&opts.Issue)
	return nil
}
