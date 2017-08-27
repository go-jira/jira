package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdUnassignRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := AssignOptions{}

	return &CommandRegistryEntry{
		"Unassign an issue",
		func() error {
			return CmdAssign(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdAssignUsage(cmd, &opts)
		},
	}
}

func CmdUnassignUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue to unassign").Required().StringVar(&opts.Issue)
	return nil
}
