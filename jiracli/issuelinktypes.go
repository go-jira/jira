package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdIssueLinkTypesRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := GlobalOptions{
		Template: figtree.NewStringOption("issuelinktypes"),
	}

	return &CommandRegistryEntry{
		"Show the issue link types",
		func() error {
			return CmdIssueLinkTypes(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdIssueLinkTypesUsage(cmd, &opts)
		},
	}
}

func CmdIssueLinkTypesUsage(cmd *kingpin.CmdClause, opts *GlobalOptions) error {
	if err := GlobalUsage(cmd, opts); err != nil {
		return err
	}
	TemplateUsage(cmd, opts)
	return nil
}

// CmdIssueLinkTypes will get issue link type data and send to "issuelinktypes" template
func CmdIssueLinkTypes(o *oreo.Client, opts *GlobalOptions) error {
	data, err := jira.GetIssueLinkTypes(o, opts.Endpoint.Value)
	if err != nil {
		return err
	}
	return runTemplate(opts.Template.Value, data, nil)
}
