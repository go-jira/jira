package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type ComponentAddOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	jiradata.Component    `yaml:",inline" json:",inline" figtree:",inline"`
}

func CmdComponentAddRegistry() *jiracli.CommandRegistryEntry {
	opts := ComponentAddOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("component-add"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Add component",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdComponentAddUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdComponentAdd(o, globals, &opts)
		},
	}
}

func CmdComponentAddUsage(cmd *kingpin.CmdClause, opts *ComponentAddOptions) error {
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("project", "project to create component in").Short('p').StringVar(&opts.Project)
	cmd.Flag("name", "name of component").Short('n').StringVar(&opts.Name)
	cmd.Flag("description", "description of component").Short('d').StringVar(&opts.Description)
	cmd.Flag("lead", "person that acts as lead for component").Short('l').StringVar(&opts.LeadUserName)
	return nil
}

// CmdComponentAdd sends the provided overrides to the "component-add" template for editing, then
// will parse the edited document as YAML and submit the document to jira.
func CmdComponentAdd(o *oreo.Client, globals *jiracli.GlobalOptions, opts *ComponentAddOptions) error {
	var err error
	component := &jiradata.Component{}
	err = jiracli.EditLoop(&opts.CommonOptions, &opts.Component, component, func() error {
		_, err = jira.CreateComponent(o, globals.Endpoint.Value, component)
		return err
	})
	if err != nil {
		return err
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", component.Project, component.Name)
	}
	return nil
}
