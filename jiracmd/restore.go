package jiracmd

import (
	"fmt"
	"strings"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type RestoreOptions struct {
	jiracli.CommonOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Project               string `yaml:"project,omitempty" json:"project,omitempty"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdRestoreRegistry() *jiracli.CommandRegistryEntry {
	opts := AssignOptions{}

	return &jiracli.CommandRegistryEntry{
		"Restore issue",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdRestoreUsage(cmd, &opts)
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			opts.Issue = jiracli.FormatIssue(opts.Issue, opts.Project)
			return CmdRestore(o, globals, &opts)
		},
	}
}

func CmdRestoreUsage(cmd *kingpin.CmdClause, opts *AssignOptions) error {
	jiracli.BrowseUsage(cmd, &opts.CommonOptions)
	cmd.Arg("ISSUE", "issue to restore").Required().StringVar(&opts.Issue)
	return nil
}

// CmdRestore will restore an issue
func CmdRestore(o *oreo.Client, globals *jiracli.GlobalOptions, opts *AssignOptions) error {
	if globals.JiraDeploymentType.Value == "" {
		serverInfo, err := jira.ServerInfo(o, globals.Endpoint.Value)
		if err != nil {
			return err
		}
		globals.JiraDeploymentType.Value = strings.ToLower(serverInfo.DeploymentType)
	}

	err := jira.IssueRestore(o, globals.Endpoint.Value, opts.Issue)
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
