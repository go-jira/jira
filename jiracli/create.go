package jiracli

import (
	"fmt"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CreateOptions struct {
	GlobalOptions
	jiradata.IssueUpdate
	Project   string
	IssueType string
	Overrides map[string]string
}

func (jc *JiraCli) CmdCreateRegistry() *CommandRegistryEntry {
	opts := CreateOptions{
		GlobalOptions: GlobalOptions{
			Template: "create",
		},
		Overrides: map[string]string{},
	}

	return &CommandRegistryEntry{
		"Create issue",
		func() error {
			return jc.CmdCreate(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdCreateUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdCreateUsage(cmd *kingpin.CmdClause, opts *CreateOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.EditorUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").BoolVar(&opts.SkipEditing)
	cmd.Flag("project", "project to create issue in").Short('p').StringVar(&opts.Project)
	cmd.Flag("issuetype", "issuetype in to create").Short('i').StringVar(&opts.IssueType)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = flagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set issue property").Short('o').StringMapVar(&opts.Overrides)
	return nil
}

// CmdCreate sends the create-metadata to the "create" template for editing, then
// will parse the edited document as YAML and submit the document to jira.
func (jc *JiraCli) CmdCreate(opts *CreateOptions) error {
	type templateInput struct {
		Meta      *jiradata.CreateMetaIssueType `yaml:"meta" json:"meta"`
		Overrides map[string]string             `yaml:"overrides" json:"overrides"`
	}

	if err := jc.defaultIssueType(&opts.Project, &opts.IssueType); err != nil {
		return err
	}
	createMeta, err := jc.GetIssueCreateMetaIssueType(opts.Project, opts.IssueType)
	if err != nil {
		return err
	}

	issueUpdate := jiradata.IssueUpdate{}
	input := templateInput{
		Meta:      createMeta,
		Overrides: opts.Overrides,
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

	link := fmt.Sprintf("%s/browse/%s", jc.Endpoint, issueResp.Key)
	fmt.Printf("OK %s %s\n", issueResp.Key, link)

	// FIXME implement browse
	return nil
}

func (jc *JiraCli) defaultIssueType(project, issuetype *string) error {
	if project == nil || *project == "" {
		return fmt.Errorf("Project undefined, please use --project argument or set the `project` config property")
	}
	if issuetype != nil && *issuetype != "" {
		return nil
	}
	projectMeta, err := jc.GetIssueCreateMetaProject(*project)
	if err != nil {
		return err
	}

	issueTypes := map[string]bool{}

	for _, issuetype := range projectMeta.Issuetypes {
		issueTypes[issuetype.Name] = true
	}

	//  prefer "Bug" type
	if _, ok := issueTypes["Bug"]; ok {
		*issuetype = "Bug"
		return nil
	}
	// next best default it "Task"
	if _, ok := issueTypes["Task"]; ok {
		*issuetype = "Task"
		return nil
	}

	return fmt.Errorf("Unable to find default issueType of Bug or Task, please set --issuetype argument or set the `issuetype` config property")
}
