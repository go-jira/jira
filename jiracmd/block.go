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

type BlockOptions struct {
	jiracli.GlobalOptions             `yaml:",inline" json:",inline" figtree:",inline"`
	jiradata.LinkIssueRequest `yaml:",inline" json:",inline" figtree:",inline"`
}

func CmdBlockRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := BlockOptions{
		GlobalOptions: jiracli.GlobalOptions{
			Template: figtree.NewStringOption("edit"),
		},
		LinkIssueRequest: jiradata.LinkIssueRequest{
			Type: &jiradata.IssueLinkType{
				Name: "Blocks",
			},
			InwardIssue:  &jiradata.IssueRef{},
			OutwardIssue: &jiradata.IssueRef{},
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Mark issues as blocker",
		func() error {
			return CmdBlock(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdBlockUsage(cmd, &opts)
		},
	}
}

func CmdBlockUsage(cmd *kingpin.CmdClause, opts *BlockOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	jiracli.EditorUsage(cmd, &opts.GlobalOptions)
	jiracli.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("comment", "Comment message when marking issue as blocker").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Comment = &jiradata.Comment{
			Body: jiracli.FlagValue(ctx, "comment"),
		}
		return nil
	}).String()
	cmd.Arg("BLOCKER", "blocker issue").Required().StringVar(&opts.OutwardIssue.Key)
	cmd.Arg("ISSUE", "issue that is blocked").Required().StringVar(&opts.InwardIssue.Key)
	return nil
}

// CmdBlock will update the given issue as being a duplicate by the given dup issue
// and will attempt to resolve the dup issue
func CmdBlock(o *oreo.Client, opts *BlockOptions) error {
	if err := jira.LinkIssues(o, opts.Endpoint.Value, &opts.LinkIssueRequest); err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.InwardIssue.Key, opts.Endpoint.Value, opts.InwardIssue.Key)
	fmt.Printf("OK %s %s/browse/%s\n", opts.OutwardIssue.Key, opts.Endpoint.Value, opts.OutwardIssue.Key)

	if opts.Browse.Value {
		if err := CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.InwardIssue.Key}); err != nil {
			return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.OutwardIssue.Key})
		}
	}

	return nil
}
