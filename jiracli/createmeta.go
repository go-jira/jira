package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CreateMetaOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Project       string `yaml:"project,omitempty" json:"project,omitempty"`
	IssueType     string `yaml:"issuetype,omitempty" json:"issuetype,omitempty"`
}

func CmdCreateMetaRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := CreateMetaOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("createmeta"),
		},
	}

	return &CommandRegistryEntry{
		"View 'create' metadata",
		func() error {
			return CmdCreateMeta(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdCreateMetaUsage(cmd, &opts)
		},
	}
}

func CmdCreateMetaUsage(cmd *kingpin.CmdClause, opts *CreateMetaOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("project", "project to fetch create metadata").Short('p').StringVar(&opts.Project)
	cmd.Flag("issuetype", "issuetype in project to fetch create metadata").Short('i').StringVar(&opts.IssueType)
	return nil
}

// Create will get issue create metadata and send to "createmeta" template
func CmdCreateMeta(o *oreo.Client, opts *CreateMetaOptions) error {
	if err := defaultIssueType(o, opts.Endpoint.Value, &opts.Project, &opts.IssueType); err != nil {
		return err
	}
	createMeta, err := jira.GetIssueCreateMetaIssueType(o, opts.Endpoint.Value, opts.Project, opts.IssueType)
	if err != nil {
		return err
	}
	return runTemplate(opts.Template.Value, createMeta, nil)
}
