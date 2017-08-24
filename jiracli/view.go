package jiracli

import (
	"github.com/coryb/figtree"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ViewOptions struct {
	GlobalOptions
	jira.IssueOptions
	Issue string
}

func (jc *JiraCli) CmdViewRegistry() *CommandRegistryEntry {
	opts := ViewOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("view"),
		},
	}

	return &CommandRegistryEntry{
		"Prints issue details",
		func() error {
			return jc.CmdView(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdViewUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdViewUsage(cmd *kingpin.CmdClause, opts *ViewOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("expand", "field to expand for the issue").StringsVar(&opts.Expand)
	cmd.Flag("field", "field to return for the issue").StringsVar(&opts.Fields)
	cmd.Flag("property", "property to return for issue").StringsVar(&opts.Properties)
	cmd.Arg("ISSUE", "issue id to view").Required().StringVar(&opts.Issue)
	return nil
}

// View will get issue data and send to "view" template
func (jc *JiraCli) CmdView(opts *ViewOptions) error {
	data, err := jc.GetIssue(opts.Issue, opts)
	if err != nil {
		return err
	}
	if err := jc.runTemplate(opts.Template.Value, data, nil); err != nil {
		return err
	}
	if opts.Browse.Value {
		return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
