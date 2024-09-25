package jiracmd

import (
	"fmt"
	"strings"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	"gopkg.in/AlecAivazis/survey.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type EditOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	jiradata.IssueUpdate  `yaml:",inline" json:",inline" figtree:",inline"`
	jira.SearchOptions    `yaml:",inline" json:",inline" figtree:",inline"`
	Overrides             map[string]string `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	Issue                 string            `yaml:"issue,omitempty" json:"issue,omitempty"`
	Queries               map[string]string `yaml:"queries,omitempty" json:"queries,omitempty"`
}

func CmdEditRegistry() *jiracli.CommandRegistryEntry {
	opts := EditOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("edit"),
		},
		Overrides: map[string]string{},
	}

	return &jiracli.CommandRegistryEntry{
		"Edit issue details",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdEditUsage(cmd, &opts, fig)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			if opts.QueryFields == "" {
				opts.QueryFields = "assignee,created,priority,reporter,status,summary,updated,issuetype,comment,description,votes,created,customfield_10110,components"
			}
			return CmdEdit(o, globals, &opts)
		},
	}
}

func CmdEditUsage(cmd *kingpin.CmdClause, opts *EditOptions, fig *figtree.FigTree) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.EditorUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("named-query", "The name of a query in the `queries` configuration").Short('n').PreAction(func(ctx *kingpin.ParseContext) error {
		name := jiracli.FlagValue(ctx, "named-query")
		if query, ok := opts.Queries[name]; ok && query != "" {
			var err error
			opts.Query, err = jiracli.ConfigTemplate(fig, query, cmd.FullCommand(), opts)
			return err
		}
		return fmt.Errorf("A valid named-query %q not found in `queries` configuration", name)
	}).String()
	cmd.Flag("query", "Jira Query Language (JQL) expression for the search to edit multiple issues").Short('q').StringVar(&opts.Query)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = jiracli.FlagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set issue property").Short('o').StringMapVar(&opts.Overrides)
	cmd.Arg("ISSUE", "issue id to edit").StringVar(&opts.Issue)
	return nil
}

// Edit will get issue data and send to "edit" template
func CmdEdit(o *oreo.Client, globals *jiracli.GlobalOptions, opts *EditOptions) error {
	if globals.JiraDeploymentType.Value == "" {
		serverInfo, err := jira.ServerInfo(o, globals.Endpoint.Value)
		if err != nil {
			return err
		}
		globals.JiraDeploymentType.Value = strings.ToLower(serverInfo.DeploymentType)
	}

	type templateInput struct {
		*jiradata.Issue `yaml:",inline"`
		Meta            *jiradata.EditMeta `yaml:"meta" json:"meta"`
		Overrides       map[string]string  `yaml:"overrides" json:"overrides"`
	}
	if opts.Issue != "" {
		issueData, err := jira.GetIssue(o, globals.Endpoint.Value, opts.Issue, nil)
		if err != nil {
			return err
		}
		editMeta, err := jira.GetIssueEditMeta(o, globals.Endpoint.Value, opts.Issue)
		if err != nil {
			return err
		}

		issueUpdate := jiradata.IssueUpdate{}
		input := templateInput{
			Issue:     issueData,
			Meta:      editMeta,
			Overrides: opts.Overrides,
		}
		err = jiracli.EditLoop(&opts.CommonOptions, &input, &issueUpdate, func() error {
			if globals.JiraDeploymentType.Value == jiracli.CloudDeploymentType {
				err := fixGDPRUserFields(o, globals.Endpoint.Value, editMeta.Fields, issueUpdate.Fields)
				if err != nil {
					return err
				}
			}
			return jira.EditIssue(o, globals.Endpoint.Value, opts.Issue, &issueUpdate)
		})
		if err != nil {
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
	results, err := jira.Search(o, globals.Endpoint.Value, opts)
	if err != nil {
		return err
	}
	for i, issueData := range results.Issues {
		editMeta, err := jira.GetIssueEditMeta(o, globals.Endpoint.Value, issueData.Key)
		if err != nil {
			return err
		}

		issueUpdate := jiradata.IssueUpdate{}
		input := templateInput{
			Issue:     issueData,
			Meta:      editMeta,
			Overrides: opts.Overrides,
		}
		err = jiracli.EditLoop(&opts.CommonOptions, &input, &issueUpdate, func() error {
			if globals.JiraDeploymentType.Value == jiracli.CloudDeploymentType {
				err := fixGDPRUserFields(o, globals.Endpoint.Value, editMeta.Fields, issueUpdate.Fields)
				if err != nil {
					return err
				}
			}
			return jira.EditIssue(o, globals.Endpoint.Value, issueData.Key, &issueUpdate)
		})
		if err == jiracli.EditLoopAbort && len(results.Issues) > i+1 {
			var answer bool
			survey.AskOne(
				&survey.Confirm{
					Message: fmt.Sprintf("Continue to edit next issue %s?", results.Issues[i+1].Key),
					Default: true,
				},
				&answer,
				nil,
			)
			if answer {
				continue
			}
			panic(jiracli.Exit{1})
		}
		if err != nil {
			return err
		}
		if !globals.Quiet.Value {
			fmt.Printf("OK %s %s\n", issueData.Key, jira.URLJoin(globals.Endpoint.Value, "browse", issueData.Key))
		}
		if opts.Browse.Value {
			return CmdBrowse(globals, issueData.Key)
		}
	}
	return nil
}

func fixUserField(ua jira.HttpClient, endpoint string, userField map[string]interface{}) error {
	if _, ok := userField["accountId"].(string); ok {
		// this field is already GDPR ready
		return nil
	}

	queryName, ok := userField["displayName"].(string)
	if !ok {
		queryName, ok = userField["emailAddress"].(string)
		if !ok {
			// no fields to search on, skip user lookup
			return nil
		}
	}
	users, err := jira.UserSearch(ua, endpoint, &jira.UserSearchOptions{
		// Query field will search users displayName and emailAddress
		Query: queryName,
	})
	if err != nil {
		return err
	}
	if len(users) != 1 {
		return fmt.Errorf("Found %d accounts for users with query %q", len(users), queryName)
	}
	userField["accountId"] = users[0].AccountID
	return nil
}

func fixGDPRUserFields(ua jira.HttpClient, endpoint string, meta jiradata.FieldMetaMap, fields map[string]interface{}) error {
	for fieldName, fieldMeta := range meta {
		// check to see if meta-field is in fields data, otherwise skip
		if _, ok := fields[fieldName]; !ok {
			continue
		}
		if fieldMeta.Schema.Type == "user" {
			userField, ok := fields[fieldName].(map[string]interface{})
			if !ok {
				// for some reason the field seems to be the wrong type in the data
				// even though the schema is a "user"
				continue
			}
			err := fixUserField(ua, endpoint, userField)
			if err != nil {
				return err
			}
			fields[fieldName] = userField
		}
		if fieldMeta.Schema.Type == "array" && fieldMeta.Schema.Items == "user" {
			listUserField, ok := fields[fieldName].([]interface{})
			if !ok {
				// for some reason the field seems to be the wrong type in the data
				// even though the schema is a list of "user"
				continue
			}
			for i, userFieldItem := range listUserField {
				userField, ok := userFieldItem.(map[string]interface{})
				if !ok {
					// for some reason the field seems to be the wrong type in the data
					// even though the schema is a "user"
					continue
				}
				err := fixUserField(ua, endpoint, userField)
				if err != nil {
					return err
				}
				listUserField[i] = userField
			}
			fields[fieldName] = listUserField
		}
	}
	return nil
}
