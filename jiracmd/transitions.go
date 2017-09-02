package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type TransitionsOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue         string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdTransitionsRegistry(fig *figtree.FigTree, o *oreo.Client, defaultTemplate string) *jiracli.CommandRegistryEntry {
	opts := TransitionsOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption(defaultTemplate),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"List valid issue transitions",
		func() error {
			return CmdTransitions(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdTransitionsUsage(cmd, &opts)
		},
	}
}

func CmdTransitionsUsage(cmd *kingpin.CmdClause, opts *TransitionsOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue to list valid transitions").Required().StringVar(&opts.Issue)
	return nil
}

// Transitions will get issue edit metadata and send to "editmeta" template
func CmdTransitions(o *oreo.Client, opts *TransitionsOptions) error {
	editMeta, err := jira.GetIssueTransitions(o, opts.Endpoint.Value, opts.Issue)
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
