package jiracmd

import (
	"fmt"

	"github.com/coryb/figtree"
	"github.com/eroshan/oreo"
	"github.com/go-jira/jira"
	"github.com/go-jira/jira/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/coryb/yaml.v2"
)

func CmdSessionRegistry() *jiracli.CommandRegistryEntry {
	opts := jiracli.CommonOptions{}
	return &jiracli.CommandRegistryEntry{
		"Attempt to login into jira server",
		func(fig *figtree.FigTree, cmd *kingpin.CmdClause) error {
			jiracli.LoadConfigs(cmd, fig, &opts)
			return nil
		},
		func(o *oreo.Client, globals *jiracli.GlobalOptions) error {
			return CmdSession(o, globals, &opts)
		},
	}
}

// CmdSession will attempt to login into jira server
func CmdSession(o *oreo.Client, globals *jiracli.GlobalOptions, opts *jiracli.CommonOptions) error {
	ua := o.WithoutRedirect().WithRetries(0).WithoutPostCallbacks()
	session, err := jira.GetSession(ua, globals.Endpoint.Value)
	var output []byte
	if err != nil {
		defer panic(jiracli.Exit{1})
		output, err = yaml.Marshal(err)
		if err != nil {
			return err
		}
	} else {
		output, err = yaml.Marshal(session)
		if err != nil {
			return err
		}
	}
	fmt.Print(string(output))
	return nil
}
