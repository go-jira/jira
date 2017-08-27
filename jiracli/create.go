package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CreateOptions struct {
	GlobalOptions        `yaml:",inline" figtree:",inline"`
	jiradata.IssueUpdate `yaml:",inline" figtree:",inline"`
	Project              string
	IssueType            string
	Overrides            map[string]string
}

func CmdCreateRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := CreateOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("create"),
		},
		Overrides: map[string]string{},
	}

	return &CommandRegistryEntry{
		"Create issue",
		func() error {
			return CmdCreate(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdCreateUsage(cmd, &opts)
		},
	}
}

func CmdCreateUsage(cmd *kingpin.CmdClause, opts *CreateOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	EditorUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
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
func CmdCreate(o *oreo.Client, opts *CreateOptions) error {
	type templateInput struct {
		Meta      *jiradata.CreateMetaIssueType `yaml:"meta" json:"meta"`
		Overrides map[string]string             `yaml:"overrides" json:"overrides"`
	}

	if err := defaultIssueType(o, opts.Endpoint.Value, &opts.Project, &opts.IssueType); err != nil {
		return err
	}
	createMeta, err := jira.GetIssueCreateMetaIssueType(o, opts.Endpoint.Value, opts.Project, opts.IssueType)
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
	input.Overrides["user"] = opts.User.Value

	var issueResp *jiradata.IssueCreateResponse
	err = editLoop(&opts.GlobalOptions, &input, &issueUpdate, func() error {
		issueResp, err = jira.CreateIssue(o, opts.Endpoint.Value, &issueUpdate)
		return err
	})
	if err != nil {
		return err
	}

	fmt.Printf("OK %s %s/browse/%s\n", issueResp.Key, opts.Endpoint.Value, issueResp.Key)

	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, issueResp.Key})
	}
	return nil
}

func defaultIssueType(o *oreo.Client, endpoint string, project, issuetype *string) error {
	if project == nil || *project == "" {
		return fmt.Errorf("Project undefined, please use --project argument or set the `project` config property")
	}
	if issuetype != nil && *issuetype != "" {
		return nil
	}
	projectMeta, err := jira.GetIssueCreateMetaProject(o, endpoint, *project)
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
