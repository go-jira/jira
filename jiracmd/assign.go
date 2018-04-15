package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type AssignOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
	Assignee              string `yaml:"assignee,omitempty" json:"assignee,omitempty"`
}

func CmdAssignRegistry() *jiracli.CommandRegistryEntry {
	opts := AssignOptions{}

	return &jiracli.CommandRegistryEntry{
		"Assign user to issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdAssignUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdAssign(o, globals, &opts)
		},
	}
}

func CmdAssignUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Flag("default", "use default user for assignee").PreAction(func(ctx *kingpin.ParseContext) error {
		if jiracli.FlagValue(ctx, "default") == "true" {
			opts.Assignee = "-1"
		}
		return nil
	}).Bool()
	cmd.Arg("ISSUE", "issue to assign").Required().StringVar(&opts.Issue)
	cmd.Arg("ASSIGNEE", "user to assign to issue").StringVar(&opts.Assignee)
	return nil
}

// CmdAssign will assign an issue to a user
func CmdAssign(o *oreo.Client, globals *jiracli.GlobalOptions, opts *AssignOptions) error {
	err := jira.IssueAssign(o, globals.Endpoint.Value, opts.Issue, opts.Assignee)
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
