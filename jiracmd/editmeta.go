package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EditMetaOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdEditMetaRegistry() *jiracli.CommandRegistryEntry {

	opts := EditMetaOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("editmeta"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"View 'edit' metadata",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdEditMetaUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdEditMeta(o, globals, &opts)
		},
	}
}

func CmdEditMetaUsage(cmd *kingpin.CmdClause, opts *EditMetaOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	jiracli.GJsonQueryUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "edit metadata for issue id").Required().StringVar(&opts.Issue)
	return nil
}

// CmdEditMeta will get issue edit metadata and send to "editmeta" template
func CmdEditMeta(o *oreo.Client, globals *jiracli.GlobalOptions, opts *EditMetaOptions) error {
	editMeta, err := jira.GetIssueEditMeta(o, globals.Endpoint.Value, opts.Issue)
	if err != nil {
		return err
	}
	if err := opts.PrintTemplate(editMeta); err != nil {
		return err
	}
	if opts.Browse.Value {
		return CmdBrowse(globals, opts.Issue)
	}
	return nil
}
