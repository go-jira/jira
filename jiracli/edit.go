package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EditOptions struct {
	GlobalOptions        `yaml:",inline" figtree:",inline"`
	jiradata.IssueUpdate `yaml:",inline" figtree:",inline"`
	jira.SearchOptions   `yaml:",inline" figtree:",inline"`
	Overrides            map[string]string
	Issue                string
}

func CmdEditRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := EditOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("edit"),
		},
		Overrides: map[string]string{},
	}

	return &CommandRegistryEntry{
		"Edit issue details",
		func() error {
			return CmdEdit(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdEditUsage(cmd, &opts)
		},
	}
}

func CmdEditUsage(cmd *kingpin.CmdClause, opts *EditOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	EditorUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
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
func CmdEdit(o *oreo.Client, opts *EditOptions) error {
	type templateInput struct {
		*jiradata.Issue `yaml:",inline"`
		Meta            *jiradata.EditMeta `yaml:"meta" json:"meta"`
		Overrides       map[string]string  `yaml:"overrides" json:"overrides"`
	}
	if opts.Issue != "" {
		issueData, err := jira.GetIssue(o, opts.Endpoint.Value, opts.Issue, nil)
		if err != nil {
			return err
		}
		editMeta, err := jira.GetIssueEditMeta(o, opts.Endpoint.Value, opts.Issue)
		if err != nil {
			return err
		}

		issueUpdate := jiradata.IssueUpdate{}
		input := templateInput{
			Issue:     issueData,
			Meta:      editMeta,
			Overrides: opts.Overrides,
		}
		err = editLoop(&opts.GlobalOptions, &input, &issueUpdate, func() error {
			return jira.EditIssue(o, opts.Endpoint.Value, opts.Issue, &issueUpdate)
		})
		if err != nil {
			return err
		}
		fmt.Printf("OK %s %s/browse/%s\n", opts.Issue, opts.Endpoint.Value, opts.Issue)

		if opts.Browse.Value {
			return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
		}
	}
	results, err := jira.Search(o, opts.Endpoint.Value, opts)
	if err != nil {
		return err
	}
	for _, issueData := range results.Issues {
		editMeta, err := jira.GetIssueEditMeta(o, opts.Endpoint.Value, issueData.Key)
		if err != nil {
			return err
		}

		issueUpdate := jiradata.IssueUpdate{}
		input := templateInput{
			Issue: issueData,
			Meta:  editMeta,
		}
		err = editLoop(&opts.GlobalOptions, &input, &issueUpdate, func() error {
			return jira.EditIssue(o, opts.Endpoint.Value, issueData.Key, &issueUpdate)
		})
		if err != nil {
			return err
		}
		fmt.Printf("OK %s %s/browse/%s\n", issueData.Key, opts.Endpoint.Value, issueData.Key)

		if opts.Browse.Value {
			return CmdBrowse(&BrowseOptions{opts.GlobalOptions, issueData.Key})
		}
	}
	return nil
}
