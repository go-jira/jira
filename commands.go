package jira

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	// "github.com/kr/pretty"
)

func (c *Cli) CmdLogin() error {
	uri := fmt.Sprintf("%s/rest/auth/1/session", c.endpoint)
	for true {
		req, _ := http.NewRequest("GET", uri, nil)
		user := c.opts.User

		fmt.Printf("Enter Password for %s: ", user)
		pwbytes, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
		passwd := string(pwbytes)

		req.SetBasicAuth(user, passwd)
		log.Info("%s %s", req.Method, req.URL.String())
		if resp, err := c.makeRequest(req); err != nil {
			return err
		} else {
			out, _ := httputil.DumpResponse(resp, true)
			log.Debug("%s", out)
			if resp.StatusCode == 403 {
				// probably got this, need to redirect the user to login manually
				// X-Authentication-Denied-Reason: CAPTCHA_CHALLENGE; login-url=https://jira/login.jsp
				if reason := resp.Header.Get("X-Authentication-Denied-Reason"); reason != "" {
					err := fmt.Errorf("Authenticaion Failed: %s", reason)
					log.Error("%s", err)
					return err
				}
				err := fmt.Errorf("Authentication Failed: Unknown Reason")
				log.Error("%s", err)
				return err

			} else if resp.StatusCode == 200 {
				// https://confluence.atlassian.com/display/JIRA043/JIRA+REST+API+%28Alpha%29+Tutorial#JIRARESTAPI%28Alpha%29Tutorial-CAPTCHAs
				// probably bad password, try again
				if reason := resp.Header.Get("X-Seraph-Loginreason"); reason == "AUTHENTICATION_DENIED" {
					log.Warning("Authentication Failed: %s", reason)
					continue
				}
			} else {
				log.Warning("Login failed")
				continue
			}
		}
		return nil
	}
	return nil
}

func (c *Cli) CmdFields() error {
	log.Debug("fields called")
	uri := fmt.Sprintf("%s/rest/api/2/field", c.endpoint)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("fields"), data, nil)
}

func (c *Cli) CmdList() error {
	log.Debug("list called")
	if data, err := c.FindIssues(); err != nil {
		return err
	} else {
		return runTemplate(c.getTemplate("list"), data, nil)
	}
}

func (c *Cli) CmdView(issue string) error {
	log.Debug("view called")
	c.Browse(issue)
	data, err := c.ViewIssue(issue)
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("view"), data, nil)
}

func (c *Cli) CmdEdit(issue string) error {
	log.Debug("edit called")

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/editmeta", c.endpoint, issue)
	editmeta, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	uri = fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
	var issueData map[string]interface{}
	if data, err := responseToJson(c.get(uri)); err != nil {
		return err
	} else {
		issueData = data.(map[string]interface{})
	}

	issueData["meta"] = editmeta.(map[string]interface{})
	issueData["overrides"] = c.opts

	return c.editTemplate(
		c.getTemplate("edit"),
		fmt.Sprintf("%s-edit-", issue),
		issueData,
		func(json string) error {
			if c.opts.DryRun {
				log.Debug("PUT: %s", json)
				log.Debug("Dryrun mode, skipping PUT")
				return nil
			}
			resp, err := c.put(uri, json)
			if err != nil {
				return err
			}

			if resp.StatusCode == 204 {
				c.Browse(issueData["key"].(string))
				if !c.opts.Quiet {
					fmt.Printf("OK %s %s/browse/%s\n", issueData["key"], c.endpoint, issueData["key"])
				}
				return nil
			} else {
				logBuffer := bytes.NewBuffer(make([]byte, 0))
				resp.Write(logBuffer)
				err := fmt.Errorf("Unexpected Response From PUT")
				log.Error("%s:\n%s", err, logBuffer)
				return err
			}
		},
	)
}

func (c *Cli) CmdEditMeta(issue string) error {
	log.Debug("editMeta called")
	c.Browse(issue)
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/editmeta", c.endpoint, issue)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("editmeta"), data, nil)
}

func (c *Cli) CmdTransitionMeta(issue string) error {
	log.Debug("tranisionMeta called")
	c.Browse(issue)
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions?expand=transitions.fields", c.endpoint, issue)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("transmeta"), data, nil)
}

func (c *Cli) CmdIssueTypes() error {
	project := c.opts.Project
	log.Debug("issueTypes called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s", c.endpoint, project)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("issuetypes"), data, nil)
}

