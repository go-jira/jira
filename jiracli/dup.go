package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type DupOptions struct {
	GlobalOptions
	jiradata.LinkIssueRequest
	Duplicate string
	Issue     string
}

func (jc *JiraCli) CmdDupRegistry() *CommandRegistryEntry {
	opts := DupOptions{
		GlobalOptions: GlobalOptions{
			Template: "edit",
		},
	}

	return &CommandRegistryEntry{
		"Mark issues as duplicate",
		func() error {
			return jc.CmdDup(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdDupUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdDupUsage(cmd *kingpin.CmdClause, opts *DupOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("comment", "Comment message when marking issue as duplicate").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Comment = &jiradata.Comment{
			Body: flagValue(ctx, "comment"),
		}
		return nil
	}).String()
	cmd.Arg("DUPLICATE", "duplicate issue to mark closed").Required().StringVar(&opts.Duplicate)
	cmd.Arg("ISSUE", "duplicate issue to leave open").Required().StringVar(&opts.Issue)
	return nil
}

// CmdDups will update the given issue as being a duplicate by the given dup issue
// and will attempt to resolve the dup issue
func (jc *JiraCli) CmdDup(opts *DupOptions) error {
	opts.Type = &jiradata.IssueLinkType{
		// FIXME is this consitent across multiple jira installs?
		Name: "Duplicate",
	}
	opts.InwardIssue = &jiradata.IssueRef{
		Key: opts.Duplicate,
	}
	opts.OutwardIssue = &jiradata.IssueRef{
		Key: opts.Issue,
	}

	if err := jc.LinkIssues(&opts.LinkIssueRequest); err != nil {
		return err
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)

	meta, err := jc.GetIssueTransitions(opts.Duplicate)
	if err != nil {
		return err
	}
	for _, trans := range []string{"close", "done", "start", "stop"} {
		transMeta := meta.Transitions.Find(trans)
		if transMeta != nil {
			issueUpdate := jiradata.IssueUpdate{
				Transition: transMeta,
			}
			if err = jc.TransitionIssue(opts.Duplicate, &issueUpdate); err != nil {
				return err
			}
			// if we just started the issue now we need to stop it
			if trans != "start" {
				break
			}
		}
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.Duplicate, jc.Endpoint, opts.Duplicate)

	// FIXME implement browse

	return nil
}
