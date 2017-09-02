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

type CommentOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Overrides     map[string]string `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	Issue         string            `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdCommentRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := CommentOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("comment"),
		},
		Overrides: map[string]string{},
	}

	return &jiracli.CommandRegistryEntry{
		"Add comment to issue",
		func() error {
			return CmdComment(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdCommentUsage(cmd, &opts)
		},
	}
}

func CmdCommentUsage(cmd *kingpin.CmdClause, opts *CommentOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	jiracli.EditorUsage(cmd, &opts.GlobalOptions)
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = jiracli.FlagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Arg("ISSUE", "issue id to update").StringVar(&opts.Issue)
	return nil
}

// CmdComment will update issue with comment
func CmdComment(o *oreo.Client, opts *CommentOptions) error {
	comment := jiradata.Comment{}
	input := struct {
		Overrides map[string]string
	}{
		opts.Overrides,
	}
	err := jiracli.EditLoop(&opts.GlobalOptions, &input, &comment, func() error {
		_, err := jira.IssueAddComment(o, opts.Endpoint.Value, opts.Issue, &comment)
		return err
	})
	if err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, opts.Endpoint.Value, opts.Issue)

	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}

	return nil
}
