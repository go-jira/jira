package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdIssueLinkTypesRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := jiracli.GlobalOptions{
		Template: figtree.NewStringOption("issuelinktypes"),
	}

	return &jiracli.CommandRegistryEntry{
		"Show the issue link types",
		func() error {
			return CmdIssueLinkTypes(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdIssueLinkTypesUsage(cmd, &opts)
		},
	}
}

func CmdIssueLinkTypesUsage(cmd *kingpin.CmdClause, opts *jiracli.GlobalOptions) error {
	if err := jiracli.GlobalUsage(cmd, opts); err != nil {
		return err
	}
	jiracli.TemplateUsage(cmd, opts)
	return nil
}

// CmdIssueLinkTypes will get issue link type data and send to "issuelinktypes" template
func CmdIssueLinkTypes(o *oreo.Client, opts *jiracli.GlobalOptions) error {
	data, err := jira.GetIssueLinkTypes(o, opts.Endpoint.Value)
	if err != nil {
		return err
	}
	return jiracli.RunTemplate(opts.Template.Value, data, nil)
}
