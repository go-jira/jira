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

// // CmdComponents sends component data for given project and sends to the "components" template
// func (c *Cli) CmdComponents(project string) error {
// 	log.Debugf("Components called")
// 	uri := fmt.Sprintf("%s/rest/api/2/project/%s/components", c.endpoint, project)
// 	data, err := responseToJSON(c.get(uri))
// 	if err != nil {
// 		return err
// 	}
// 	return runTemplate(c.getTemplate("components"), data, nil)
// }

// // CmdWatch will add the given watcher to the issue (or remove the watcher
// // given the 'remove' flag)
// func (c *Cli) CmdWatch(issue string, watcher string, remove bool) error {
// 	log.Debugf("watch called: watcher: %q, remove: %n", watcher, remove)

// 	var uri string
// 	json, err := jsonEncode(watcher)
// 	if err != nil {
// 		return err
// 	}

// 	if c.getOptBool("dryrun", false) {
// 		if !remove {
// 			log.Debugf("POST: %s", json)
// 			log.Debugf("Dryrun mode, skipping POST")
// 		} else {
// 			log.Debugf("DELETE: %s", watcher)
// 			log.Debugf("Dryrun mode, skipping POST")
// 		}
// 		return nil
// 	}

// 	var resp *http.Response
// 	if !remove {
// 		uri = fmt.Sprintf("%s/rest/api/2/issue/%s/watchers", c.endpoint, issue)
// 		resp, err = c.post(uri, json)
// 	} else {
// 		uri = fmt.Sprintf("%s/rest/api/2/issue/%s/watchers?username=%s", c.endpoint, issue, watcher)
// 		resp, err = c.delete(uri)
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	if resp.StatusCode == 204 {
// 		c.Browse(issue)
// 		if !c.GetOptBool("quiet", false) {
// 			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
// 		}
// 	} else {
// 		logBuffer := bytes.NewBuffer(make([]byte, 0))
// 		resp.Write(logBuffer)
// 		if !remove {
// 			err = fmt.Errorf("Unexpected Response From POST")
// 		} else {
// 			err = fmt.Errorf("Unexpected Response From DELETE")
// 		}
// 		log.Errorf("%s:\n%s", err, logBuffer)
// 		return err
// 	}
// 	return nil
// }

// // CmdComment will open up editor with "comment" template and submit
// // YAML output to jira
// func (c *Cli) CmdComment(issue string) error {
// 	log.Debugf("comment called")

// 	handlePost := func(json string) error {
// 		uri := fmt.Sprintf("%s/rest/api/2/issue/%s/comment", c.endpoint, issue)
// 		if c.getOptBool("dryrun", false) {
// 			log.Debugf("POST: %s", json)
// 			log.Debugf("Dryrun mode, skipping POST")
// 			return nil
// 		}
// 		resp, err := c.post(uri, json)
// 		if err != nil {
// 			return err
// 		}

// 		if resp.StatusCode == 201 {
// 			c.Browse(issue)
// 			if !c.GetOptBool("quiet", false) {
// 				fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
// 			}
// 			return nil
// 		}
// 		logBuffer := bytes.NewBuffer(make([]byte, 0))
// 		resp.Write(logBuffer)
// 		err = fmt.Errorf("Unexpected Response From POST")
// 		log.Errorf("%s:\n%s", err, logBuffer)
// 		return err
// 	}

// 	if comment, ok := c.opts["comment"]; ok && comment != "" {
// 		json, err := jsonEncode(map[string]interface{}{
// 			"body": comment,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return handlePost(json)
// 	}
// 	return c.editTemplate(
// 		c.getTemplate("comment"),
// 		fmt.Sprintf("%s-create-", issue),
// 		map[string]interface{}{},
// 		handlePost,
// 	)
// }

// // CmdComponent will add a new component to given project
// func (c *Cli) CmdComponent(action string, project string, name string, desc string, lead string) error {
// 	log.Debugf("component called")

// 	switch action {
// 	case "add":
// 	default:
// 		return fmt.Errorf("CmdComponent: %q is not a valid action", action)
// 	}

// 	json, err := jsonEncode(map[string]interface{}{
// 		"name":         name,
// 		"description":  desc,
// 		"leadUserName": lead,
// 		"project":      project,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	uri := fmt.Sprintf("%s/rest/api/2/component", c.endpoint)
// 	if c.getOptBool("dryrun", false) {
// 		log.Debugf("POST: %s", json)
// 		log.Debugf("Dryrun mode, skipping POST")
// 		return nil
// 	}
// 	resp, err := c.post(uri, json)
// 	if err != nil {
// 		return err
// 	}
// 	if resp.StatusCode == 201 {
// 		if !c.GetOptBool("quiet", false) {
// 			fmt.Printf("OK %s %s\n", project, name)
// 		}
// 	} else {
// 		logBuffer := bytes.NewBuffer(make([]byte, 0))
// 		resp.Write(logBuffer)
// 		err := fmt.Errorf("Unexpected Response From POST")
// 		log.Errorf("%s:\n%s", err, logBuffer)
// 		return err
// 	}
// 	return nil
// }

// // CmdLabels will add, remove or set labels on a given issue
// func (c *Cli) CmdLabels(action string, issue string, labels []string) error {
// 	log.Debugf("label called")

// 	if action != "add" && action != "remove" && action != "set" {
// 		return fmt.Errorf("action must be 'add', 'set' or 'remove': %q is invalid", action)
// 	}

// 	handlePut := func(json string) error {
// 		uri := fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
// 		if c.getOptBool("dryrun", false) {
// 			log.Debugf("PUT: %s", json)
// 			log.Debugf("Dryrun mode, skipping POST")
// 			return nil
// 		}
// 		resp, err := c.put(uri, json)
// 		if err != nil {
// 			return err
// 		}

// 		if resp.StatusCode == 204 {
// 			c.Browse(issue)
// 			if !c.GetOptBool("quiet", false) {
// 				fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
// 			}
// 			return nil
// 		}
// 		logBuffer := bytes.NewBuffer(make([]byte, 0))
// 		resp.Write(logBuffer)
// 		err = fmt.Errorf("Unexpected Response From PUT")
// 		log.Errorf("%s:\n%s", err, logBuffer)
// 		return err
// 	}

// 	var labelsJSON string
// 	var err error
// 	if action == "set" {
// 		labelsActions := make([]map[string][]string, 1)
// 		labelsActions[0] = map[string][]string{
// 			"set": labels,
// 		}
// 		labelsJSON, err = jsonEncode(map[string]interface{}{
// 			"labels": labelsActions,
// 		})
// 	} else {
// 		labelsActions := make([]map[string]string, len(labels))
// 		for i, label := range labels {
// 			labelActionMap := map[string]string{
// 				action: label,
// 			}
// 			labelsActions[i] = labelActionMap
// 		}
// 		labelsJSON, err = jsonEncode(map[string]interface{}{
// 			"labels": labelsActions,
// 		})
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	json := fmt.Sprintf("{ \"update\": %s }", labelsJSON)
// 	return handlePut(json)

// }

// // CmdAssign will assign the given user to be the owner of the given issue
// func (c *Cli) CmdAssign(issue string, user string) error {
// 	log.Debugf("assign called")

// 	var userVal interface{} = user
// 	// https://docs.atlassian.com/jira/REST/cloud/#api/2/issue-assign
// 	// If the name is "-1" automatic assignee is used. A null name will remove the assignee.
// 	if user == "" {
// 		userVal = nil
// 	}
// 	if c.GetOptBool("default", false) {
// 		userVal = "-1"
// 	}

// 	json, err := jsonEncode(map[string]interface{}{
// 		"name": userVal,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/assignee", c.endpoint, issue)
// 	if c.getOptBool("dryrun", false) {
// 		log.Debugf("PUT: %s", json)
// 		log.Debugf("Dryrun mode, skipping PUT")
// 		return nil
// 	}
// 	resp, err := c.put(uri, json)
// 	if err != nil {
// 		return err
// 	}
// 	if resp.StatusCode == 204 {
// 		c.Browse(issue)
// 		if !c.GetOptBool("quiet", false) {
// 			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
// 		}
// 	} else {
// 		logBuffer := bytes.NewBuffer(make([]byte, 0))
// 		resp.Write(logBuffer)
// 		err := fmt.Errorf("Unexpected Response From PUT")
// 		log.Errorf("%s:\n%s", err, logBuffer)
// 		return err
// 	}
// 	return nil
// }

// func (c *Cli) CmdUnassign(issue string) error {
// 	return c.CmdAssign(issue, "")
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
