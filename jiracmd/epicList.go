package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EpicListOptions struct {
	ListOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Epic        string `yaml:"epic,omitempty" json:"epic,omitempty"`
}

func CmdEpicListRegistry() *jiracli.CommandRegistryEntry {
	opts := EpicListOptions{
		ListOptions: ListOptions{
			CommonOptions: jiracli.CommonOptions{
				Template: figtree.NewStringOption("epic-list"),
			},
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Prints list of issues for an epic with optional search criteria",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdEpicListUsage(cmd, &opts, fig)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Epic = jiracli.FormatIssue(opts.Epic, opts.Project)
			if opts.MaxResults == 0 {
				opts.MaxResults = 500
			}
			if opts.QueryFields == "" {
				opts.QueryFields = "assignee,created,priority,reporter,status,summary,updated,issuetype"
			}
			if opts.Sort == "" {
				opts.Sort = "priority asc, key"
			}
			return CmdEpicList(o, globals, &opts)
		},
	}
}

func CmdEpicListUsage(cmd *kingpin.CmdClause, opts *EpicListOptions, fig *figtree.FigTree) error {
	CmdListUsage(cmd, &opts.ListOptions, fig)
	cmd.Arg("EPIC", "Epic Key or ID to list").Required().StringVar(&opts.Epic)
	return nil
}

func CmdEpicList(o *oreo.Client, globals *jiracli.GlobalOptions, opts *EpicListOptions) error {
	data, err := jira.EpicSearch(o, globals.Endpoint.Value, opts.Epic, opts)
	if err != nil {
		return err
	}
	return opts.PrintTemplate(data)
}
