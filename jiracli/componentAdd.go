package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ComponentAddOptions struct {
	GlobalOptions
	Overrides map[string]string
}

func (jc *JiraCli) CmdComponentAddRegistry() *CommandRegistryEntry {
	opts := ComponentAddOptions{
		GlobalOptions: GlobalOptions{
			Template: "component-add",
		},
		Overrides: map[string]string{},
	}

	return &CommandRegistryEntry{
		"ComponentAdd issue",
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
	cmd.Flag("project", "project to create component in").Short('p').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["project"] = flagValue(ctx, "project")
		return nil
	}).String()
	cmd.Flag("name", "name of component").Short('n').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["name"] = flagValue(ctx, "name")
		return nil
	}).String()
	cmd.Flag("description", "description of component").Short('d').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["description"] = flagValue(ctx, "description")
		return nil
	}).String()
	cmd.Flag("lead", "person that acts as lead for component").Short('l').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["lead"] = flagValue(ctx, "lead")
		return nil
	}).String()
	return nil
}

// CmdComponentAdd sends the provided overrides to the "component-add" template for editing, then
// will parse the edited document as YAML and submit the document to jira.
func (jc *JiraCli) CmdComponentAdd(opts *ComponentAddOptions) error {
	input := struct {
		Overrides map[string]string
	}{
		opts.Overrides,
	}

	var err error
	component := &jiradata.Component{}
	err = jc.editLoop(&opts.GlobalOptions, &input, component, func() error {
		component, err = jc.CreateComponent(component)
		return err
	})
	if err != nil {
		return err
	}

	fmt.Printf("OK %s %s\n", component.Project, component.Name)
	return nil
}
