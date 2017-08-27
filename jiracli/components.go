package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ComponentsOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Project       string
}

func CmdComponentsRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := ComponentsOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("components"),
		},
	}

	return &CommandRegistryEntry{
		"Show components for a project",
		func() error {
			return CmdComponents(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdComponentsUsage(cmd, &opts)
		},
	}
}

func CmdComponentsUsage(cmd *kingpin.CmdClause, opts *ComponentsOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("project", "project to list components").Short('p').StringVar(&opts.Project)

	return nil
}

// CmdComponents will get available components for project and send to the  "components" template
func CmdComponents(o *oreo.Client, opts *ComponentsOptions) error {
	if opts.Project == "" {
		return fmt.Errorf("Project Required.")
	}
	data, err := jira.GetProjectComponents(o, opts.Endpoint.Value, opts.Project)
	if err != nil {
		return err
	}
	return runTemplate(opts.Template.Value, data, nil)
}
