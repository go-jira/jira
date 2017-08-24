package jiracli

import kingpin "gopkg.in/alecthomas/kingpin.v2"

type EditMetaOptions struct {
	GlobalOptions
	Issue string
}

func (jc *JiraCli) CmdEditMetaRegistry() *CommandRegistryEntry {

	opts := EditMetaOptions{
		GlobalOptions: GlobalOptions{
			Template: "editmeta",
		},
	}

	return &CommandRegistryEntry{
		"View 'edit' metadata",
		func() error {
			return jc.CmdEditMeta(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdEditMetaUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdEditMetaUsage(cmd *kingpin.CmdClause, opts *EditMetaOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.BrowseUsage(cmd, &opts.GlobalOptions)
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "edit metadata for issue id").Required().StringVar(&opts.Issue)
	return nil
}

// EditMeta will get issue edit metadata and send to "editmeta" template
func (jc *JiraCli) CmdEditMeta(opts *EditMetaOptions) error {
	editMeta, err := jc.GetIssueEditMeta(opts.Issue)
	if err != nil {
		return err
	}
	if err := jc.runTemplate(opts.Template, editMeta, nil); err != nil {
		return err
	}
	if opts.Browse {
		return jc.CmdBrowse(&BrowseOptions{opts.GlobalOptions, opts.Issue})
	}
	return nil
}
