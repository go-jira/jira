package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type LabelsRemoveOptions struct {
	GlobalOptions
	Issue  string
	Labels []string
}

func (jc *JiraCli) CmdLabelsRemoveRegistry() *CommandRegistryEntry {
	opts := LabelsRemoveOptions{}
	return &CommandRegistryEntry{
		"Remove labels from an issue",
		func() error {
			return jc.CmdLabelsRemove(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdLabelsRemoveUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdLabelsRemoveUsage(cmd *kingpin.CmdClause, opts *LabelsRemoveOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	cmd.Arg("ISSUE", "issue id to modify labels").Required().StringVar(&opts.Issue)
	cmd.Arg("LABEL", "label to remove from issue").Required().StringsVar(&opts.Labels)
	return nil
}

// CmdLabels will remove labels on a given issue
func (jc *JiraCli) CmdLabelsRemove(opts *LabelsRemoveOptions) error {
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

	err := jc.EditIssue(opts.Issue, &issueUpdate)
	if err != nil {
		return err
	}
	fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)
	return nil
}
