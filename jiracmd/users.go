package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type UsersOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
}

func CmdUsersRegistry() *jiracli.CommandRegistryEntry {
	opts := UsersOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("json"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"list org users",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdUsersUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdUsers(o, globals, &opts)
		},
	}
}

func CmdUsersUsage(cmd *kingpin.CmdClause, opts *UsersOptions) error {
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	return nil
}

func CmdUsers(o *oreo.Client, globals *jiracli.GlobalOptions, opts *UsersOptions) error {
	users, err := jira.UserSearch(o, globals.Endpoint.Value, &jira.UserSearchOptions{
		Query: "*",
		MaxResults: 1000,
	})
	if err != nil {
		return err
	}

	if err := opts.PrintTemplate(users); err != nil {
		return err
	}
	return nil
}
