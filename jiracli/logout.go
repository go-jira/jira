package jiracli

import (
	"fmt"

	"github.com/mgutz/ansi"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func (jc *JiraCli) CmdLogoutRegistry() *CommandRegistryEntry {
	opts := GlobalOptions{}
	return &CommandRegistryEntry{
		"Deactivate sesssion with Jira server",
		func() error {
			return jc.CmdLogout(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.GlobalUsage(cmd, &opts)
		},
	}
}

// CmdLogout will attempt to terminate an active Jira session
func (jc *JiraCli) CmdLogout(opts *GlobalOptions) error {
	jc.UA = jc.oreoAgent.WithoutRedirect().WithRetries(0)
	err := jc.DeleteSession()
	if err == nil {
		fmt.Println(ansi.Color("OK", "green"), "Terminated session for", opts.User)
	} else {
		fmt.Printf("%s Failed to terminate session for %s: %s", ansi.Color("ERROR", "red"), opts.User, err)
	}
	return nil
}
