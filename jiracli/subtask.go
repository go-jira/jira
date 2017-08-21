package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type SubtaskOptions struct {
	GlobalOptions
	jiradata.IssueUpdate
	Project   string
	IssueType string
	Overrides map[string]string
	Issue     string
}

func (jc *JiraCli) CmdSubtaskRegistry() *CommandRegistryEntry {
	opts := SubtaskOptions{
		GlobalOptions: GlobalOptions{
			Template: "subtask",
		},
		IssueType: "Sub-task",
		Overrides: map[string]string{},
	}

	return &CommandRegistryEntry{
		"Subtask issue",
		func() error {
			return jc.CmdSubtask(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdSubtaskUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdSubtaskUsage(cmd *kingpin.CmdClause, opts *SubtaskOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").BoolVar(&opts.SkipEditing)
	cmd.Flag("project", "project to subtask issue in").Short('p').StringVar(&opts.Project)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = flagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set issue property").Short('o').StringMapVar(&opts.Overrides)
	cmd.Arg("ISSUE", "Parent issue for subtask").StringVar(&opts.Issue)
	return nil
}

// CmdSubtask sends the subtask-metadata to the "subtask" template for editing, then
// will parse the edited document as YAML and submit the document to jira.
func (jc *JiraCli) CmdSubtask(opts *SubtaskOptions) error {
	type templateInput struct {
		Meta      *jiradata.CreateMetaIssueType `yaml:"meta" json:"meta"`
		Overrides map[string]string             `yaml:"overrides" json:"overrides"`
		Parent    *jiradata.Issue               `yaml:"parent" json:"parent"`
	}

	parent, err := jc.GetIssue(opts.Issue, nil)
	if err != nil {
		return err
	}

	if project, ok := parent.Fields["project"].(map[string]interface{}); ok {
		if key, ok := project["key"].(string); ok {
			opts.Project = key
		} else {
			return fmt.Errorf("Failed to find Project Key in parent issue")
		}
	} else {
		return fmt.Errorf("Failed to find Project field in parent issue")
	}

	createMeta, err := jc.GetIssueCreateMetaIssueType(opts.Project, opts.IssueType)
	if err != nil {
		return err
	}

	issueUpdate := jiradata.IssueUpdate{}
	input := templateInput{
		Meta:      createMeta,
		Overrides: opts.Overrides,
		Parent:    parent,
	}
	input.Overrides["project"] = opts.Project
	input.Overrides["issuetype"] = opts.IssueType
	input.Overrides["user"] = opts.User

	var issueResp *jiradata.IssueCreateResponse
	err = jc.editLoop(&opts.GlobalOptions, &input, &issueUpdate, func() error {
		issueResp, err = jc.CreateIssue(&issueUpdate)
		return err
	})
	if err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", issueResp.Key, jc.Endpoint, issueResp.Key)

	if opts.Browse {
		return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
