package jiracmd

import (
	"sort"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/go-jira/jira/jiradata"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type AttachListOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdAttachListRegistry() *jiracli.CommandRegistryEntry {
	opts := AttachListOptions{
		CommonOptions: jiracli.CommonOptions{
			Template: figtree.NewStringOption("attach-list"),
		},
	}

	return &jiracli.CommandRegistryEntry{
		"Prints attachment details for issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdAttachListUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdAttachList(o, globals, &opts)
		},
	}
}

func CmdAttachListUsage(cmd *kingpin.CmdClause, opts *AttachListOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	jiracli.TemplateUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "Issue id to lookup attachments").Required().StringVar(&opts.Issue)
	return nil
}

func CmdAttachList(o *oreo.Client, globals *jiracli.GlobalOptions, opts *AttachListOptions) error {
	data, err := jira.GetIssue(o, globals.Endpoint.Value, opts.Issue, nil)
	if err != nil {
		return err
	}

	// need to conver the interface{} "attachment" field to an actual
	// ListOfAttachment object so we can sort it
	var attachments jiradata.ListOfAttachment
	err = jiracli.ConvertType(data.Fields["attachment"], &attachments)
	if err != nil {
		return err
	}
	sort.Sort(&attachments)

	if err := opts.PrintTemplate(attachments); err != nil {
		return err
	}
	if opts.Browse.Value {
		return CmdBrowse(globals, opts.Issue)
	}
	return nil
}
