package jiracli

import (
	"strings"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type CommandRegistryEntry struct {
	Help        string
	ExecuteFunc func() error
	UsageFunc   func(*kingpin.CmdClause) error
}

type CommandRegistry struct {
	Command string
	Aliases []string
	Entry   *CommandRegistryEntry
	Default bool
}

// either kingpin.Application or kingpin.CmdClause fit this interface
type kingpinAppOrCommand interface {
	Command(string, string) *kingpin.CmdClause
	GetCommand(string) *kingpin.CmdClause
}

func (jc *JiraCli) Register(app *kingpin.Application, reg []CommandRegistry) {
	for _, command := range reg {
		copy := command
		commandFields := strings.Fields(copy.Command)
		var appOrCmd kingpinAppOrCommand = app
		if len(commandFields) > 1 {
			for _, name := range commandFields[0 : len(commandFields)-1] {
				tmp := appOrCmd.GetCommand(name)
				if tmp == nil {
					tmp = appOrCmd.Command(name, "")
				}
				appOrCmd = tmp
			}
		}

		cmd := appOrCmd.Command(commandFields[len(commandFields)-1], copy.Entry.Help)
		for _, alias := range copy.Aliases {
			cmd = cmd.Alias(alias)
		}
		if copy.Default {
			cmd = cmd.Default()
		}
		if copy.Entry.UsageFunc != nil {
			copy.Entry.UsageFunc(cmd)
		}

		cmd.Action(
			func(_ *kingpin.ParseContext) error {
				return copy.Entry.ExecuteFunc()
			},
		)
	}
}

// // CmdRequest will use the given uri to make a request and potentially send provided content.
// func (c *Cli) CmdRequest(uri, content string) (err error) {
// 	log.Debugf("request called")

// 	if !strings.HasPrefix(uri, "http") {
// 		uri = fmt.Sprintf("%s%s", c.endpoint, uri)
// 	}

// 	method := strings.ToUpper(c.opts["method"].(string))
// 	var data interface{}
// 	if method == "GET" {
// 		data, err = responseToJSON(c.get(uri))
// 	} else if method == "POST" {
// 		data, err = responseToJSON(c.post(uri, content))
// 	} else if method == "PUT" {
// 		data, err = responseToJSON(c.put(uri, content))
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return runTemplate(c.getTemplate("request"), data, nil)
// }
