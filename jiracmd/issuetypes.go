package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type IssueTypesOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project       string `yaml:"project,omitempty" json:"project,omitempty"`
}

func CmdIssueTypesRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := IssueTypesOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("issuetypes"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Show issue types for a project",
		func() error {
			return CmdIssueTypes(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdIssueTypesUsage(cmd, &opts)
		},
	}
}

func CmdIssueTypesUsage(cmd *kingpin.CmdClause, opts *IssueTypesOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
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
	return jiracli.RunTemplate(opts.Template.Value, data, nil)
}
