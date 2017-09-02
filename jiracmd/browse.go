package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/pkg/browser"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type BrowseOptions struct {
	jiracli.GlobalOptions `yaml:",inline" json:",inline" figtree:",inline"`
	Issue                 string `yaml:"issue,omitempty" json:"issue,omitempty"`
}

func CmdBrowseRegistry(fig *figtree.FigTree) *jiracli.CommandRegistryEntry {
	opts := BrowseOptions{}

	return &jiracli.CommandRegistryEntry{
		"Open issue in browser",
		func() error {
			return CmdBrowse(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return CmdBrowseUsage(cmd, &opts)
		},
	}
}

func CmdBrowseUsage(cmd *kingpin.CmdClause, opts *BrowseOptions) error {
	if err := jiracli.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	cmd.Arg("ISSUE", "Issue to browse to").Required().StringVar(&opts.Issue)

	return nil
}

// CmdBrowse open the default system browser to the provided issue
func CmdBrowse(opts *BrowseOptions) error {
	return browser.OpenURL(fmt.Sprintf("%s/browse/%s", opts.Endpoint.Value, opts.Issue))
}
