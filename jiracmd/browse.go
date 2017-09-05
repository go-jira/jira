package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/pkg/browser"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdBrowseRegistry() *jiracli.CommandRegistryEntry {
	issue := ""

	return &jiracli.CommandRegistryEntry{
		"Open issue in browser",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			cmd.Arg("ISSUE", "Issue to browse to").Required().StringVar(&issue)
			return nil
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdBrowse(globals, issue)
		},
	}
}

// CmdBrowse open the default system browser to the provided issue
func CmdBrowse(globals *jiracli.GlobalOptions, issue string) error {
	return browser.OpenURL(fmt.Sprintf("%s/browse/%s", globals.Endpoint.Value, issue))
}
