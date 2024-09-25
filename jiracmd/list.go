package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ListOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	jira.SearchOptions    `yaml:",inline" json:",inline" figtree:",inline"`
	Queries               map[string]string `yaml:"queries,omitempty" json:"queries,omitempty"`
}

func CmdListRegistry() *jiracli.CommandRegistryEntry {
	opts := ListOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("list"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Prints list of issues for given search criteria",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdListUsage(cmd, &opts, fig)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			if opts.QueryFields == "" {
				opts.QueryFields = "assignee,created,priority,reporter,status,summary,updated,issuetype"
			}
			if opts.Sort == "" {
				opts.Sort = "priority asc, key"
			}
			return CmdList(o, globals, &opts)
		},
	}
}

func CmdListUsage(cmd *kingpin.CmdClause, opts *ListOptions, fig *figtree.FigTree) error {
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	jiracli.GJsonQueryUsage(cmd, &opts.CommonOptions)
	cmd.Flag("assignee", "User assigned the issue").Short('a').StringVar(&opts.Assignee)
	cmd.Flag("component", "Component to search for").Short('c').StringVar(&opts.Component)
	cmd.Flag("issuetype", "Issue type to search for").Short('i').StringVar(&opts.IssueType)
	cmd.Flag("limit", "Maximum number of results to return in search").Short('l').IntVar(&opts.MaxResults)
	cmd.Flag("project", "Project to search for").Short('p').StringVar(&opts.Project)
	cmd.Flag("named-query", "The name of a query in the `queries` configuration").Short('n').PreAction(func(ctx *kingpin.ParseContext) error {
		name := jiracli.FlagValue(ctx, "named-query")
		if query, ok := opts.Queries[name]; ok && query != "" {
			var err error
			opts.Query, err = jiracli.ConfigTemplate(fig, query, cmd.FullCommand(), opts)
			return err
		}
		return fmt.Errorf("A valid named-query %q not found in `queries` configuration", name)
	}).String()
	cmd.Flag("query", "Jira Query Language (JQL) expression for the search").Short('q').StringVar(&opts.Query)
	cmd.Flag("queryfields", "Fields that are used in \"list\" template").Short('f').StringVar(&opts.QueryFields)
	cmd.Flag("reporter", "Reporter to search for").Short('r').StringVar(&opts.Reporter)
	cmd.Flag("status", "Filter on issue status").Short('S').StringVar(&opts.Status)
	cmd.Flag("sort", "Sort order to return").Short('s').StringVar(&opts.Sort)
	cmd.Flag("watcher", "Watcher to search for").Short('w').StringVar(&opts.Watcher)
	return nil
}

// List will query jira and send data to "list" template
func CmdList(o *oreo.Client, globals *jiracli.GlobalOptions, opts *ListOptions) error {
	data, err := jira.Search(o, globals.Endpoint.Value, opts, jira.WithAutoPagination())
	if err != nil {
		return err
	}
	return opts.PrintTemplate(data)
}
