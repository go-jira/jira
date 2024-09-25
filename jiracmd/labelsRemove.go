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

type LabelsRemoveOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string   `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string   `yaml:"issue,omitempty" json:"issue,omitempty"`
	Labels                []string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

func CmdLabelsRemoveRegistry() *jiracli.CommandRegistryEntry {
	opts := LabelsRemoveOptions{}
	return &jiracli.CommandRegistryEntry{
		"Remove labels from an issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdLabelsRemoveUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdLabelsRemove(o, globals, &opts)
		},
	}
}

func CmdLabelsRemoveUsage(cmd *kingpin.CmdClause, opts *LabelsRemoveOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "issue id to modify labels").Required().StringVar(&opts.Issue)
	cmd.Arg("LABEL", "label to remove from issue").Required().StringsVar(&opts.Labels)
	return nil
}

// CmdLabelsRemove will remove labels on a given issue
func CmdLabelsRemove(o *oreo.Client, globals *jiracli.GlobalOptions, opts *LabelsRemoveOptions) error {
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

	err := jira.EditIssue(o, globals.Endpoint.Value, opts.Issue, &issueUpdate)
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
