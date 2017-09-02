package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CreateMetaOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project       string `yaml:"project,omitempty" json:"project,omitempty"`
	IssueType     string `yaml:"issuetype,omitempty" json:"issuetype,omitempty"`
}

func CmdCreateMetaRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := CreateMetaOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("createmeta"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"View 'create' metadata",
		func() error {
			return CmdCreateMeta(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdCreateMetaUsage(cmd, &opts)
		},
	}
}

func CmdCreateMetaUsage(cmd *kingpin.CmdClause, opts *CreateMetaOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
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
	return jiracli.RunTemplate(opts.Template.Value, createMeta, nil)
}
