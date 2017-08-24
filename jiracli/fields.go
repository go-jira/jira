package jiracli

import (
	"github.com/coryb/figtree"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func (jc *JiraCli) CmdFieldsRegistry() *CommandRegistryEntry {
	opts := GlobalOptions{
		Template: figtree.NewStringOption("fields"),
	}
	return &CommandRegistryEntry{
		"Prints all fields, both System and Custom",
		func() error {
			return jc.CmdFields(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			err := jc.GlobalUsage(cmd, &opts)
			jc.TemplateUsage(cmd, &opts)
			return err
		},
	}
}

// Fields will send data from /rest/api/2/field API to "fields" template
func (jc *JiraCli) CmdFields(opts *GlobalOptions) error {
	data, err := jc.GetFields()
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template.Value, data, nil)
}
