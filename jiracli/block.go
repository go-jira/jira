package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type BlockOptions struct {
	GlobalOptions
	jiradata.LinkIssueRequest
	Blocker string
	Issue   string
}

func (jc *JiraCli) CmdBlockRegistry() *CommandRegistryEntry {
	opts := BlockOptions{
		GlobalOptions: GlobalOptions{
			Template: "edit",
		},
	}

	return &CommandRegistryEntry{
		"Mark issues as blocker",
		func() error {
			return jc.CmdBlock(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdBlockUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdBlockUsage(cmd *kingpin.CmdClause, opts *BlockOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("comment", "Comment message when marking issue as blocker").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Comment = &jiradata.Comment{
			Body: flagValue(ctx, "comment"),
		}
		return nil
	}).String()
	cmd.Arg("BLOCKER", "blocker issue").Required().StringVar(&opts.Blocker)
	cmd.Arg("ISSUE", "issue that is blocked").Required().StringVar(&opts.Issue)
	return nil
}

// CmdBlock will update the given issue as being a duplicate by the given dup issue
// and will attempt to resolve the dup issue
func (jc *JiraCli) CmdBlock(opts *BlockOptions) error {
	opts.Type = &jiradata.IssueLinkType{
		// FIXME is this consitent across multiple jira installs?
		Name: "Blocks",
	}
	opts.InwardIssue = &jiradata.IssueRef{
		Key: opts.Issue,
	}
	opts.OutwardIssue = &jiradata.IssueRef{
		Key: opts.Blocker,
	}

	if err := jc.LinkIssues(&opts.LinkIssueRequest); err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)
	fmt.Printf("OK %s %s/browse/%s\n", opts.Blocker, jc.Endpoint, opts.Blocker)

	// FIXME implement browse

	return nil
}
