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

// func NewCommand(app kingpinAppOrCommand, name string, entry *CommandRegistryEntry) *kingpin.CmdClause {
// 	returnapp.Command(name, entry.Help)
// }

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

// // CmdIssueTypes will send issue 'create' metadata to the 'issuetypes'
// func (c *Cli) CmdIssueTypes() error {
// 	project := c.opts["project"].(string)
// 	log.Debugf("issueTypes called")
// 	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s", c.endpoint, project)
// 	data, err := responseToJSON(c.get(uri))
// 	if err != nil {
// 		return err
// 	}

// 	return runTemplate(c.getTemplate("issuetypes"), data, nil)
// }

// func (c *Cli) defaultIssueType() string {
// 	project := c.opts["project"].(string)
// 	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s", c.endpoint, project)
// 	data, _ := responseToJSON(c.get(uri))
// 	issueTypeNames := make(map[string]bool)

// 	if data, ok := data.(map[string]interface{}); ok {
// 		if projects, ok := data["projects"].([]interface{}); ok {
// 			for _, project := range projects {
// 				if project, ok := project.(map[string]interface{}); ok {
// 					if issuetypes, ok := project["issuetypes"].([]interface{}); ok {
// 						if len(issuetypes) > 0 {
// 							for _, issuetype := range issuetypes {
// 								issueTypeNames[issuetype.(map[string]interface{})["name"].(string)] = true
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// 	if _, ok := issueTypeNames["Bug"]; ok {
// 		return "Bug"
// 	} else if _, ok := issueTypeNames["Task"]; ok {
// 		return "Task"
// 	}
// 	return ""
// }

// // CmdExportTemplates will export the default templates to the template directory.
// func (c *Cli) CmdExportTemplates() error {
// 	dir := c.opts["directory"].(string)
// 	if err := mkdir(dir); err != nil {
// 		return err
// 	}

// 	for name, template := range allTemplates {
// 		if wanted, ok := c.opts["template"]; ok && wanted != name {
// 			continue
// 		}
// 		templateFile := fmt.Sprintf("%s/%s", dir, name)
// 		if _, err := os.Stat(templateFile); err == nil {
// 			log.Warning("Skipping %s, already exists", templateFile)
// 			continue
// 		}
// 		fh, err := os.OpenFile(templateFile, os.O_WRONLY|os.O_CREATE, 0644)
// 		if err != nil {
// 			log.Errorf("Failed to open %s for writing: %s", templateFile, err)
// 			return err
// 		}
// 		defer fh.Close()
// 		log.Noticef("Creating %s", templateFile)
// 		fh.Write([]byte(template))
// 	}
// 	return nil
// }

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
