package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ComponentsOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project       string `yaml:"project,omitempty" json:"project,omitempty"`
}

func CmdComponentsRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := ComponentsOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("components"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Show components for a project",
		func() error {
			return CmdComponents(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdComponentsUsage(cmd, &opts)
		},
	}
}

func CmdComponentsUsage(cmd *kingpin.CmdClause, opts *ComponentsOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
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
	return jiracli.RunTemplate(opts.Template.Value, data, nil)
}
