package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdFieldsRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := jiracli.GlobalOptions{
		Template: figtree.NewStringOption("fields"),
	}
	return &jiracli.CommandRegistryEntry{
		"Prints all fields, both System and Custom",
		func() error {
			return CmdFields(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			err := jiracli.GlobalUsage(cmd, &opts)
			jiracli.TemplateUsage(cmd, &opts)
			return err
		},
	}
}

// Fields will send data from /rest/api/2/field API to "fields" template
func CmdFields(o *oreo.Client, opts *jiracli.GlobalOptions) error {
	data, err := jira.GetFields(o, opts.Endpoint.Value)
	if err != nil {
		return err
	}
	return jiracli.RunTemplate(opts.Template.Value, data, nil)
}
