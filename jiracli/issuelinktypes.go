package jiracli

import kingpin "gopkg.in/alecthomas/kingpin.v2"

func (jc *JiraCli) CmdIssueLinkTypesRegistry() *CommandRegistryEntry {
	opts := GlobalOptions{
		Template: "issuelinktypes",
	}

	return &CommandRegistryEntry{
		"Show the issue link types",
		func() error {
			return jc.CmdIssueLinkTypes(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdIssueLinkTypesUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdIssueLinkTypesUsage(cmd *kingpin.CmdClause, opts *GlobalOptions) error {
	if err := jc.GlobalUsage(cmd, opts); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, opts)
	return nil
}

// CmdIssueLinkTypes will get issue link type data and send to "issuelinktypes" template
func (jc *JiraCli) CmdIssueLinkTypes(opts *GlobalOptions) error {
	data, err := jc.GetIssueLinkTypes()
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template, data, nil)
}
