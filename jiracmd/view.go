package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ViewOptions struct {
	jiracli.GlobalOptions     `yaml:",inline" json:",inline" figtree:",inline"`
	jira.IssueOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue             string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdViewRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := ViewOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("view"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Prints issue details",
		func() error {
			return CmdView(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdViewUsage(cmd, &opts)
		},
	}
}

func CmdViewUsage(cmd *kingpin.CmdClause, opts *ViewOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("expand", "field to expand for the issue").StringsVar(&opts.Expand)
	cmd.Flag("field", "field to return for the issue").StringsVar(&opts.Fields)
	cmd.Flag("property", "property to return for issue").StringsVar(&opts.Properties)
	cmd.Arg("ISSUE", "issue id to view").Required().StringVar(&opts.Issue)
	return nil
}

// View will get issue data and send to "view" template
func CmdView(o *oreo.Client, opts *ViewOptions) error {
	data, err := jira.GetIssue(o, opts.Endpoint.Value, opts.Issue, opts)
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
