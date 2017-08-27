package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EditMetaOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Issue         string
}

func CmdEditMetaRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {

	opts := EditMetaOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("editmeta"),
		},
	}

	return &CommandRegistryEntry{
		"View 'edit' metadata",
		func() error {
			return CmdEditMeta(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdEditMetaUsage(cmd, &opts)
		},
	}
}

func CmdEditMetaUsage(cmd *kingpin.CmdClause, opts *EditMetaOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "edit metadata for issue id").Required().StringVar(&opts.Issue)
	return nil
}

// EditMeta will get issue edit metadata and send to "editmeta" template
func CmdEditMeta(o *oreo.Client, opts *EditMetaOptions) error {
	editMeta, err := jira.GetIssueEditMeta(o, opts.Endpoint.Value, opts.Issue)
	if err != nil {
		return err
	}
	if err := runTemplate(opts.Template.Value, editMeta, nil); err != nil {
		return err
	}
	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
