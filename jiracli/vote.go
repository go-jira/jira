package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type VoteAction int

const (
	VoteUP VoteAction = iota
	VoteDown
)

type VoteOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Issue         string
	Action        VoteAction
}

func CmdVoteRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := VoteOptions{
		GlobalOptions: GlobalOptions{},
		Action:        VoteUP,
	}

	return &CommandRegistryEntry{
		"Vote up/down an issue",
		func() error {
			return CmdVote(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdVoteUsage(cmd, &opts)
		},
	}
}

func CmdVoteUsage(cmd *kingpin.CmdClause, opts *VoteOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("down", "downvote the issue").Short('d').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Action = VoteDown
		return nil
	}).Bool()
	cmd.Arg("ISSUE", "issue id to vote").StringVar(&opts.Issue)
	return nil
}

// Vote will up/down vote an issue
func CmdVote(o *oreo.Client, opts *VoteOptions) error {
	if opts.Action == VoteUP {
		if err := jira.IssueAddVote(o, opts.Endpoint.Value, opts.Issue); err != nil {
			return err
		}
	} else {
		if err := jira.IssueRemoveVote(o, opts.Endpoint.Value, opts.Issue); err != nil {
			return err
		}
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, opts.Endpoint.Value, opts.Issue)

	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
