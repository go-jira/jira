package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type LabelsSetOptions struct {
	GlobalOptions
	Issue  string
	Labels []string
}

func (jc *JiraCli) CmdLabelsSetRegistry() *CommandRegistryEntry {
	opts := LabelsSetOptions{}
	return &CommandRegistryEntry{
		"Set labels on an issue",
		func() error {
			return jc.CmdLabelsSet(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdLabelsSetUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdLabelsSetUsage(cmd *kingpin.CmdClause, opts *LabelsSetOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue id to modify labels").Required().StringVar(&opts.Issue)
	cmd.Arg("LABEL", "label to set on issue").Required().StringsVar(&opts.Labels)
	return nil
}

// CmdLabels will set labels on a given issue
func (jc *JiraCli) CmdLabelsSet(opts *LabelsSetOptions) error {
	issueUpdate := jiradata.IssueUpdate{
		Update: jiradata.FieldOperationsMap{
			"labels": jiradata.FieldOperations{
				jiradata.FieldOperation{
					"set": opts.Labels,
				},
			},
		},
	}

	if err := jc.EditIssue(opts.Issue, &issueUpdate); err != nil {
		return err
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)
	if opts.Browse.Value {
		return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
