package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type AssignOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Issue         string
	Assignee      string
}

func CmdAssignRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := AssignOptions{}

	return &CommandRegistryEntry{
		"Assign user to issue",
		func() error {
			return CmdAssign(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdAssignUsage(cmd, &opts)
		},
	}
}

func CmdAssignUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
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
func CmdAssign(o *oreo.Client, opts *AssignOptions) error {
	err := jira.IssueAssign(o, opts.Endpoint.Value, opts.Issue, opts.Assignee)
	if err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, opts.Endpoint.Value, opts.Issue)

	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}

	return nil
}
