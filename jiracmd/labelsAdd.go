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

type LabelsAddOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue         string   `yaml:"issue,omitempty" json:"issue,omitempty"`
	Labels        []string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

func CmdLabelsAddRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := LabelsAddOptions{}
	return &jiracli.CommandRegistryEntry{
		"Add labels to an issue",
		func() error {
			return CmdLabelsAdd(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdLabelsAddUsage(cmd, &opts)
		},
	}
}

func CmdLabelsAddUsage(cmd *kingpin.CmdClause, opts *LabelsAddOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jiracli.BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue id to modify labels").Required().StringVar(&opts.Issue)
	cmd.Arg("LABEL", "label to add to issue").Required().StringsVar(&opts.Labels)
	return nil
}

// CmdLabels will add labels on a given issue
func CmdLabelsAdd(o *oreo.Client, opts *LabelsAddOptions) error {
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

	if err := jira.EditIssue(o, opts.Endpoint.Value, opts.Issue, &issueUpdate); err != nil {
		return err
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, opts.Endpoint.Value, opts.Issue)
	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
