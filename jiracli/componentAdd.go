package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ComponentAddOptions struct {
	GlobalOptions
	jiradata.Component
}

func (jc *JiraCli) CmdComponentAddRegistry() *CommandRegistryEntry {
	opts := ComponentAddOptions{
		GlobalOptions: GlobalOptions{
			Template: "component-add",
		},
	}

	return &CommandRegistryEntry{
		"Add component",
		func() error {
			return jc.CmdComponentAdd(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdComponentAddUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdComponentAddUsage(cmd *kingpin.CmdClause, opts *ComponentAddOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").BoolVar(&opts.SkipEditing)
	cmd.Flag("project", "project to create component in").Short('p').StringVar(&opts.Project)
	cmd.Flag("name", "name of component").Short('n').StringVar(&opts.Name)
	cmd.Flag("description", "description of component").Short('d').StringVar(&opts.Description)
	cmd.Flag("lead", "person that acts as lead for component").Short('l').StringVar(&opts.LeadUserName)
	return nil
}

// CmdComponentAdd sends the provided overrides to the "component-add" template for editing, then
// will parse the edited document as YAML and submit the document to jira.
func (jc *JiraCli) CmdComponentAdd(opts *ComponentAddOptions) error {
	var err error
	component := &jiradata.Component{}
	var resp *jiradata.Component
	err = jc.editLoop(&opts.GlobalOptions, &opts.Component, component, func() error {
		resp, err = jc.CreateComponent(component)
		return err
	})
	if err != nil {
		return err
	}

	fmt.Printf("OK %s %s\n", component.Project, component.Name)
	return nil
}
