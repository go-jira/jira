package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type IssueLinkOptions struct {
	GlobalOptions             `yaml:",inline" figtree:",inline"`
	jiradata.LinkIssueRequest `yaml:",inline" figtree:",inline"`
	LinkType                  string
}

func CmdIssueLinkRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := IssueLinkOptions{
		LinkIssueRequest: jiradata.LinkIssueRequest{
			Type:         &jiradata.IssueLinkType{},
			InwardIssue:  &jiradata.IssueRef{},
			OutwardIssue: &jiradata.IssueRef{},
		},
	}
	return &CommandRegistryEntry{
		"Link two issues",
		func() error {
			return CmdIssueLink(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdIssueLinkUsage(cmd, &opts)
		},
	}
}

func CmdIssueLinkUsage(cmd *kingpin.CmdClause, opts *IssueLinkOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	EditorUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
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
func CmdIssueLink(o *oreo.Client, opts *IssueLinkOptions) error {
	if err := jira.LinkIssues(o, opts.Endpoint.Value, &opts.LinkIssueRequest); err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.InwardIssue.Key, opts.Endpoint.Value, opts.InwardIssue.Key)
	fmt.Printf("OK %s %s/browse/%s\n", opts.OutwardIssue.Key, opts.Endpoint.Value, opts.OutwardIssue.Key)

	if opts.Browse.Value {
		if err := CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.OutwardIssue.Key}); err != nil {
			return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.InwardIssue.Key})
		}
	}

	return nil
}
