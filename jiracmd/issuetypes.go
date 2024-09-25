package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type IssueTypesOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
}

func CmdIssueTypesRegistry() *jiracli.CommandRegistryEntry {
	opts := IssueTypesOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("issuetypes"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Show issue types for a project",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdIssueTypesUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdIssueTypes(o, globals, &opts)
		},
	}
}

func CmdIssueTypesUsage(cmd *kingpin.CmdClause, opts *IssueTypesOptions) error {
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	jiracli.GJsonQueryUsage(cmd, &opts.CommonOptions)
	cmd.Flag("project", "project to list issueTypes").Short('p').StringVar(&opts.Project)

	return nil
}

// CmdIssueTypes will get available issueTypes for project and send to the  "issueTypes" template
func CmdIssueTypes(o *oreo.Client, globals *jiracli.GlobalOptions, opts *IssueTypesOptions) error {
	if opts.Project == "" {
		return fmt.Errorf("Project Required.")
	}
	data, err := jira.GetIssueCreateMetaProject(o, globals.Endpoint.Value, opts.Project)
	if err != nil {
		return err
	}
	return opts.PrintTemplate(data)
}
