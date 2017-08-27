package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdFieldsRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := GlobalOptions{
		Template: figtree.NewStringOption("fields"),
	}
	return &CommandRegistryEntry{
		"Prints all fields, both System and Custom",
		func() error {
			return CmdFields(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			err := GlobalUsage(cmd, &opts)
			TemplateUsage(cmd, &opts)
			return err
		},
	}
}

// Fields will send data from /rest/api/2/field API to "fields" template
func CmdFields(o *oreo.Client, opts *GlobalOptions) error {
	data, err := jira.GetFields(o, opts.Endpoint.Value)
	if err != nil {
		return err
	}
	return runTemplate(opts.Template.Value, data, nil)
}