func (c *Cli) CmdCreateMeta() error {
	project := c.opts.Project
	issuetype := c.opts.IssueType
	if issuetype == "" {
		issuetype = "Bug"
	}

	log.Debug("createMeta called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&issuetypeNames=%s&expand=projects.issuetypes.fields", c.endpoint, project, issuetype)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	if val, ok := data.(map[string]interface{})["projects"]; ok {
		if len(val.([]interface{})) == 0 {
			err = fmt.Errorf("Project '%s' or issuetype '%s' unknown.  Unable to createmeta.", project, issuetype)
			log.Error("%s", err)
			return err
		}
		if val, ok = val.([]interface{})[0].(map[string]interface{})["issuetypes"]; ok {
			data = val.([]interface{})[0]
		}
	}

	return runTemplate(c.getTemplate("createmeta"), data, nil)
}

func (c *Cli) CmdTransitions(issue string) error {
	log.Debug("Transitions called")
	c.Browse(issue)
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.endpoint, issue)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("transitions"), data, nil)
}

func (c *Cli) CmdCreate() error {
	project := c.opts.Project
	issuetype := c.opts.IssueType
	if issuetype == "" {
		issuetype = "Bug"
	}
	log.Debug("create called")

	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&issuetypeNames=%s&expand=projects.issuetypes.fields", c.endpoint, project, issuetype)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	issueData := make(map[string]interface{})
	issueData["overrides"] = c.opts
	issueData["overrides"].(map[string]interface{})["issuetype"] = issuetype

	if val, ok := data.(map[string]interface{})["projects"]; ok {
		if len(val.([]interface{})) == 0 {
			err = fmt.Errorf("Project '%s' or issuetype '%s' unknown.  Unable to create issue.", project, issuetype)
			log.Error("%s", err)
			return err
		}
		if val, ok = val.([]interface{})[0].(map[string]interface{})["issuetypes"]; ok {
			if len(val.([]interface{})) == 0 {
				err = fmt.Errorf("Project '%s' does not support issuetype '%s'.  Unable to create issue.", project, issuetype)
				log.Error("%s", err)
				return err
			}
			issueData["meta"] = val.([]interface{})[0]
		}
	}

	sanitizedType := strings.ToLower(strings.Replace(issuetype, " ", "", -1))
	return c.editTemplate(
		c.getTemplate(fmt.Sprintf("create-%s", sanitizedType)),
		fmt.Sprintf("create-%s-", sanitizedType),
		issueData,
		func(json string) error {
			log.Debug("JSON: %s", json)
			uri := fmt.Sprintf("%s/rest/api/2/issue", c.endpoint)
			if c.opts.DryRun {
				log.Debug("POST: %s", json)
				log.Debug("Dryrun mode, skipping POST")
				return nil
			}
			resp, err := c.post(uri, json)
			if err != nil {
				return err
			}

			if resp.StatusCode == 201 {
				// response: {"id":"410836","key":"PROJ-238","self":"https://jira/rest/api/2/issue/410836"}
				if json, err := responseToJson(resp, nil); err != nil {
					return err
				} else {
					key := json.(map[string]interface{})["key"].(string)
					link := fmt.Sprintf("%s/browse/%s", c.endpoint, key)
					c.Browse(key)
					c.SaveData(map[string]string{
						"issue": key,
						"link":  link,
					})
					if !c.opts.Quiet {
						fmt.Printf("OK %s %s\n", key, link)
					}
				}
				return nil
			} else {
				logBuffer := bytes.NewBuffer(make([]byte, 0))
				resp.Write(logBuffer)
				err := fmt.Errorf("Unexpected Response From POST")
				log.Error("%s:\n%s", err, logBuffer)
				return err
			}
		},
	)
	return nil
}

func (c *Cli) CmdIssueLinkTypes() error {
	log.Debug("Transitions called")
	uri := fmt.Sprintf("%s/rest/api/2/issueLinkType", c.endpoint)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("issuelinktypes"), data, nil)
}

