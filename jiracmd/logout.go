package jiracmd

import (
	"fmt"
	"os"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	"github.com/mgutz/ansi"
	"golang.org/x/crypto/ssh/terminal"
	survey "gopkg.in/AlecAivazis/survey.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdLogoutRegistry() *jiracli.CommandRegistryEntry {
	opts := jiracli.CommonOptions{}
	return &jiracli.CommandRegistryEntry{
		"Deactivate session with Jira server",
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
	if globals.AuthMethodIsToken() {
		log.Noticef("No need to logout when using api-token or bearer-token authentication method")
		if globals.GetPass() != "" && terminal.IsTerminal(int(os.Stdin.Fd())) && terminal.IsTerminal(int(os.Stdout.Fd())) {
			delete := false
			err := survey.AskOne(
				&survey.Confirm{
					Message: fmt.Sprintf("Delete token from password provider [%s]: ", globals.PasswordSource),
					Default: false,
				},
				&delete,
				nil,
			)
			if err != nil {
				log.Errorf("%s", err)
				panic(jiracli.Exit{Code: 1})
			}
			if delete {
				globals.SetPass("")
			}
		}
		return nil
	}
	ua := o.WithoutRedirect().WithRetries(0).WithoutCallbacks()
	err := jira.DeleteSession(ua, globals.Endpoint.Value)
	if err == nil {
		if !globals.Quiet.Value {
			fmt.Println(ansi.Color("OK", "green"), "Terminated session for", globals.User)
		}
	} else {
		fmt.Printf("%s Failed to terminate session for %s: %s", ansi.Color("ERROR", "red"), globals.User, err)
	}
	return nil
}
