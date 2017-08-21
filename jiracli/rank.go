package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type RankOptions struct {
	GlobalOptions
	First  string
	Second string
	Order  string
}

func (jc *JiraCli) CmdRankRegistry() *CommandRegistryEntry {
	opts := RankOptions{
		GlobalOptions: GlobalOptions{},
	}

	return &CommandRegistryEntry{
		"Mark issues as blocker",
		func() error {
			return jc.CmdRank(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdRankUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdRankUsage(cmd *kingpin.CmdClause, opts *RankOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("FIRST-ISSUE", "first issue").Required().StringVar(&opts.First)
	cmd.Arg("after|before", "rank ordering").Required().HintOptions("after", "before").EnumVar(&opts.Order, "after", "before")
	cmd.Arg("SECOND-ISSUE", "second issue").Required().StringVar(&opts.Second)
	return nil
}

// CmdRank order two issue
func (jc *JiraCli) CmdRank(opts *RankOptions) error {
	req := &jiradata.RankRequest{
		Issues: []string{opts.First},
	}

	if opts.Order == "after" {
		req.RankAfterIssue = opts.Second
	} else {
		req.RankBeforeIssue = opts.Second
	}

	if err := jc.RankIssues(req); err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.First, jc.Endpoint, opts.First)
	fmt.Printf("OK %s %s/browse/%s\n", opts.Second, jc.Endpoint, opts.Second)

	if opts.Browse {
		if err := jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.First}); err != nil {
			return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Second})
		}
	}

	return nil
}
