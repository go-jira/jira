package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type WorklogListOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdWorklogListRegistry() *jiracli.CommandRegistryEntry {
	opts := WorklogListOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("worklogs"),
		},
	}
	return &jiracli.CommandRegistryEntry{
		"Prints the worklog data for given issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdWorklogListUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdWorklogList(o, globals, &opts)
		},
	}
}

func CmdWorklogListUsage(cmd *kingpin.CmdClause, opts *WorklogListOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	jiracli.GJsonQueryUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "issue id to fetch worklogs").Required().StringVar(&opts.Issue)
	return nil
}

// CmdWorklogList will get worklog data for given issue and sent to the "worklogs" template
func CmdWorklogList(o *oreo.Client, globals *jiracli.GlobalOptions, opts *WorklogListOptions) error {
	data, err := jira.GetIssueWorklog(o, globals.Endpoint.Value, opts.Issue)
	if err != nil {
		return err
	}
	if err := opts.PrintTemplate(struct {
		Worklogs *jiradata.Worklogs `json:"worklogs,omitempty" yaml:"worklogs,omitempty"`
	}{data}); err != nil {
		return err
	}
	if opts.Browse.Value {
		return CmdBrowse(globals, opts.Issue)
	}
	return nil
}
