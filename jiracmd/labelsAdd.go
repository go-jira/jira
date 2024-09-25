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

type LabelsAddOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string   `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string   `yaml:"issue,omitempty" json:"issue,omitempty"`
	Labels                []string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

func CmdLabelsAddRegistry() *jiracli.CommandRegistryEntry {
	opts := LabelsAddOptions{}
	return &jiracli.CommandRegistryEntry{
		"Add labels to an issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdLabelsAddUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdLabelsAdd(o, globals, &opts)
		},
	}
}

func CmdLabelsAddUsage(cmd *kingpin.CmdClause, opts *LabelsAddOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "issue id to modify labels").Required().StringVar(&opts.Issue)
	cmd.Arg("LABEL", "label to add to issue").Required().StringsVar(&opts.Labels)
	return nil
}

// CmdLabelsAdd will add labels on a given issue
func CmdLabelsAdd(o *oreo.Client, globals *jiracli.GlobalOptions, opts *LabelsAddOptions) error {
	ops := jiradata.FieldOperations{}
	for _, label := range opts.Labels {
		ops = append(ops, jiradata.FieldOperation{
			"add": label,
		})
	}
	issueUpdate := jiradata.IssueUpdate{
		Update: jiradata.FieldOperationsMap{
			"labels": ops,
		},
	}

	if err := jira.EditIssue(o, globals.Endpoint.Value, opts.Issue, &issueUpdate); err != nil {
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
