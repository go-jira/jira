package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type BlockOptions struct {
	GlobalOptions             `yaml:",inline" figtree:",inline"`
	jiradata.LinkIssueRequest `yaml:",inline" figtree:",inline"`
	Blocker                   string
	Issue                     string
}

func CmdBlockRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := BlockOptions{
		GlobalOptions: GlobalOptions{
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

	return &CommandRegistryEntry{
		"Mark issues as blocker",
		func() error {
			return CmdBlock(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdBlockUsage(cmd, &opts)
		},
	}
}

func CmdBlockUsage(cmd *kingpin.CmdClause, opts *BlockOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	EditorUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("comment", "Comment message when marking issue as blocker").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Comment = &jiradata.Comment{
			Body: flagValue(ctx, "comment"),
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

	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, opts.Endpoint.Value, opts.Issue)
	fmt.Printf("OK %s %s/browse/%s\n", opts.Blocker, opts.Endpoint.Value, opts.Blocker)

	if opts.Browse.Value {
		if err := CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue}); err != nil {
			return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Blocker})
		}
	}

	return nil
}
