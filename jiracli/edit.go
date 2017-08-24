package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EditOptions struct {
	GlobalOptions
	jiradata.IssueUpdate
	jira.SearchOptions
	Overrides map[string]string
	Issue     string
}

func (jc *JiraCli) CmdEditRegistry() *CommandRegistryEntry {
	opts := EditOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("edit"),
		},
		Overrides: map[string]string{},
	}

	return &CommandRegistryEntry{
		"Edit issue details",
		func() error {
			return jc.CmdEdit(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdEditUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdEditUsage(cmd *kingpin.CmdClause, opts *EditOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("query", "Jira Query Language (JQL) expression for the search to edit multiple issues").Short('q').StringVar(&opts.Query)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = flagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set issue property").Short('o').StringMapVar(&opts.Overrides)
	cmd.Arg("ISSUE", "issue id to edit").StringVar(&opts.Issue)
	return nil
}

// Edit will get issue data and send to "edit" template
func (jc *JiraCli) CmdEdit(opts *EditOptions) error {
	type templateInput struct {
		*jiradata.Issue `yaml:",inline"`
		Meta            *jiradata.EditMeta `yaml:"meta" json:"meta"`
		Overrides       map[string]string  `yaml:"overrides" json:"overrides"`
	}
	if opts.Issue != "" {
		issueData, err := jc.GetIssue(opts.Issue, nil)
		if err != nil {
			return err
		}
		editMeta, err := jc.GetIssueEditMeta(opts.Issue)
		if err != nil {
			return err
		}

		issueUpdate := jiradata.IssueUpdate{}
		input := templateInput{
			Issue:     issueData,
			Meta:      editMeta,
			Overrides: opts.Overrides,
		}
		err = jc.editLoop(&opts.GlobalOptions, &input, &issueUpdate, func() error {
			return jc.EditIssue(opts.Issue, &issueUpdate)
		})
		if err != nil {
			return err
		}
		fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, jc.Endpoint, opts.Issue)

		if opts.Browse.Value {
			return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
		}
	}
	results, err := jc.Search(opts)
	if err != nil {
		return err
	}
	for _, issueData := range results.Issues {
		editMeta, err := jc.GetIssueEditMeta(issueData.Key)
		if err != nil {
			return err
		}

		issueUpdate := jiradata.IssueUpdate{}
		input := templateInput{
			Issue: issueData,
			Meta:  editMeta,
		}
		err = jc.editLoop(&opts.GlobalOptions, &input, &issueUpdate, func() error {
			return jc.EditIssue(issueData.Key, &issueUpdate)
		})
		if err != nil {
			return err
		}
		fmt.Printf("OK %s %s/browse/%s\n", issueData.Key, jc.Endpoint, issueData.Key)

		if opts.Browse.Value {
			return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, issueData.Key})
		}
	}
	return nil
}
