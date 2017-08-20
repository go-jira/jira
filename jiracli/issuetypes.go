package jiracli

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type IssueTypesOptions struct {
	GlobalOptions
	Project string
}

func (jc *JiraCli) CmdIssueTypesRegistry() *CommandRegistryEntry {
	opts := IssueTypesOptions{
		GlobalOptions: GlobalOptions{
			Template: "issuetypes",
		},
	}

	return &CommandRegistryEntry{
		"Show issue types for a project",
		func() error {
			return jc.CmdIssueTypes(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdIssueTypesUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdIssueTypesUsage(cmd *kingpin.CmdClause, opts *IssueTypesOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("project", "project to list issueTypes").Short('p').StringVar(&opts.Project)

	return nil
}

// CmdIssueTypes will get available issueTypes for project and send to the  "issueTypes" template
func (jc *JiraCli) CmdIssueTypes(opts *IssueTypesOptions) error {
	if opts.Project == "" {
		return fmt.Errorf("Project Required.")
	}
	data, err := jc.GetIssueCreateMetaProject(opts.Project)
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template, data, nil)
}
