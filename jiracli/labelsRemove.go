package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type LabelsRemoveOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Issue         string
	Labels        []string
}

func CmdLabelsRemoveRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := LabelsRemoveOptions{}
	return &CommandRegistryEntry{
		"Remove labels from an issue",
		func() error {
			return CmdLabelsRemove(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdLabelsRemoveUsage(cmd, &opts)
		},
	}
}

func CmdLabelsRemoveUsage(cmd *kingpin.CmdClause, opts *LabelsRemoveOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue id to modify labels").Required().StringVar(&opts.Issue)
	cmd.Arg("LABEL", "label to remove from issue").Required().StringsVar(&opts.Labels)
	return nil
}

// CmdLabels will remove labels on a given issue
func CmdLabelsRemove(o *oreo.Client, opts *LabelsRemoveOptions) error {
	ops := jiradata.FieldOperations{}
	for _, label := range opts.Labels {
		ops = append(ops, jiradata.FieldOperation{
			"remove": label,
		})
	}
	issueUpdate := jiradata.IssueUpdate{
		Update: jiradata.FieldOperationsMap{
			"labels": ops,
		},
	}

	err := jira.EditIssue(o, opts.Endpoint.Value, opts.Issue, &issueUpdate)
	if err != nil {
		return err
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, opts.Endpoint.Value, opts.Issue)
	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
