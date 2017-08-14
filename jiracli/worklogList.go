package jiracli

import kingpin "gopkg.in/alecthomas/kingpin.v2"

func (jc *JiraCli) CmdWorklogListRegistry() *CommandRegistryEntry {
	issue := ""
	opts := GlobalOptions{
		Template: "worklogs",
	}
	return &CommandRegistryEntry{
		"Prints the worklog data for given issue",
		func() error {
			return jc.CmdWorklogList(issue, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdWorklogListUsage(cmd, &issue, &opts)
		},
	}
}

func (jc *JiraCli) CmdWorklogListUsage(cmd *kingpin.CmdClause, issue *string, opts *GlobalOptions) error {
	if err := jc.GlobalUsage(cmd, opts); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, opts)
	cmd.Arg("ISSUE", "issue id to fetch worklogs").Required().StringVar(issue)
	return nil
}

// // CmdWorklogList will get worklog data for given issue and sent to the "worklogs" template
func (jc *JiraCli) CmdWorklogList(issue string, opts *GlobalOptions) error {
	data, err := jc.GetIssueWorklog(issue)
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template, data, nil)
}
