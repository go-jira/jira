package jiracli

import kingpin "gopkg.in/alecthomas/kingpin.v2"

func (jc *JiraCli) CmdUnassignRegistry() *CommandRegistryEntry {
	opts := AssignOptions{}

	return &CommandRegistryEntry{
		"Unassign an issue",
		func() error {
			return jc.CmdAssign(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdAssignUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdUnassignUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	cmd.Arg("ISSUE", "issue to unassign").Required().StringVar(&opts.Issue)
	return nil
}
