package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type DupOptions struct {
	GlobalOptions             `yaml:",inline" figtree:",inline"`
	jiradata.LinkIssueRequest `yaml:",inline" figtree:",inline"`
	Duplicate                 string
	Issue                     string
}

func CmdDupRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := DupOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("edit"),
		},
		LinkIssueRequest: jiradata.LinkIssueRequest{
			Type: &jiradata.IssueLinkType{
				Name: "Duplicate",
			},
			InwardIssue:  &jiradata.IssueRef{},
			OutwardIssue: &jiradata.IssueRef{},
		},
	}

	return &CommandRegistryEntry{
		"Mark issues as duplicate",
		func() error {
			return CmdDup(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdDupUsage(cmd, &opts)
		},
	}
}

func CmdDupUsage(cmd *kingpin.CmdClause, opts *DupOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	EditorUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("comment", "Comment message when marking issue as duplicate").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Comment = &jiradata.Comment{
			Body: flagValue(ctx, "comment"),
		}
		return nil
	}).String()
	cmd.Arg("DUPLICATE", "duplicate issue to mark closed").Required().StringVar(&opts.InwardIssue.Key)
	cmd.Arg("ISSUE", "duplicate issue to leave open").Required().StringVar(&opts.OutwardIssue.Key)
	return nil
}

// CmdDups will update the given issue as being a duplicate by the given dup issue
// and will attempt to resolve the dup issue
func CmdDup(o *oreo.Client, opts *DupOptions) error {
	if err := jira.LinkIssues(o, opts.Endpoint.Value, &opts.LinkIssueRequest); err != nil {
		return err
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.OutwardIssue.Key, opts.Endpoint.Value, opts.OutwardIssue.Key)

	meta, err := jira.GetIssueTransitions(o, opts.Endpoint.Value, opts.InwardIssue.Key)
	if err != nil {
		return err
	}
	for _, trans := range []string{"close", "done", "start", "stop"} {
		transMeta := meta.Transitions.Find(trans)
		if transMeta != nil {
			issueUpdate := jiradata.IssueUpdate{
				Transition: transMeta,
			}
			if err = jira.TransitionIssue(o, opts.Endpoint.Value, opts.InwardIssue.Key, &issueUpdate); err != nil {
				return err
			}
			// if we just started the issue now we need to stop it
			if trans != "start" {
				break
			}
		}
	}

	fmt.Printf("OK %s %s/browse/%s\n", opts.OutwardIssue.Key, opts.Endpoint.Value, opts.OutwardIssue.Key)
	fmt.Printf("OK %s %s/browse/%s\n", opts.InwardIssue.Key, opts.Endpoint.Value, opts.InwardIssue.Key)

	if opts.Browse.Value {
		if err := CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.OutwardIssue.Key}); err != nil {
			return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.InwardIssue.Key})
		}
	}

	return nil
}
