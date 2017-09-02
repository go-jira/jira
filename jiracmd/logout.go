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

func CmdLogoutRegistry(fig *figtree.FigTree, o *oreo.Client) *jiracli.CommandRegistryEntry {
	opts := jiracli.GlobalOptions{}
	return &jiracli.CommandRegistryEntry{
		"Deactivate sesssion with Jira server",
		func() error {
			return CmdLogout(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return jiracli.GlobalUsage(cmd, &opts)
		},
	}
}

// CmdLogout will attempt to terminate an active Jira session
func CmdLogout(o *oreo.Client, opts *jiracli.GlobalOptions) error {
	ua := o.WithoutRedirect().WithRetries(0).WithoutCallbacks()
	err := jira.DeleteSession(ua, opts.Endpoint.Value)
	if err == nil {
		fmt.Println(ansi.Color("OK", "green"), "Terminated session for", opts.User)
	} else {
		fmt.Printf("%s Failed to terminate session for %s: %s", ansi.Color("ERROR", "red"), opts.User, err)
	}
	return nil
}
