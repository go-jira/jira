package jiracmd

import (
	"fmt"
	"strings"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type AssignOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
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
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
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
	cmd.Arg("ASSIGNEE", "email or display name of user to assign to issue").StringVar(&opts.Assignee)
	return nil
}

// CmdAssign will assign an issue to a user
func CmdAssign(o *oreo.Client, globals *jiracli.GlobalOptions, opts *AssignOptions) error {
	if globals.JiraDeploymentType.Value == "" {
		serverInfo, err := jira.ServerInfo(o, globals.Endpoint.Value)
		if err != nil {
			return err
		}
		globals.JiraDeploymentType.Value = strings.ToLower(serverInfo.DeploymentType)
	}

	assignFunc := jira.IssueAssign
	if globals.JiraDeploymentType.Value == jiracli.CloudDeploymentType {
		if opts.Assignee != "" && opts.Assignee != "-1" {
			users, err := jira.UserSearch(o, globals.Endpoint.Value, &jira.UserSearchOptions{
				Query: opts.Assignee,
			})
			if err != nil {
				return err
			}
			if len(users) > 1 {
				return fmt.Errorf("Found %d accounts for users with username %q", len(users), opts.Assignee)
			} else if len(users) == 1 {
				opts.Assignee = users[0].AccountID
			}
		}
		assignFunc = jira.IssueAssignAccountID
	}

	err := assignFunc(o, globals.Endpoint.Value, opts.Issue, opts.Assignee)
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
