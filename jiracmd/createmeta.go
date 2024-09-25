package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CreateMetaOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
	IssueType             string `yaml:"issuetype,omitempty" json:"issuetype,omitempty"`
}

func CmdCreateMetaRegistry() *jiracli.CommandRegistryEntry {
	opts := CreateMetaOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("createmeta"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"View 'create' metadata",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdCreateMetaUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdCreateMeta(o, globals, &opts)
		},
	}
}

func CmdCreateMetaUsage(cmd *kingpin.CmdClause, opts *CreateMetaOptions) error {
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	jiracli.GJsonQueryUsage(cmd, &opts.CommonOptions)
	cmd.Flag("project", "project to fetch create metadata").Short('p').StringVar(&opts.Project)
	cmd.Flag("issuetype", "issuetype in project to fetch create metadata").Short('i').StringVar(&opts.IssueType)
	return nil
}

// Create will get issue create metadata and send to "createmeta" template
func CmdCreateMeta(o *oreo.Client, globals *jiracli.GlobalOptions, opts *CreateMetaOptions) error {
	if err := defaultIssueType(o, globals.Endpoint.Value, &opts.Project, &opts.IssueType); err != nil {
		return err
	}
	createMeta, err := jira.GetIssueCreateMetaIssueType(o, globals.Endpoint.Value, opts.Project, opts.IssueType)
	if err != nil {
		return err
	}
	return opts.PrintTemplate(createMeta)
}
