package jiracmd

import (
	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type TransitionsOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdTransitionsRegistry(defaultTemplate string) *jiracli.CommandRegistryEntry {
	opts := TransitionsOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption(defaultTemplate),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"List valid issue transitions",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdTransitionsUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdTransitions(o, globals, &opts)
		},
	}
}

func CmdTransitionsUsage(cmd *kingpin.CmdClause, opts *TransitionsOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	jiracli.GJsonQueryUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "issue to list valid transitions").Required().StringVar(&opts.Issue)
	return nil
}

// Transitions will get issue edit metadata and send to "editmeta" template
func CmdTransitions(o *oreo.Client, globals *jiracli.GlobalOptions, opts *TransitionsOptions) error {
	editMeta, err := jira.GetIssueTransitions(o, globals.Endpoint.Value, opts.Issue)
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
