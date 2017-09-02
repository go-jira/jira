package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WorklogListOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue         string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdWorklogListRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := WorklogListOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("worklogs"),
		},
	}
	return &jiracli.CommandRegistryEntry{
		"Prints the worklog data for given issue",
		func() error {
			return CmdWorklogList(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdWorklogListUsage(cmd, &opts)
		},
	}
}

func CmdWorklogListUsage(cmd *kingpin.CmdClause, opts *WorklogListOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue id to fetch worklogs").Required().StringVar(&opts.Issue)
	return nil
}

// // CmdWorklogList will get worklog data for given issue and sent to the "worklogs" template
func CmdWorklogList(o *oreo.Client, opts *WorklogListOptions) error {
	data, err := jira.GetIssueWorklog(o, opts.Endpoint.Value, opts.Issue)
	if err != nil {
		return err
	}
	if err := jiracli.RunTemplate(opts.Template.Value, data, nil); err != nil {
		return err
	}
	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
