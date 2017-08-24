package jiracli

import (
	"github.com/coryb/figtree"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ListOptions struct {
	GlobalOptions      `yaml:",inline" figtree:",inline"`
	jira.SearchOptions `yaml:",inline" figtree:",inline"`
}

func (jc *JiraCli) CmdListRegistry() *CommandRegistryEntry {
	opts := ListOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("list"),
		},
		SearchOptions: jira.SearchOptions{
			MaxResults:  500,
			QueryFields: "assignee,created,priority,reporter,status,summary,updated",
			Sort:        "priority asc, key",
		},
	}

	return &CommandRegistryEntry{
		"Prints list of issues for given search criteria",
		func() error {
			return jc.CmdList(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdListUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdListUsage(cmd *kingpin.CmdClause, opts *ListOptions) error {
	log.Debugf("Configs: %#v", opts)
	jc.LoadConfigs(cmd, opts)
	log.Debugf("Configs: %#v", opts)
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("assignee", "User assigned the issue").Short('a').StringVar(&opts.Assignee)
	cmd.Flag("component", "Component to search for").Short('c').StringVar(&opts.Component)
	cmd.Flag("issuetype", "Issue type to search for").Short('i').StringVar(&opts.IssueType)
	cmd.Flag("limit", "Maximum number of results to return in search").Short('l').IntVar(&opts.MaxResults)
	cmd.Flag("project", "Project to search for").Short('p').StringVar(&opts.Project)
	cmd.Flag("query", "Jira Query Language (JQL) expression for the search").Short('q').StringVar(&opts.Query)
	cmd.Flag("queryfields", "Fields that are used in \"list\" template").Short('f').StringVar(&opts.QueryFields)
	cmd.Flag("reporter", "Reporter to search for").Short('r').StringVar(&opts.Reporter)
	cmd.Flag("sort", "Sort order to return").Short('s').StringVar(&opts.Sort)
	cmd.Flag("watcher", "Watcher to search for").Short('w').StringVar(&opts.Watcher)
	return nil
}

// List will query jira and send data to "list" template
func (jc *JiraCli) CmdList(opts *ListOptions) error {
	log.Debugf("Configs: %#v", opts)
	data, err := jc.Search(opts)
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template.Value, data, nil)
}
