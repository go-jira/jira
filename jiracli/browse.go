package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/pkg/browser"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type BrowseOptions struct {
	GlobalOptions `yaml:",inline" figtree:",inline"`
	Issue         string
}

func CmdBrowseRegistry(fig *figtree.FigTree) *CommandRegistryEntry {
	opts := BrowseOptions{}

	return &CommandRegistryEntry{
		"Open issue in browser",
		func() error {
			return CmdBrowse(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return CmdBrowseUsage(cmd, &opts)
		},
	}
}

func CmdBrowseUsage(cmd *kingpin.CmdClause, opts *BrowseOptions) error {
	if err := GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	cmd.Arg("ISSUE", "Issue to browse to").Required().StringVar(&opts.Issue)

	return nil
}

// CmdBrowse open the default system browser to the provided issue
func CmdBrowse(opts *BrowseOptions) error {
	return browser.OpenURL(fmt.Sprintf("%s/browse/%s", opts.Endpoint.Value, opts.Issue))
}
