package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type IssueTypesOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Project       string
}

func CmdIssueTypesRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := IssueTypesOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("issuetypes"),
		},
	}

	return &CommandRegistryEntry{
		"Show issue types for a project",
		func() error {
			return CmdIssueTypes(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdIssueTypesUsage(cmd, &opts)
		},
	}
}

func CmdIssueTypesUsage(cmd *kingpin.CmdClause, opts *IssueTypesOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("project", "project to list issueTypes").Short('p').StringVar(&opts.Project)

	return nil
}

// CmdIssueTypes will get available issueTypes for project and send to the  "issueTypes" template
func CmdIssueTypes(o *oreo.Client, opts *IssueTypesOptions) error {
	if opts.Project == "" {
		return fmt.Errorf("Project Required.")
	}
	data, err := jira.GetIssueCreateMetaProject(o, opts.Endpoint.Value, opts.Project)
	if err != nil {
		return err
	}
	return runTemplate(opts.Template.Value, data, nil)
}
