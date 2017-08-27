package jiracli

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/mgutz/ansi"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdLogoutRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := GlobalOptions{}
	return &CommandRegistryEntry{
		"Deactivate sesssion with Jira server",
		func() error {
			return CmdLogout(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return GlobalUsage(cmd, &opts)
		},
	}
}

// CmdLogout will attempt to terminate an active Jira session
func CmdLogout(o *oreo.Client, opts *GlobalOptions) error {
	ua := o.WithoutRedirect().WithRetries(0)
	err := jira.DeleteSession(ua, opts.Endpoint.Value)
	if err == nil {
		fmt.Println(ansi.Color("OK", "green"), "Terminated session for", opts.User)
	} else {
		fmt.Printf("%s Failed to terminate session for %s: %s", ansi.Color("ERROR", "red"), opts.User, err)
	}
	return nil
}
