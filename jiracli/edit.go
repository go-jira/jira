package jiracli

import (
	"fmt"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EditOptions struct {
	GlobalOptions
	jiradata.IssueUpdate
	jira.SearchOptions
	Overrides map[string]string
}

func (jc *JiraCli) CmdEditRegistry() *CommandRegistryEntry {
	issue := ""
	opts := EditOptions{
		GlobalOptions: GlobalOptions{
			Template: "edit",
		},
		Overrides: map[string]string{},
	}

	return &CommandRegistryEntry{
		"Edit issue details",
		func() error {
			return jc.CmdEdit(issue, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdEditUsage(cmd, &issue, &opts)
		},
	}
}

func (jc *JiraCli) CmdEditUsage(cmd *kingpin.CmdClause, issue *string, opts *EditOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").BoolVar(&opts.SkipEditing)
	// cmd.Flag("assignee", "User assigned the issue").Short('a').StringVar(&opts.Assignee)
	// cmd.Flag("component", "Component to search for").Short('c').StringVar(&opts.Component)
	// cmd.Flag("issuetype", "Issue type to search for").Short('i').StringVar(&opts.IssueType)
	// cmd.Flag("limit", "Maximum number of results to return in search").Short('l').Default("500").IntVar(&opts.MaxResults)
	// cmd.Flag("project", "Project to search for").Short('p').StringVar(&opts.Project)
	cmd.Flag("query", "Jira Query Language (JQL) expression for the search to edit multiple issues").Short('q').StringVar(&opts.Query)
	// cmd.Flag("reporter", "Reporter to search for").Short('r').StringVar(&opts.Reporter)
	// cmd.Flag("sort", "Sort order to return").Short('s').Default("priority asc, key").StringVar(&opts.Sort)
	// cmd.Flag("watcher", "Watcher to search for").Short('w').StringVar(&opts.Watcher)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = flagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set issue property").Short('o').StringMapVar(&opts.Overrides)
	cmd.Arg("ISSUE", "issue id to edit").StringVar(issue)
	return nil
}

// Edit will get issue data and send to "edit" template
func (jc *JiraCli) CmdEdit(issue string, opts *EditOptions) error {
	type templateInput struct {
		*jiradata.Issue `yaml:",inline"`
		Meta            *jiradata.EditMeta `yaml:"meta" json:"meta"`
		Overrides       map[string]string  `yaml:"overrides" json:"overrides"`
	}
	if issue != "" {
		issueData, err := jc.GetIssue(issue, nil)
		if err != nil {
			return err
		}
		editMeta, err := jc.GetIssueEditMeta(issue)
		if err != nil {
			return err
		}

		issueUpdate := jiradata.IssueUpdate{}
		input := templateInput{
			Issue:     issueData,
			Meta:      editMeta,
			Overrides: opts.Overrides,
		}
		return jc.editLoop(&opts.GlobalOptions, &input, &issueUpdate, func() error {
			return jc.EditIssue(issue, &issueUpdate)
		})
		// FIXME implement browse
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
		link := fmt.Sprintf("%s/browse/%s", jc.Endpoint, issueData.Key)
		fmt.Printf("OK %s %s\n", issueData.Key, link)
		// FIXME implement browse
	}
	return nil
}
