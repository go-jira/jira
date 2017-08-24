package jiracli

import (
	"github.com/coryb/figtree"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CreateMetaOptions struct {
	GlobalOptions
	Project   string `yaml:"project,omitempty" json:"project,omitempty"`
	IssueType string `yaml:"issuetype,omitempty" json:"issuetype,omitempty"`
}

func (jc *JiraCli) CmdCreateMetaRegistry() *CommandRegistryEntry {
	opts := CreateMetaOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("createmeta"),
		},
	}

	return &CommandRegistryEntry{
		"View 'create' metadata",
		func() error {
			return jc.CmdCreateMeta(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdCreateMetaUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdCreateMetaUsage(cmd *kingpin.CmdClause, opts *CreateMetaOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("project", "project to fetch create metadata").Short('p').StringVar(&opts.Project)
	cmd.Flag("issuetype", "issuetype in project to fetch create metadata").Short('i').StringVar(&opts.IssueType)
	return nil
}

// Create will get issue create metadata and send to "createmeta" template
func (jc *JiraCli) CmdCreateMeta(opts *CreateMetaOptions) error {
	if err := jc.defaultIssueType(&opts.Project, &opts.IssueType); err != nil {
		return err
	}
	createMeta, err := jc.GetIssueCreateMetaIssueType(opts.Project, opts.IssueType)
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template.Value, createMeta, nil)
}
