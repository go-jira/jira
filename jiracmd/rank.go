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
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	First         string `yaml:"first,omitempty" json:"first,omitempty"`
	Second        string `yaml:"second,omitempty" json:"second,omitempty"`
	Order         string `yaml:"order,omitempty" json:"order,omitempty"`
}

func CmdRankRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := RankOptions{}

	return &jiracli.CommandRegistryEntry{
		"Mark issues as blocker",
		func() error {
			return CmdRank(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdRankUsage(cmd, &opts)
		},
	}
}

func CmdRankUsage(cmd *kingpin.CmdClause, opts *RankOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("FIRST-ISSUE", "first issue").Required().StringVar(&opts.First)
	cmd.Arg("after|before", "rank ordering").Required().HintOptions("after", "before").EnumVar(&opts.Order, "after", "before")
	cmd.Arg("SECOND-ISSUE", "second issue").Required().StringVar(&opts.Second)
	return nil
}

// CmdRank order two issue
func CmdRank(o *oreo.Client, opts *RankOptions) error {
	req := &jiradata.RankRequest{
		Issues: []string{opts.First},
	}

	if opts.Order == "after" {
		req.RankAfterIssue = opts.Second
	} else {
		req.RankBeforeIssue = opts.Second
	}

	if err := jira.RankIssues(o, opts.Endpoint.Value, req); err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.First, opts.Endpoint.Value, opts.First)
	fmt.Printf("OK %s %s/browse/%s\n", opts.Second, opts.Endpoint.Value, opts.Second)

	if opts.Browse.Value {
		if err := CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.First}); err != nil {
			return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Second})
		}
	}

	return nil
}