func (c *Cli) CmdBlocks(blocker string, issue string) error {
	log.Debug("blocks called")

	json, err := jsonEncode(map[string]interface{}{
		"type": map[string]string{
			"name": "Depends", // TODO This is probably not constant across Jira installs
		},
		"inwardIssue": map[string]string{
			"key": issue,
		},
		"outwardIssue": map[string]string{
			"key": blocker,
		},
	})
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issueLink", c.endpoint)
	if c.opts.DryRun {
		log.Debug("POST: %s", json)
		log.Debug("Dryrun mode, skipping POST")
		return nil
	}
	resp, err := c.post(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		c.Browse(issue)
		if !c.opts.Quiet {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From POST")
		log.Error("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

func (c *Cli) CmdDups(duplicate string, issue string) error {
	log.Debug("dups called")

	json, err := jsonEncode(map[string]interface{}{
		"type": map[string]string{
			"name": "Duplicate", // TODO This is probably not constant across Jira installs
		},
		"inwardIssue": map[string]string{
			"key": duplicate,
		},
		"outwardIssue": map[string]string{
			"key": issue,
		},
	})
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issueLink", c.endpoint)
	if c.opts.DryRun {
		log.Debug("POST: %s", json)
		log.Debug("Dryrun mode, skipping POST")
		return nil
	}
	resp, err := c.post(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		c.Browse(issue)
		if !c.opts.Quiet {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From POST")
		log.Error("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

func (c *Cli) CmdWatch(issue string) error {
	watcher := c.opts.Watcher
	if watcher == "" {
		watcher = c.opts.User
	}

	log.Debug("watch called")

	json, err := jsonEncode(watcher)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/watchers", c.endpoint, issue)
	if c.opts.DryRun {
		log.Debug("POST: %s", json)
		log.Debug("Dryrun mode, skipping POST")
		return nil
	}
	resp, err := c.post(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		c.Browse(issue)
		if !c.opts.Quiet {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From POST")
		log.Error("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

func (c *Cli) CmdTransition(issue string, trans string) error {
	log.Debug("transition called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions?expand=transitions.fields", c.endpoint, issue)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	transitions := data.(map[string]interface{})["transitions"].([]interface{})
	var transId, transName string
	var transMeta map[string]interface{}
	found := make([]string, 0, len(transitions))
	for _, transition := range transitions {
		name := transition.(map[string]interface{})["name"].(string)
		id := transition.(map[string]interface{})["id"].(string)
		found = append(found, name)
		if strings.Contains(strings.ToLower(name), strings.ToLower(trans)) {
			transName = name
			transId = id
			transMeta = transition.(map[string]interface{})
		}
	}
	if transId == "" {
		err := fmt.Errorf("Invalid Transition '%s', Available: %s", trans, strings.Join(found, ", "))
		log.Error("%s", err)
		return err
	}

	handlePost := func(json string) error {
		log.Debug("POST: %s", json)
		// os.Exit(0)
		uri = fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.endpoint, issue)
		if c.opts.DryRun {
			log.Debug("POST: %s", json)
			log.Debug("Dryrun mode, skipping POST")
			return nil
		}
		resp, err := c.post(uri, json)
		if err != nil {
			return err
		}
		if resp.StatusCode == 204 {
			c.Browse(issue)
			if !c.opts.Quiet {
				fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
			}
		} else {
			logBuffer := bytes.NewBuffer(make([]byte, 0))
			resp.Write(logBuffer)
			err := fmt.Errorf("Unexpected Response From POST")
			log.Error("%s:\n%s", err, logBuffer)
			return err
		}
		return nil
	}

	uri = fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
	var issueData map[string]interface{}
	if data, err := responseToJson(c.get(uri)); err != nil {
		return err
	} else {
		issueData = data.(map[string]interface{})
	}
	issueData["meta"] = transMeta
	issueData["overrides"] = c.opts
	issueData["transition"] = map[string]interface{}{
		"name": transName,
		"id":   transId,
	}

	return c.editTemplate(
		c.getTemplate("transition"),
		fmt.Sprintf("%s-trans-%s-", issue, trans),
		issueData,
		handlePost,
	)
}

func (c *Cli) CmdComment(issue string) error {
	log.Debug("comment called")

	handlePost := func(json string) error {
		log.Debug("JSON: %s", json)
		uri := fmt.Sprintf("%s/rest/api/2/issue/%s/comment", c.endpoint, issue)
		if c.opts.DryRun {
			log.Debug("POST: %s", json)
			log.Debug("Dryrun mode, skipping POST")
			return nil
		}
		resp, err := c.post(uri, json)
		if err != nil {
			return err
		}

		if resp.StatusCode == 201 {
			c.Browse(issue)
			if !c.opts.Quiet {
				fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
			}
			return nil
		} else {
			logBuffer := bytes.NewBuffer(make([]byte, 0))
			resp.Write(logBuffer)
			err := fmt.Errorf("Unexpected Response From POST")
			log.Error("%s:\n%s", err, logBuffer)
			return err
		}
	}

	if c.opts.Comment != "" {
		json, err := jsonEncode(map[string]interface{}{
			"body": c.opts.Comment,
		})
		if err != nil {
			return err
		}
		return handlePost(json)
	} else {
		return c.editTemplate(
			c.getTemplate("comment"),
			fmt.Sprintf("%s-create-", issue),
			map[string]interface{}{},
			handlePost,
		)
	}
	return nil
}

func (c *Cli) CmdLabels(action string, issue string, labels []string) error {
	log.Debug("label called")

	if action != "add" && action != "remove" && action != "set" {
		return fmt.Errorf("action must be 'add', 'set' or 'remove': %q is invalid", action)
	}

	handlePut := func(json string) error {
		log.Debug("JSON: %s", json)
		uri := fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
		if c.opts.DryRun {
			log.Debug("PUT: %s", json)
			log.Debug("Dryrun mode, skipping POST")
			return nil
		}
		resp, err := c.put(uri, json)
		if err != nil {
			return err
		}

		if resp.StatusCode == 204 {
			c.Browse(issue)
			if !c.opts.Quiet {
				fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
			}
			return nil
		} else {
			logBuffer := bytes.NewBuffer(make([]byte, 0))
			resp.Write(logBuffer)
			err := fmt.Errorf("Unexpected Response From PUT")
			log.Error("%s:\n%s", err, logBuffer)
			return err
		}
	}

	var labels_json string
	var err error
	if action == "set" {
		labelsActions := make([]map[string][]string, 1)
		labelsActions[0] = map[string][]string{
			"set": labels,
		}
		labels_json, err = jsonEncode(map[string]interface{}{
			"labels": labelsActions,
		})
	} else {
		labelsActions := make([]map[string]string, len(labels))
		for i, label := range labels {
			labelActionMap := map[string]string{
				action: label,
			}
			labelsActions[i] = labelActionMap
		}
		labels_json, err = jsonEncode(map[string]interface{}{
			"labels": labelsActions,
		})
	}
	if err != nil {
		return err
	}
	json := fmt.Sprintf("{ \"update\": %s }", labels_json)
	return handlePut(json)

}

func (c *Cli) CmdAssign(issue string, user string) error {
	log.Debug("assign called")

	json, err := jsonEncode(map[string]interface{}{
		"name": user,
	})
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/assignee", c.endpoint, issue)
	if c.opts.DryRun {
		log.Debug("PUT: %s", json)
		log.Debug("Dryrun mode, skipping PUT")
		return nil
	}
	resp, err := c.put(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		c.Browse(issue)
		if !c.opts.Quiet {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From PUT")
		log.Error("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

func (c *Cli) CmdExportTemplates() error {
	dir := c.opts.Directory
	if err := mkdir(dir); err != nil {
		return err
	}

	for name, template := range all_templates {
		if (c.opts.Template != "") && (c.opts.Template != name) {
			continue
		}

		templateFile := fmt.Sprintf("%s/%s", dir, name)
		if _, err := os.Stat(templateFile); err == nil {
			log.Warning("Skipping %s, already exists", templateFile)
			continue
		}
		if fh, err := os.OpenFile(templateFile, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
			log.Error("Failed to open %s for writing: %s", templateFile, err)
			return err
		} else {
			defer fh.Close()
			log.Notice("Creating %s", templateFile)
			fh.Write([]byte(template))
		}
	}
	return nil
}

func (c *Cli) CmdRequest(uri, content string) (err error) {
	log.Debug("request called")

	if !strings.HasPrefix(uri, "http") {
		uri = fmt.Sprintf("%s%s", c.endpoint, uri)
	}

	method := strings.ToUpper(c.opts.Method)
	var data interface{}
	if method == "GET" {
		data, err = responseToJson(c.get(uri))
	} else if method == "POST" {
		data, err = responseToJson(c.post(uri, content))
	} else if method == "PUT" {
		data, err = responseToJson(c.put(uri, content))
	}
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("request"), data, nil)
}
