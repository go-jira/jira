package jiracli

import (
	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type TransitionsOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Issue         string
}

func CmdTransitionsRegistry(fig *figtree.FigTree, o *oreo.Client, defaultTemplate string) *CommandRegistryEntry {
	opts := TransitionsOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption(defaultTemplate),
		},
	}

	return &CommandRegistryEntry{
		"List valid issue transitions",
		func() error {
			return CmdTransitions(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdTransitionsUsage(cmd, &opts)
		},
	}
}

func CmdTransitionsUsage(cmd *kingpin.CmdClause, opts *TransitionsOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue to list valid transitions").Required().StringVar(&opts.Issue)
	return nil
}

// Transitions will get issue edit metadata and send to "editmeta" template
func CmdTransitions(o *oreo.Client, opts *TransitionsOptions) error {
	editMeta, err := jira.GetIssueTransitions(o, opts.Endpoint.Value, opts.Issue)
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
