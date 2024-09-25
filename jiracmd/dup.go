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

type DupOptions struct {
	jiracli.CommonOptions     `yaml:",inline" json:",inline" figtree:",inline"`
	jiradata.LinkIssueRequest `yaml:",inline" json:",inline" figtree:",inline"`
	Project                   string `yaml:"project,omitempty" json:"project,omitempty"`
}

func CmdDupRegistry() *jiracli.CommandRegistryEntry {
	opts := DupOptions{
		CommonOptions: jiracli.CommonOptions{
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

	return &jiracli.CommandRegistryEntry{
		"Mark issues as duplicate",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdDupUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.OutwardIssue.Key = jiracli.FormatIssue(opts.OutwardIssue.Key, opts.Project)
			opts.InwardIssue.Key = jiracli.FormatIssue(opts.InwardIssue.Key, opts.Project)
			return CmdDup(o, globals, &opts)
		},
	}
}

func CmdDupUsage(cmd *kingpin.CmdClause, opts *DupOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("comment", "Comment message when marking issue as duplicate").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Comment = &jiradata.Comment{
			Body: jiracli.FlagValue(ctx, "comment"),
		}
		return nil
	}).String()
	cmd.Arg("DUPLICATE", "duplicate issue to mark closed").Required().StringVar(&opts.InwardIssue.Key)
	cmd.Arg("ISSUE", "duplicate issue to leave open").Required().StringVar(&opts.OutwardIssue.Key)
	return nil
}

// CmdDup will update the given issue as being a duplicate by the given dup issue
// and will attempt to resolve the dup issue
func CmdDup(o *oreo.Client, globals *jiracli.GlobalOptions, opts *DupOptions) error {
	if err := jira.LinkIssues(o, globals.Endpoint.Value, &opts.LinkIssueRequest); err != nil {
		return err
	}
	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", opts.OutwardIssue.Key, jira.URLJoin(globals.Endpoint.Value, "browse", opts.OutwardIssue.Key))
	}

	meta, err := jira.GetIssueTransitions(o, globals.Endpoint.Value, opts.InwardIssue.Key)
	if err != nil {
		return err
	}
	for _, trans := range []string{"close", "done", "cancel", "start", "stop"} {
		transMeta := meta.Transitions.Find(trans)
		if transMeta == nil {
			continue
		}
		issueUpdate := jiradata.IssueUpdate{
			Transition: transMeta,
		}
		resolution := defaultResolution(transMeta)
		if resolution != "" {
			issueUpdate.Fields = map[string]interface{}{
				"resolution": map[string]interface{}{
					"name": resolution,
				},
			}
		}
		if err = jira.TransitionIssue(o, globals.Endpoint.Value, opts.InwardIssue.Key, &issueUpdate); err != nil {
			return err
		}
		if trans != "start" {
			break
		}
		// if we are here then we must be stopping, so need to reset the meta
		meta, err = jira.GetIssueTransitions(o, globals.Endpoint.Value, opts.InwardIssue.Key)
		if err != nil {
			return err
		}
	}

	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", opts.InwardIssue.Key, jira.URLJoin(globals.Endpoint.Value, "browse", opts.InwardIssue.Key))
	}

	if opts.Browse.Value {
		if err := CmdBrowse(globals, opts.OutwardIssue.Key); err != nil {
			return CmdBrowse(globals, opts.InwardIssue.Key)
		}
	}

	return nil
}
