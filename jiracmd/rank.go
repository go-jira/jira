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

type RankOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	First                 string `yaml:"first,omitempty" json:"first,omitempty"`
	Second                string `yaml:"second,omitempty" json:"second,omitempty"`
	Order                 string `yaml:"order,omitempty" json:"order,omitempty"`
}

func CmdRankRegistry() *jiracli.CommandRegistryEntry {
	opts := RankOptions{}

	return &jiracli.CommandRegistryEntry{
		"Mark issues as blocker",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdRankUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdRank(o, globals, &opts)
		},
	}
}

func CmdRankUsage(cmd *kingpin.CmdClause, opts *RankOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Arg("FIRST-ISSUE", "first issue").Required().StringVar(&opts.First)
	cmd.Arg("after|before", "rank ordering").Required().HintOptions("after", "before").EnumVar(&opts.Order, "after", "before")
	cmd.Arg("SECOND-ISSUE", "second issue").Required().StringVar(&opts.Second)
	return nil
}

// CmdRank order two issue
func CmdRank(o *oreo.Client, globals *jiracli.GlobalOptions, opts *RankOptions) error {
	req := &jiradata.RankRequest{
		Issues: []string{opts.First},
	}

	if opts.Order == "after" {
		req.RankAfterIssue = opts.Second
	} else {
		req.RankBeforeIssue = opts.Second
	}

	if err := jira.RankIssues(o, globals.Endpoint.Value, req); err != nil {
		return err
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", opts.First, jira.URLJoin(globals.Endpoint.Value, "browse", opts.First))
		fmt.Printf("OK %s %s\n", opts.Second, jira.URLJoin(globals.Endpoint.Value, "browse", opts.Second))
	}

	if opts.Browse.Value {
		if err := CmdBrowse(globals, opts.First); err != nil {
			return CmdBrowse(globals, opts.Second)
		}
	}

	return nil
}
