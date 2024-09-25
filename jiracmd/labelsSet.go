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

type LabelsSetOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string   `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string   `yaml:"issue,omitempty" json:"issue,omitempty"`
	Labels                []string `yaml:"labels,omitempty" json:"labels,omitempty"`
}

func CmdLabelsSetRegistry() *jiracli.CommandRegistryEntry {
	opts := LabelsSetOptions{}
	return &jiracli.CommandRegistryEntry{
		"Set labels on an issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdLabelsSetUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdLabelsSet(o, globals, &opts)
		},
	}
}

func CmdLabelsSetUsage(cmd *kingpin.CmdClause, opts *LabelsSetOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "issue id to modify labels").Required().StringVar(&opts.Issue)
	cmd.Arg("LABEL", "label to set on issue").Required().StringsVar(&opts.Labels)
	return nil
}

// CmdLabelsSet will set labels on a given issue
func CmdLabelsSet(o *oreo.Client, globals *jiracli.GlobalOptions, opts *LabelsSetOptions) error {
	issueUpdate := jiradata.IssueUpdate{
		Update: jiradata.FieldOperationsMap{
			"labels": jiradata.FieldOperations{
				jiradata.FieldOperation{
					"set": opts.Labels,
				},
			},
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
