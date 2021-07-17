package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type BlockOptions struct {
	jiracli.CommonOptions     `yaml:",inline" json:",inline" figtree:",inline"`
	jiradata.LinkIssueRequest `yaml:",inline" json:",inline" figtree:",inline"`
	Project                   string `yaml:"project,omitempty" json:"project,omitempty"`
}

func CmdBlockRegistry() *jiracli.CommandRegistryEntry {
	opts := BlockOptions{
		CommonOptions: jiracli.CommonOptions{
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
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdBlockUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.OutwardIssue.Key = jiracli.FormatIssue(opts.OutwardIssue.Key, opts.Project)
			opts.InwardIssue.Key = jiracli.FormatIssue(opts.InwardIssue.Key, opts.Project)
			return CmdBlock(o, globals, &opts)
		},
	}
}

func CmdBlockUsage(cmd *kingpin.CmdClause, opts *BlockOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("comment", "Comment message when marking issue as blocker").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Comment = &jiradata.Comment{
			Body: jiracli.FlagValue(ctx, "comment"),
		}
		return nil
	}).String()
	cmd.Arg("ISSUE", "issue that is blocked").Required().StringVar(&opts.OutwardIssue.Key)
	cmd.Arg("BLOCKER", "blocker issue").Required().StringVar(&opts.InwardIssue.Key)
	return nil
}

// CmdBlock will update the given issue as being a duplicate by the given dup issue
// and will attempt to resolve the dup issue
func CmdBlock(o *oreo.Client, globals *jiracli.GlobalOptions, opts *BlockOptions) error {
	if err := jira.LinkIssues(o, globals.Endpoint.Value, &opts.LinkIssueRequest); err != nil {
		return err
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", opts.InwardIssue.Key, jira.URLJoin(globals.Endpoint.Value, "browse", opts.InwardIssue.Key))
		fmt.Printf("OK %s %s\n", opts.OutwardIssue.Key, jira.URLJoin(globals.Endpoint.Value, "browse", opts.OutwardIssue.Key))
	}

	if opts.Browse.Value {
		if err := CmdBrowse(globals, opts.InwardIssue.Key); err != nil {
			return CmdBrowse(globals, opts.OutwardIssue.Key)
		}
	}

	return nil
}
