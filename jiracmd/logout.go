package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/mgutz/ansi"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdLogoutRegistry() *jiracli.CommandRegistryEntry {
	opts := jiracli.CommonOptions{}
	return &jiracli.CommandRegistryEntry{
		"Deactivate sesssion with Jira server",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return nil
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdLogout(o, globals, &opts)
		},
	}
}

// CmdLogout will attempt to terminate an active Jira session
func CmdLogout(o *oreo.Client, globals *jiracli.GlobalOptions, opts *jiracli.CommonOptions) error {
	ua := o.WithoutRedirect().WithRetries(0).WithoutCallbacks()
	err := jira.DeleteSession(ua, globals.Endpoint.Value)
	if err == nil {
		fmt.Println(ansi.Color("OK", "green"), "Terminated session for", globals.User)
	} else {
		fmt.Printf("%s Failed to terminate session for %s: %s", ansi.Color("ERROR", "red"), globals.User, err)
	}
	return nil
}
