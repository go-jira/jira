package jiracli

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ComponentsOptions struct {
	GlobalOptions
	Project string
}

func (jc *JiraCli) CmdComponentsRegistry() *CommandRegistryEntry {
	opts := ComponentsOptions{
		GlobalOptions: GlobalOptions{
			Template: "components",
		},
	}

	return &CommandRegistryEntry{
		"Show components for a project",
		func() error {
			return jc.CmdComponents(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdComponentsUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdComponentsUsage(cmd *kingpin.CmdClause, opts *ComponentsOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("project", "project to list components").Short('p').StringVar(&opts.Project)

	return nil
}

// CmdComponents will get available components for project and send to the  "components" template
func (jc *JiraCli) CmdComponents(opts *ComponentsOptions) error {
	if opts.Project == "" {
		return fmt.Errorf("Project Required.")
	}
	data, err := jc.GetProjectComponents(opts.Project)
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template, data, nil)
}
