package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ComponentsOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
}

func CmdComponentsRegistry() *jiracli.CommandRegistryEntry {
	opts := ComponentsOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("components"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Show components for a project",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdComponentsUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdComponents(o, globals, &opts)
		},
	}
}

func CmdComponentsUsage(cmd *kingpin.CmdClause, opts *ComponentsOptions) error {
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	jiracli.GJsonQueryUsage(cmd, &opts.CommonOptions)
	cmd.Flag("project", "project to list components").Short('p').StringVar(&opts.Project)

	return nil
}

// CmdComponents will get available components for project and send to the  "components" template
func CmdComponents(o *oreo.Client, globals *jiracli.GlobalOptions, opts *ComponentsOptions) error {
	if opts.Project == "" {
		return fmt.Errorf("Project Required.")
	}
	data, err := jira.GetProjectComponents(o, globals.Endpoint.Value, opts.Project)
	if err != nil {
		return err
	}
	return opts.PrintTemplate(data)
}
