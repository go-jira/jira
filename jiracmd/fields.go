package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdFieldsRegistry() *jiracli.CommandRegistryEntry {
	opts := jiracli.CommonOptions{
		Template: figtree.NewStringOption("fields"),
	}
	return &jiracli.CommandRegistryEntry{
		"Prints all fields, both System and Custom",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			jiracli.TemplateUsage(cmd, &opts)
			jiracli.GJsonQueryUsage(cmd, &opts)
			return nil
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdFields(o, globals, &opts)
		},
	}
}

// Fields will send data from /rest/api/2/field API to "fields" template
func CmdFields(o *oreo.Client, globals *jiracli.GlobalOptions, opts *jiracli.CommonOptions) error {
	data, err := jira.GetFields(o, globals.Endpoint.Value)
	if err != nil {
		return err
	}
	return opts.PrintTemplate(data)
}
