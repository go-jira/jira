package jiracmd

import (
	"fmt"
	"strings"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"

	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type TransitionOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string            `yaml:"project,omitempty" json:"project,omitempty"`
	Overrides             map[string]string `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	Transition            string            `yaml:"transition,omitempty" json:"transition,omitempty"`
	Issue                 string            `yaml:"issue,omitempty" json:"issue,omitempty"`
	Resolution            string            `yaml:"resolution,omitempty" json:"resolution,omitempty"`
}

func CmdTransitionRegistry(transition string) *jiracli.CommandRegistryEntry {
	opts := TransitionOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("transition"),
		},
		Overrides: map[string]string{},
	}

	help := "Transition issue to given state"
	if transition != "" {
		help = fmt.Sprintf("Transition issue to %s state", transition)
		opts.SkipEditing = figtree.NewBoolOption(true)
	}

	return &jiracli.CommandRegistryEntry{
		help,
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			if opts.Transition == "" {
				opts.Transition = transition
			}
			return CmdTransitionUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdTransition(o, globals, &opts)
		},
	}
}

func CmdTransitionUsage(cmd *kingpin.CmdClause, opts *TransitionOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Flag("noedit", "Disable opening the editor").SetValue(&opts.SkipEditing)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = jiracli.FlagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set issue property").Short('o').StringMapVar(&opts.Overrides)
	if opts.Transition == "" {
		cmd.Arg("TRANSITION", "State to transition issue to").Required().StringVar(&opts.Transition)
	}
	cmd.Arg("ISSUE", "issue to transition").Required().StringVar(&opts.Issue)
	cmd.Flag("resolution", "Set resolution on transition").StringVar(&opts.Resolution)
	return nil
}

func defaultResolution(transMeta *jiradata.Transition) string {
	if resField, ok := transMeta.Fields["resolution"]; ok {
		for _, allowedValueRaw := range resField.AllowedValues {
			if allowedValue, ok := allowedValueRaw.(map[string]interface{}); ok {
				if allowedValue["name"] == "Fixed" {
					return "Fixed"
				} else if allowedValue["name"] == "Done" {
					return "Done"
				}
			}
		}
	}
	return ""
}

// CmdTransition will move state of the given issue to the given transtion
func CmdTransition(o *oreo.Client, globals *jiracli.GlobalOptions, opts *TransitionOptions) error {
	if globals.JiraDeploymentType.Value == "" {
		serverInfo, err := jira.ServerInfo(o, globals.Endpoint.Value)
		if err != nil {
			return err
		}
		globals.JiraDeploymentType.Value = strings.ToLower(serverInfo.DeploymentType)
	}

	issueData, err := jira.GetIssue(o, globals.Endpoint.Value, opts.Issue, nil)
	if err != nil {
		return jiracli.CliError(err)
	}

	meta, err := jira.GetIssueTransitions(o, globals.Endpoint.Value, opts.Issue)
	if err != nil {
		return jiracli.CliError(err)
	}
	transMeta := meta.Transitions.Find(opts.Transition)

	if transMeta == nil {
		possible := []string{}
		for _, trans := range meta.Transitions {
			possible = append(possible, trans.Name)
		}

		if status, ok := issueData.Fields["status"].(map[string]interface{}); ok {
			if name, ok := status["name"].(string); ok {
				return jiracli.CliError(fmt.Errorf("Invalid Transition %q from %q, Available: %s", opts.Transition, name, strings.Join(possible, ", ")))
			}
		}
		return jiracli.CliError(fmt.Errorf("No valid transition found matching %s", opts.Transition))
	}

	// need to default the Resolution, usually Fixed works but sometime need Done
	if resField, ok := transMeta.Fields["resolution"]; ok && opts.Resolution == "" {
		for _, allowedValueRaw := range resField.AllowedValues {
			if allowedValue, ok := allowedValueRaw.(map[string]interface{}); ok {
				if allowedValue["name"] == "Fixed" {
					opts.Resolution = "Fixed"
				} else if allowedValue["name"] == "Done" {
					opts.Resolution = "Done"
				}
			}
		}
	}
	opts.Overrides["resolution"] = opts.Resolution

	type templateInput struct {
		*jiradata.Issue `yaml:",inline"`
		// Yes, Meta and Transition are redundant, but this is for backwards compatibility
		// with old templates
		Meta       *jiradata.Transition `yaml:"meta,omitempty" json:"meta,omitempty"`
		Transition *jiradata.Transition `yaml:"transition,omitempty" json:"transition,omitempty"`
		Overrides  map[string]string    `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	}

	if _, ok := transMeta.Fields["comment"]; !ok && opts.Overrides["comment"] != "" {
		comment := jiradata.Comment{
			Body: opts.Overrides["comment"],
		}
		if _, err := jira.IssueAddComment(o, globals.Endpoint.Value, opts.Issue, &comment); err != nil {
			return err
		}
	}

	issueUpdate := jiradata.IssueUpdate{}
	input := templateInput{
		Issue:      issueData,
		Meta:       transMeta,
		Transition: transMeta,
		Overrides:  opts.Overrides,
	}
	err = jiracli.EditLoop(&opts.CommonOptions, &input, &issueUpdate, func() error {
		if globals.JiraDeploymentType.Value == jiracli.CloudDeploymentType {
			err := fixGDPRUserFields(o, globals.Endpoint.Value, transMeta.Fields, issueUpdate.Fields)
			if err != nil {
				return err
			}
		}

		// if issueUpdate contains fields lets see if we can map them
		// to their ids
		if len(issueUpdate.Fields) > 0 {
			fields, err := jira.GetFields(o, globals.Endpoint.Value)
			if err != nil {
				return err
			}
			for k, v := range issueUpdate.Fields {
				for _, f := range fields {
					if f.Name == k {
						// re-map to field.id
						issueUpdate.Fields[f.ID] = v
						delete(issueUpdate.Fields, k)
						break
					}
				}
			}
		}

		return jira.TransitionIssue(o, globals.Endpoint.Value, opts.Issue, &issueUpdate)
	})
	if err != nil {
		return jiracli.CliError(err)
	}
	if !globals.Quiet.Value {
		fmt.Printf("OK %s %s\n", issueData.Key, jira.URLJoin(globals.Endpoint.Value, "browse", issueData.Key))
	}

	if opts.Browse.Value {
		return CmdBrowse(globals, opts.Issue)
	}
	return nil
}
