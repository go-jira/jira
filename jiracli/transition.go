package jiracli

import (
	"fmt"
	"strings"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type TransitionOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Overrides     map[string]string
	Transition    string
	Issue         string
	Resolution    string
}

func CmdTransitionRegistry(fig *figtree.FigTree, o *oreo.Client, transition string) *CommandRegistryEntry {
	opts := TransitionOptions{
		GlobalOptions: GlobalOptions{
			Template: figtree.NewStringOption("transition"),
		},
		Transition: transition,
		Overrides:  map[string]string{},
	}

	help := "Transition issue to given state"
	if transition == "" {
		help = fmt.Sprintf("Transition issue to %s state", transition)
	}

	return &CommandRegistryEntry{
		help,
		func() error {
			return CmdTransition(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdTransitionUsage(cmd, &opts)
		},
	}
}

func CmdTransitionUsage(cmd *kingpin.CmdClause, opts *TransitionOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	BrowseUsage(cmd, &opts.GlobalOptions)
	TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Flag("comment", "Comment message for issue").Short('m').PreAction(func(ctx *kingpin.ParseContext) error {
		opts.Overrides["comment"] = flagValue(ctx, "comment")
		return nil
	}).String()
	cmd.Flag("override", "Set issue property").Short('o').StringMapVar(&opts.Overrides)
	if opts.Transition == "" {
		cmd.Arg("TRANSITION", "State to transition issue to").Required().StringVar(&opts.Transition)
	}
	cmd.Arg("ISSUE", "issue to transition").Required().StringVar(&opts.Issue)
	return nil
}

// CmdTransition will move state of the given issue to the given transtion
func CmdTransition(o *oreo.Client, opts *TransitionOptions) error {
	issueData, err := jira.GetIssue(o, opts.Endpoint.Value, opts.Issue, nil)
	if err != nil {
		return err
	}

	meta, err := jira.GetIssueTransitions(o, opts.Endpoint.Value, opts.Issue)
	if err != nil {
		return err
	}
	transMeta := meta.Transitions.Find(opts.Transition)

	if transMeta == nil {
		possible := []string{}
		for _, trans := range meta.Transitions {
			possible = append(possible, trans.Name)
		}

		if status, ok := issueData.Fields["status"].(map[string]interface{}); ok {
			if name, ok := status["name"].(string); ok {
				return fmt.Errorf("Invalid Transition %q from %q, Available: %s", opts.Transition, name, strings.Join(possible, ", "))
			}
		}
		return fmt.Errorf("No valid transition found matching %s", opts.Transition)
	}

	// need to default the Resolution, usually Fixed works but sometime need Done
	if opts.Resolution == "" {
		if resField, ok := transMeta.Fields["resolution"]; ok {
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
	}
	opts.Overrides["resolution"] = opts.Resolution

	type templateInput struct {
		*jiradata.Issue `yaml:",inline"`
		// Yes, Meta and Transition are redundant, but this is for backwards compatibility
		// with old templates
		Meta       *jiradata.Transition `yaml:"meta,omitempty" json:"meta,omitemtpy"`
		Transition *jiradata.Transition `yaml:"transition,omitempty" json:"transition,omitempty"`
		Overrides  map[string]string    `yaml:"overrides,omitempty" json:"overrides,omitempty"`
	}

	issueUpdate := jiradata.IssueUpdate{}
	input := templateInput{
		Issue:      issueData,
		Meta:       transMeta,
		Transition: transMeta,
		Overrides:  opts.Overrides,
	}
	err = editLoop(&opts.GlobalOptions, &input, &issueUpdate, func() error {
		return jira.TransitionIssue(o, opts.Endpoint.Value, opts.Issue, &issueUpdate)
	})
	if err != nil {
		return err
	}
	fmt.Printf("OK %s %s/browse/%s\n", issueData.Key, opts.Endpoint.Value, issueData.Key)

	if opts.Browse.Value {
		return CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
