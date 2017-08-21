package jiracli

import (
	"fmt"

	"github.com/pkg/browser"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type BrowseOptions struct {
	GlobalOptions
	Issue string
}

func (jc *JiraCli) CmdBrowseRegistry() *CommandRegistryEntry {
	opts := BrowseOptions{}

	return &CommandRegistryEntry{
		"Open issue in browser",
		func() error {
			return jc.CmdBrowse(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.CmdBrowseUsage(cmd, &opts)
		},
	}
}

func (jc *JiraCli) CmdBrowseUsage(cmd *kingpin.CmdClause, opts *BrowseOptions) error {
	if err := jc.GlobalUsage(cmd, &opts.GlobalOptions); err != nil {
		return err
	}
	cmd.Arg("ISSUE", "Issue to browse to").Required().StringVar(&opts.Issue)

	return nil
}

// CmdBrowse open the default system browser to the provided issue
func (jc *JiraCli) CmdBrowse(opts *BrowseOptions) error {
	return browser.OpenURL(fmt.Sprintf("%s/browse/%s", jc.Endpoint, opts.Issue))
}
