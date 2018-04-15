package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EpicAddOptions struct {
	jiradata.EpicIssues `yaml:",inline" json:",inline" figtree:",inline"`
	Epic                string `yaml:"epic,omitempty" json:"epic,omitempty"`
}

func CmdEpicAddRegistry() *jiracli.CommandRegistryEntry {
	opts := EpicAddOptions{}

	return &jiracli.CommandRegistryEntry{
		"Add issues to Epic",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdEpicAddUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdEpicAdd(o, globals, &opts)
		},
	}
}

func CmdEpicAddUsage(cmd *kingpin.CmdClause, opts *EpicAddOptions) error {
	cmd.Arg("EPIC", "Epic Key or ID to add issues to").Required().StringVar(&opts.Epic)
	cmd.Arg("ISSUE", "Issues to add to epic").Required().StringsVar(&opts.Issues)
	return nil
}

func CmdEpicAdd(o *oreo.Client, globals *jiracli.GlobalOptions, opts *EpicAddOptions) error {
	if err := jira.EpicAddIssues(o, globals.Endpoint.Value, opts.Epic, &opts.EpicIssues); err != nil {
		return err
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", opts.Epic, jira.URLJoin(globals.Endpoint.Value, "browse", opts.Epic))
		for _, issue := range opts.Issues {
			fmt.Printf("OK %s %s\n", issue, jira.URLJoin(globals.Endpoint.Value, "browse", issue))
		}
	}

	return nil
}
