package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CommentOptions struct {
	GlobalOptions
	Overrides map[string]string
	Issue     string
}

func (jc *JiraCli) CmdCommentRegistry() *CommandRegistryEntry {
	opts := CommentOptions{
		GlobalOptions: GlobalOptions{
			Template: "comment",
		},
		Overrides: map[string]string{},
	}

	return &CommandRegistryEntry{
		"Add comment to issue",
		func() error {
			return jc.CmdComment(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdCommentUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdCommentUsage(cmd *kingpin.CmdClause, opts *CommentOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = flagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Arg("ISSUE", "issue id to update").StringVar(&opts.Issue)
	return nil
}

// CmdComment will update issue with comment
func (jc *JiraCli) CmdComment(opts *CommentOptions) error {
	comment := jiradata.Comment{}
	input := struct {
		Overrides map[string]string
	}{
		opts.Overrides,
	}
	err := jc.editLoop(&opts.GlobalOptions, &input, &comment, func() error {
		_, err := jc.IssueAddComment(opts.Issue, &comment)
		return err
	})
	if err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)

	if opts.Browse {
		return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}

	return nil
}
