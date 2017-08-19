package jiracli

import kingpin "gopkg.in/alecthomas/kingpin.v2"

func (jc *JiraCli) CmdTakeRegistry() *CommandRegistryEntry {
	opts := AssignOptions{}

	return &CommandRegistryEntry{
		"Assign issue to yourself",
		func() error {
			opts.Assignee = opts.User
			return jc.CmdAssign(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdAssignUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdTakeUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	cmd.Arg("ISSUE", "issue to assign").Required().StringVar(&opts.Issue)
	return nil
}
