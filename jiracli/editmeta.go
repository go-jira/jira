package jiracli

import kingpin "gopkg.in/alecthomas/kingpin.v2"

func (jc *JiraCli) CmdEditMetaRegistry() *CommandRegistryEntry {
	issue := ""
	opts := GlobalOptions{
		Template: "editmeta",
	}

	return &CommandRegistryEntry{
		"View 'edit' metadata",
		func() error {
			return jc.CmdEditMeta(issue, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdEditMetaUsage(cmd, &issue, &opts)
		},
	}
}

func (jc *JiraCli) CmdEditMetaUsage(cmd *kingpin.CmdClause, issue *string, opts *GlobalOptions) error {
	if err := jc.GlobalUsage(cmd, opts); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, opts)
	cmd.Arg("ISSUE", "edit metadata for issue id").Required().StringVar(issue)
	return nil
}

// EditMeta will get issue edit metadata and send to "editmeta" template
func (jc *JiraCli) CmdEditMeta(issue string, opts *GlobalOptions) error {
	editMeta, err := jc.GetIssueEditMeta(issue)
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template, editMeta, nil)
}
