package jiracli

import kingpin "gopkg.in/alecthomas/kingpin.v2"

type TransitionsOptions struct {
	GlobalOptions
	Issue string
}

func (jc *JiraCli) CmdTransitionsRegistry(defaultTemplate string) *CommandRegistryEntry {
	opts := TransitionsOptions{
		GlobalOptions: GlobalOptions{
			Template: defaultTemplate,
		},
	}

	return &CommandRegistryEntry{
		"List valid issue transitions",
		func() error {
			return jc.CmdTransitions(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdTransitionsUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdTransitionsUsage(cmd *kingpin.CmdClause, opts *TransitionsOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue to list valid transitions").Required().StringVar(&opts.Issue)
	return nil
}

// Transitions will get issue edit metadata and send to "editmeta" template
func (jc *JiraCli) CmdTransitions(opts *TransitionsOptions) error {
	editMeta, err := jc.GetIssueTransitions(opts.Issue)
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template, editMeta, nil)
}
