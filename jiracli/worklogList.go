package jiracli

import kingpin "gopkg.in/alecthomas/kingpin.v2"

type WorklogListOptions struct {
	GlobalOptions
	Issue string
}

func (jc *JiraCli) CmdWorklogListRegistry() *CommandRegistryEntry {
	opts := WorklogListOptions{
		GlobalOptions: GlobalOptions{
			Template: "worklogs",
		},
	}
	return &CommandRegistryEntry{
		"Prints the worklog data for given issue",
		func() error {
			return jc.CmdWorklogList(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdWorklogListUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdWorklogListUsage(cmd *kingpin.CmdClause, opts *WorklogListOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	jc.TemplateUsage(cmd, &opts.GlobalOptions)
	cmd.Arg("ISSUE", "issue id to fetch worklogs").Required().StringVar(&opts.Issue)
	return nil
}

// // CmdWorklogList will get worklog data for given issue and sent to the "worklogs" template
func (jc *JiraCli) CmdWorklogList(opts *WorklogListOptions) error {
	data, err := jc.GetIssueWorklog(opts.Issue)
	if err != nil {
		return err
	}
	return jc.runTemplate(opts.Template, data, nil)
}
