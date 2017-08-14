package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type IssueLinkOptions struct {
	GlobalOptions
	jiradata.LinkIssueRequest
	LinkType string
}

func (jc *JiraCli) CmdIssueLinkRegistry() *CommandRegistryEntry {
	opts := IssueLinkOptions{
		GlobalOptions: GlobalOptions{
			Template: "edit",
		},
		LinkIssueRequest: jiradata.LinkIssueRequest{
			Type:         &jiradata.IssueLinkType{},
			InwardIssue:  &jiradata.IssueRef{},
			OutwardIssue: &jiradata.IssueRef{},
		},
	}
	return &CommandRegistryEntry{
		"Link two issues",
		func() error {
			return jc.CmdIssueLink(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdIssueLinkUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdIssueLinkUsage(cmd *kingpin.CmdClause, opts *IssueLinkOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("comment", "Comment message when linking issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Comment = &jiradata.Comment{
			Body: flagValue(ctx, "comment"),
		}
		return nil
	}).String()
	cmd.Arg("OUTWARDISSUE", "outward issue").Required().StringVar(&opts.OutwardIssue.Key)
	cmd.Arg("ISSUELINKTYPE", "issue link type").Required().StringVar(&opts.Type.Name)
	cmd.Arg("INWARDISSUE", "inward issue").Required().StringVar(&opts.InwardIssue.Key)
	return nil
}

// CmdBlock will update the given issue as being a duplicate by the given dup issue
// and will attempt to resolve the dup issue
func (jc *JiraCli) CmdIssueLink(opts *IssueLinkOptions) error {
	if err := jc.LinkIssues(&opts.LinkIssueRequest); err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.InwardIssue.Key, jc.Endpoint, opts.InwardIssue.Key)
	fmt.Printf("OK %s %s/browse/%s\n", opts.OutwardIssue.Key, jc.Endpoint, opts.OutwardIssue.Key)

	// FIXME implement browse

	return nil
}
