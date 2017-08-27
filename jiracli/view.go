package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ViewOptions struct {
	GlobalOptions     `yaml:",inline" figtree:",inline"`
	jira.IssueOptions `yaml:",inline" figtree:",inline"`
	Issue             string
}

func CmdViewRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := ViewOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("view"),
		},
	}

	return &CommandRegistryEntry{
		"Prints issue details",
		func() error {
			return CmdView(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdViewUsage(cmd, &opts)
		},
	}
}

func CmdViewUsage(cmd *kingpin.CmdClause, opts *ViewOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
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
	if err := runTemplate(opts.Template.Value, data, nil); err != nil {
		return err
	}
	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
