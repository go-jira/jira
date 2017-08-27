package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ComponentAddOptions struct {
	GlobalOptions      `yaml:",inline" figtree:",inline"`
	jiradata.Component `yaml:",inline" figtree:",inline"`
}

func CmdComponentAddRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := ComponentAddOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("component-add"),
		},
	}

	return &CommandRegistryEntry{
		"Add component",
		func() error {
			return CmdComponentAdd(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdComponentAddUsage(cmd, &opts)
		},
	}
}

func CmdComponentAddUsage(cmd *kingpin.CmdClause, opts *ComponentAddOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	EditorUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("project", "project to create component in").Short('p').StringVar(&opts.Project)
	cmd.Flag("name", "name of component").Short('n').StringVar(&opts.Name)
	cmd.Flag("description", "description of component").Short('d').StringVar(&opts.Description)
	cmd.Flag("lead", "person that acts as lead for component").Short('l').StringVar(&opts.LeadUserName)
	return nil
}

// CmdComponentAdd sends the provided overrides to the "component-add" template for editing, then
// will parse the edited document as YAML and submit the document to jira.
func CmdComponentAdd(o *oreo.Client, opts *ComponentAddOptions) error {
	var err error
	component := &jiradata.Component{}
	var resp *jiradata.Component
	err = editLoop(&opts.GlobalOptions, &opts.Component, component, func() error {
		resp, err = jira.CreateComponent(o, opts.Endpoint.Value, component)
		return err
	})
	if err != nil {
		return err
	}

	fmt.Printf("OK %s %s\n", component.Project, component.Name)
	return nil
}
