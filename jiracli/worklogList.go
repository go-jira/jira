package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WorklogListOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Issue         string
}

func CmdWorklogListRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := WorklogListOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("worklogs"),
		},
	}
	return &CommandRegistryEntry{
		"Prints the worklog data for given issue",
		func() error {
			return CmdWorklogList(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdWorklogListUsage(cmd, &opts)
		},
	}
}

func CmdWorklogListUsage(cmd *kingpin.CmdClause, opts *WorklogListOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue id to fetch worklogs").Required().StringVar(&opts.Issue)
	return nil
}

// // CmdWorklogList will get worklog data for given issue and sent to the "worklogs" template
func CmdWorklogList(o *oreo.Client, opts *WorklogListOptions) error {
	data, err := jira.GetIssueWorklog(o, opts.Endpoint.Value, opts.Issue)
	if err != nil {
		return err
	}
	if err := runTemplate(opts.Template.Value, data, nil); err != nil {
		return err
	}
	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
