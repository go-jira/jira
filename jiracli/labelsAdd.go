package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type LabelsAddOptions struct {
	GlobalOptions
	Issue  string
	Labels []string
}

func (jc *JiraCli) CmdLabelsAddRegistry() *CommandRegistryEntry {
	opts := LabelsAddOptions{}
	return &CommandRegistryEntry{
		"Add labels to an issue",
		func() error {
			return jc.CmdLabelsAdd(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdLabelsAddUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdLabelsAddUsage(cmd *kingpin.CmdClause, opts *LabelsAddOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	cmd.Arg("ISSUE", "issue id to modify labels").Required().StringVar(&opts.Issue)
	cmd.Arg("LABEL", "label to add to issue").Required().StringsVar(&opts.Labels)
	return nil
}

// CmdLabels will add labels on a given issue
func (jc *JiraCli) CmdLabelsAdd(opts *LabelsAddOptions) error {
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

	err := jc.EditIssue(opts.Issue, &issueUpdate)
	if err != nil {
		return err
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)
	return nil
}
