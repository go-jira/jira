package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EpicRemoveOptions struct {
	jiradata.EpicIssues `yaml:",inline" json:",inline" figtree:",inline"`
	Project             string `yaml:"project,omitempty" json:"project,omitempty"`
}

func CmdEpicRemoveRegistry() *jiracli.CommandRegistryEntry {
	opts := EpicRemoveOptions{}

	return &jiracli.CommandRegistryEntry{
		"Remove issues from Epic",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdEpicRemoveUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			for i := range opts.Issues {
				opts.Issues[i] = jiracli.FormatIssue(opts.Issues[i], opts.Project)
			}
			return CmdEpicRemove(o, globals, &opts)
		},
	}
}

func CmdEpicRemoveUsage(cmd *kingpin.CmdClause, opts *EpicRemoveOptions) error {
	cmd.Arg("ISSUE", "Issues to remove from any epic").Required().StringsVar(&opts.Issues)
	return nil
}

func CmdEpicRemove(o *oreo.Client, globals *jiracli.GlobalOptions, opts *EpicRemoveOptions) error {
	if err := jira.EpicRemoveIssues(o, globals.Endpoint.Value, &opts.EpicIssues); err != nil {
		return err
	}

	if !globals.Quiet.Value {
		for _, issue := range opts.Issues {
			fmt.Printf("OK %s %s\n", issue, jira.URLJoin(globals.Endpoint.Value, "browse", issue))
		}
	}

	return nil
}
