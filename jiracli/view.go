package jiracli

import (
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ViewOptions struct {
	GlobalOptions
	jira.IssueOptions
}

func (jc *JiraCli) CmdViewRegistry() *CommandRegistryEntry {
	issue := ""
	opts := ViewOptions{
		GlobalOptions: GlobalOptions{
			Template: "view",
		},
	}

	return &CommandRegistryEntry{
		"Prints issue details",
		func() error {
			return jc.CmdView(issue, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdViewUsage(cmd, &issue, &opts)
		},
	}
}

func (jc *JiraCli) CmdViewUsage(cmd *kingpin.CmdClause, issue *string, opts *ViewOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("expand", "field to expand for the issue").StringsVar(&opts.Expand)
	cmd.Flag("field", "field to return for the issue").StringsVar(&opts.Fields)
	cmd.Flag("property", "property to return for issue").StringsVar(&opts.Properties)
	cmd.Arg("ISSUE", "issue id to view").Required().StringVar(issue)
	return nil
}

// View will get issue data and send to "view" template
func (jc *JiraCli) CmdView(issue string, opts *ViewOptions) error {
	data, err := jc.GetIssue(issue, opts)
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template, data, nil)
}
