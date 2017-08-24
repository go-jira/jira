package jiracli

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type AssignOptions struct {
	GlobalOptions
	Issue    string
	Assignee string
}

func (jc *JiraCli) CmdAssignRegistry() *CommandRegistryEntry {
	opts := AssignOptions{}

	return &CommandRegistryEntry{
		"Assign user to issue",
		func() error {
			return jc.CmdAssign(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdAssignUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdAssignUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("default", "use default user for assignee").PreAction(func(ctx *kingpin.ParseContext) error {
		if flagValue(ctx, "default") == "true" {
			opts.Assignee = "-1"
		}
		return nil
	}).Bool()
	cmd.Arg("ISSUE", "issue to assign").Required().StringVar(&opts.Issue)
	cmd.Arg("ASSIGNEE", "user to assign to issue").StringVar(&opts.Assignee)
	return nil
}

// CmdAssign will assign an issue to a user
func (jc *JiraCli) CmdAssign(opts *AssignOptions) error {
	err := jc.IssueAssign(opts.Issue, opts.Assignee)
	if err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)

	if opts.Browse.Value {
		return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}

	return nil
}
