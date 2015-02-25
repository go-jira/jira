package cli

import (
	"bytes"
	"code.google.com/p/gopass"
	"fmt"
	"net/http"
	"os"
	"strings"
	// "github.com/kr/pretty"
)

func (c *Cli) CmdLogin() error {
	uri := fmt.Sprintf("%s/rest/auth/1/session", c.endpoint)
	for true {
		req, _ := http.NewRequest("GET", uri, nil)
		user, _ := c.opts["user"]

		prompt := fmt.Sprintf("Enter Password for %s: ", user)
		passwd, _ := gopass.GetPass(prompt)

		req.SetBasicAuth(user, passwd)
		if resp, err := c.makeRequest(req); err != nil {
			return err
		} else {
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

	var query string
	var ok bool
	// project = BAKERY and status not in (Resolved, Closed)
	if query, ok = c.opts["query"]; !ok {
		qbuff := bytes.NewBufferString("resolution = unresolved")
		if project, ok := c.opts["project"]; !ok {
			err := fmt.Errorf("Missing required arguments, either 'query' or 'project' are required")
			log.Error("%s", err)
			return err
		} else {
			qbuff.WriteString(fmt.Sprintf(" AND project = '%s'", project))
		}

		if component, ok := c.opts["component"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND component = '%s'", component))
		}

		if assignee, ok := c.opts["assignee"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND assignee = '%s'", assignee))
		}

		if issuetype, ok := c.opts["issuetype"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND issuetype = '%s'", issuetype))
		}

		if watcher, ok := c.opts["watcher"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND watcher = '%s'", watcher))
		}

		if reporter, ok := c.opts["reporter"]; ok {
			qbuff.WriteString(fmt.Sprintf(" AND reporter = '%s'", reporter))
		}

		if sort, ok := c.opts["sort"]; ok && sort != "" {
			qbuff.WriteString(fmt.Sprintf(" ORDER BY %s", sort ))
		}

		query = qbuff.String()
	}

	fields := make([]string, 0)
	if qf, ok := c.opts["queryfields"]; ok {
		fields = strings.Split(qf, ",")
	} else {
		fields = append(fields, "summary")
	}

	json, err := jsonEncode(map[string]interface{}{
		"jql":        query,
		"startAt":    "0",
		"maxResults": "500",
		"fields":     fields,
	})
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/search", c.endpoint)
	data, err := responseToJson(c.post(uri, json))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("list"), data, nil)
}

func (c *Cli) CmdView(issue string) error {
	log.Debug("view called")
	c.Browse(issue)
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
	data, err := responseToJson(c.get(uri))
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
			resp, err := c.put(uri, json)
			if err != nil {
				return err
			}

			if resp.StatusCode == 204 {
				c.Browse(issueData["key"].(string))
				fmt.Printf("OK %s %s/browse/%s\n", issueData["key"], c.endpoint, issueData["key"])
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

func (c *Cli) CmdIssueTypes(project string) error {
	log.Debug("issueTypes called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s", c.endpoint, project)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("issuetypes"), data, nil)
}

func (c *Cli) CmdCreateMeta(project string, issuetype string) error {
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

func (c *Cli) CmdCreate(project string, issuetype string) error {
	log.Debug("create called")

	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&issuetypeNames=%s&expand=projects.issuetypes.fields", c.endpoint, project, issuetype)
	data, err := responseToJson(c.get(uri))
	if err != nil {
		return err
	}

	issueData := make(map[string]interface{})
	issueData["overrides"] = c.opts
	issueData["overrides"].(map[string]string)["issuetype"] = issuetype

	if val, ok := data.(map[string]interface{})["projects"]; ok {
		if len(val.([]interface{})) == 0 {
			err = fmt.Errorf("Project '%s' or issuetype '%s' unknown.  Unable to create issue.", project, issuetype)
			log.Error("%s", err)
			return err
		}
		if val, ok = val.([]interface{})[0].(map[string]interface{})["issuetypes"]; ok {
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
			resp, err := c.post(uri, json)
			if err != nil {
				return err
			}

			if resp.StatusCode == 201 {
				// response: {"id":"410836","key":"PROJ-238","self":"https://jira/rest/api/2/issue/410836"}
				if json, err := responseToJson(resp, nil); err != nil {
					return err
				} else {
					key := json.(map[string]interface{})["key"]
					c.Browse(key.(string))
					fmt.Printf("OK %s %s/browse/%s\n", key, c.endpoint, key)

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
	resp, err := c.post(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		c.Browse(issue)
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
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
	resp, err := c.post(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		c.Browse(issue)
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From POST")
		log.Error("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

func (c *Cli) CmdWatch(issue string, watcher string) error {
	log.Debug("watch called")

	json, err := jsonEncode(watcher)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/watchers", c.endpoint, issue)
	resp, err := c.post(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		c.Browse(issue)
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
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
		if strings.Contains(strings.ToLower(name), trans) {
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
		resp, err := c.post(uri, json)
		if err != nil {
			return err
		}
		if resp.StatusCode == 204 {
			c.Browse(issue)
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
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
		resp, err := c.post(uri, json)
		if err != nil {
			return err
		}

		if resp.StatusCode == 201 {
			c.Browse(issue)
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
			return nil
		} else {
			logBuffer := bytes.NewBuffer(make([]byte, 0))
			resp.Write(logBuffer)
			err := fmt.Errorf("Unexpected Response From POST")
			log.Error("%s:\n%s", err, logBuffer)
			return err
		}
	}

	if comment, ok := c.opts["comment"]; ok && comment != "" {
		json, err := jsonEncode(map[string]interface{}{
			"body": comment,
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

func (c *Cli) CmdAssign(issue string, user string) error {
	log.Debug("assign called")

	json, err := jsonEncode(map[string]interface{}{
		"name": user,
	})
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/assignee", c.endpoint, issue)
	resp, err := c.put(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		c.Browse(issue)
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
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
	dir := c.opts["directory"]
	if err := mkdir(dir); err != nil {
		return err
	}

	for name, template := range all_templates {
		if wanted, ok := c.opts["template"]; ok && wanted != name {
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
