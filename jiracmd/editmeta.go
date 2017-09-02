package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EditMetaOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue         string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdEditMetaRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {

	opts := EditMetaOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("editmeta"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"View 'edit' metadata",
		func() error {
			return CmdEditMeta(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdEditMetaUsage(cmd, &opts)
		},
	}
}

func CmdEditMetaUsage(cmd *kingpin.CmdClause, opts *EditMetaOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "edit metadata for issue id").Required().StringVar(&opts.Issue)
	return nil
}

// EditMeta will get issue edit metadata and send to "editmeta" template
func CmdEditMeta(o *oreo.Client, opts *EditMetaOptions) error {
	editMeta, err := jira.GetIssueEditMeta(o, opts.Endpoint.Value, opts.Issue)
	if err != nil {
		return err
	}
	if err := jiracli.RunTemplate(opts.Template.Value, editMeta, nil); err != nil {
		return err
	}
	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
