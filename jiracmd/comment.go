package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CommentOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string            `yaml:"project,omitempty" json:"project,omitempty"`
	Overrides             map[string]string `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	Issue                 string            `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdCommentRegistry() *jiracli.CommandRegistryEntry {
	opts := CommentOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("comment"),
		},
		Overrides: map[string]string{},
	}

	return &jiracli.CommandRegistryEntry{
		"Add comment to issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdCommentUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdComment(o, globals, &opts)
		},
	}
}

func CmdCommentUsage(cmd *kingpin.CmdClause, opts *CommentOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = jiracli.FlagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Arg("ISSUE", "issue id to update").StringVar(&opts.Issue)
	return nil
}

// CmdComment will update issue with comment
func CmdComment(o *oreo.Client, globals *jiracli.GlobalOptions, opts *CommentOptions) error {
	comment := jiradata.Comment{}
	input := struct {
		Overrides map[string]string `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	}{
		opts.Overrides,
	}
	err := jiracli.EditLoop(&opts.CommonOptions, &input, &comment, func() error {
		_, err := jira.IssueAddComment(o, globals.Endpoint.Value, opts.Issue, &comment)
		return err
	})
	if err != nil {
		return err
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", opts.Issue, jira.URLJoin(globals.Endpoint.Value, "browse", opts.Issue))
	}

	if opts.Browse.Value {
		return CmdBrowse(globals, opts.Issue)
	}

	return nil
}
