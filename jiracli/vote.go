package jiracli

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type VoteAction int

const (
	VoteUP VoteAction = iota
	VoteDown
)

type VoteOptions struct {
	GlobalOptions
	Issue  string
	Action VoteAction
}

func (jc *JiraCli) CmdVoteRegistry() *CommandRegistryEntry {
	opts := VoteOptions{
		GlobalOptions: GlobalOptions{},
		Action:        VoteUP,
	}

	return &CommandRegistryEntry{
		"Vote up/down an issue",
		func() error {
			return jc.CmdVote(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdVoteUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdVoteUsage(cmd *kingpin.CmdClause, opts *VoteOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	cmd.Flag("down", "downvote the issue").Short('d').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Action = VoteDown
		return nil
	}).Bool()
	cmd.Arg("ISSUE", "issue id to vote").StringVar(&opts.Issue)
	return nil
}

// Vote will up/down vote an issue
func (jc *JiraCli) CmdVote(opts *VoteOptions) error {
	if opts.Action == VoteUP {
		if err := jc.IssueAddVote(opts.Issue); err != nil {
			return err
		}
	} else {
		if err := jc.IssueRemoveVote(opts.Issue); err != nil {
			return err
		}
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)

	// FIXME implement browse
	return nil
}
