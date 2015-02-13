package cli

import (
	"net/http"
	"fmt"
	"code.google.com/p/gopass"
	"bytes"
	"strings"
	// "github.com/kr/pretty"
)

func (c *Cli) CmdLogin() (error) {
	uri := fmt.Sprintf("%s/rest/auth/1/session", c.endpoint)
    for ; true ; {
		req, _ := http.NewRequest("GET", uri, nil)
		user, _ := c.opts["user"]
		
		prompt := fmt.Sprintf("Enter Password for %s: ", user)
		passwd, _ := gopass.GetPass(prompt);

		req.SetBasicAuth(user, passwd)
		if resp, err := c.makeRequest(req); err != nil {
			return err
		} else {
			if resp.StatusCode == 403 {
				// probably got this, need to redirect the user to login manually
				// X-Authentication-Denied-Reason: CAPTCHA_CHALLENGE; login-url=https://jira/login.jsp
				if reason := resp.Header.Get("X-Authentication-Denied-Reason"); reason != "" {
					log.Error("Authentication Failed: %s", reason)
					return fmt.Errorf("Authenticaion Failed: %s", reason)
				}
				log.Error("Authentication Failead: Unknown")
				return fmt.Errorf("Authentication Failead")
				
			}
			if resp.StatusCode != 200 {
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
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}

	return runTemplate(c.getTemplate(".jira.d/templates/fields", default_fields_template), data, nil)
}


func (c *Cli) CmdList() error {
	log.Debug("list called")

	if query, ok := c.opts["query"]; !ok {
		log.Error("No query argument found, either use --query or set query attribute in .jira file")
		return fmt.Errorf("Missing query")
	} else {
		json, err := jsonEncode(map[string]string{
			"jql": query,
			"startAt": "0",
			"maxResults": "500",
		}); if err != nil {
			return err
		}

		uri := fmt.Sprintf("%s/rest/api/2/search", c.endpoint)
		data, err := responseToJson(c.post(uri, json)); if err != nil {
			return err
		}

		return runTemplate(c.getTemplate(".jira.d/templates/list", default_list_template), data, nil)
	}
}

func (c *Cli) CmdView(issue string) error {
	log.Debug("view called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}

	return runTemplate(c.getTemplate(".jira.d/templates/view", default_view_template), data, nil)
}

func (c *Cli) CmdEdit(issue string) error {
	log.Debug("edit called")

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/editmeta", c.endpoint, issue)
	editmeta, err := responseToJson(c.get(uri)); if err != nil {
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
		c.getTemplate(".jira.d/templates/edit", default_edit_template),
		fmt.Sprintf("%s-edit-", issue),
		issueData,
		func(json string) error {
			resp, err := c.put(uri, json); if err != nil {
				return err
			}
			
			if resp.StatusCode == 204 {
				fmt.Printf("OK %s %s/browse/%s\n", issueData["key"], c.endpoint, issueData["key"])
				return nil
			} else {
				logBuffer := bytes.NewBuffer(make([]byte,0))
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
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/editmeta", c.endpoint, issue)
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}
	
	return runTemplate(c.getTemplate(".jira.d/templates/editmeta", default_fields_template), data, nil)
}

func (c *Cli) CmdTransitionMeta(issue string) error {
	log.Debug("tranisionMeta called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions?expand=transitions.fields", c.endpoint, issue)
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}
	
	return runTemplate(c.getTemplate(".jira.d/templates/transmeta", default_fields_template), data, nil)
}

func (c *Cli) CmdIssueTypes(project string) error {
	log.Debug("issueTypes called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s", c.endpoint, project)
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}

	return runTemplate(c.getTemplate(".jira.d/templates/issuetypes", default_issuetypes_template), data, nil)
}

func (c *Cli) CmdCreateMeta(project string, issuetype string) error {
	log.Debug("createMeta called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&issuetypeNames=%s&expand=projects.issuetypes.fields", c.endpoint, project, issuetype)
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}

	if val, ok := data.(map[string]interface{})["projects"]; ok {
		if val, ok = val.([]interface{})[0].(map[string]interface{})["issuetypes"]; ok {
			data = val.([]interface{})[0]
		}
	}

	return runTemplate(c.getTemplate(".jira.d/templates/createmeta", default_fields_template), data, nil)
}

func (c *Cli) CmdTransitions(issue string) error {
	log.Debug("Transitions called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.endpoint, issue)
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}
	return runTemplate(c.getTemplate(".jira.d/templates/transitions", default_transitions_template), data, nil)
}

func (c *Cli) CmdCreate(project string, issuetype string) error {
	log.Debug("create called")

	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&issuetypeNames=%s&expand=projects.issuetypes.fields", c.endpoint, project, issuetype)
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}

	issueData := make(map[string]interface{})
	issueData["overrides"] = c.opts

	if val, ok := data.(map[string]interface{})["projects"]; ok {
		if val, ok = val.([]interface{})[0].(map[string]interface{})["issuetypes"]; ok {
			issueData["meta"] = val.([]interface{})[0]
		}
	}
	
	sanitizedType := strings.ToLower(strings.Replace(issuetype, " ", "", -1))
	return c.editTemplate(
		c.getTemplate(fmt.Sprintf(".jira.d/templates/create-%s", sanitizedType), default_create_template),
		fmt.Sprintf("create-%s-", sanitizedType),
		issueData,
		func(json string) error {
			log.Debug("JSON: %s", json)
			uri := fmt.Sprintf("%s/rest/api/2/issue", c.endpoint)
			resp, err := c.post(uri, json); if err != nil {
				return err
			}
			
			if resp.StatusCode == 201 {
				// response: {"id":"410836","key":"PROJ-238","self":"https://jira/rest/api/2/issue/410836"}
				if json, err := responseToJson(resp, nil); err != nil {
					return err
				} else {
					key := json.(map[string]interface{})["key"]
					fmt.Printf("OK %s %s/browse/%s\n", key, c.endpoint, key)
				}
				return nil
			} else {
				logBuffer := bytes.NewBuffer(make([]byte,0))
				resp.Write(logBuffer)
				err := fmt.Errorf("Unexpected Response From PUT")
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
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}
	return runTemplate(c.getTemplate(".jira.d/templates/issuelinktypes", default_fields_template), data, nil)
}

func (c *Cli) CmdBlocks(blocker string, issue string) error {
	log.Debug("blocks called")

	json, err := jsonEncode(map[string]interface{}{
		"type": map[string]string{
			"name": "Depends",  // TODO This is probably not constant across Jira installs
		},
		"inwardIssue": map[string]string{
			"key": issue,
		},
		"outwardIssue": map[string]string{
			"key": blocker,
		},
	}); if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issueLink", c.endpoint)
	resp, err := c.post(uri, json); if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
	} else {
		logBuffer := bytes.NewBuffer(make([]byte,0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From PUT")
		log.Error("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

func (c *Cli) CmdDups(duplicate string, issue string) error {
	log.Debug("dups called")

	json, err := jsonEncode(map[string]interface{}{
		"type": map[string]string{
			"name": "Duplicate",  // TODO This is probably not constant across Jira installs
		},
		"inwardIssue": map[string]string{
			"key": duplicate,
		},
		"outwardIssue": map[string]string{
			"key": issue,
		},
	}); if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issueLink", c.endpoint)
	resp, err := c.post(uri, json); if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
	} else {
		logBuffer := bytes.NewBuffer(make([]byte,0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From PUT")
		log.Error("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}


func (c *Cli) CmdWatch(issue string, watcher string) error {
	log.Debug("watch called")

	json, err := jsonEncode(watcher); if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/watchers", c.endpoint, issue)
	resp, err := c.post(uri, json); if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
	} else {
		logBuffer := bytes.NewBuffer(make([]byte,0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From PUT")
		log.Error("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

func (c *Cli) CmdTransition(issue string, trans string) error {
	log.Debug("transition called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.endpoint, issue)
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}

	transitions := data.(map[string]interface{})["transitions"].([]interface{})
	var	transId string
	found := make([]string, 0, len(transitions))
	for _, transition := range transitions {
		name := transition.(map[string]interface{})["name"].(string)
		id  :=  transition.(map[string]interface{})["id"].(string)
		found = append(found, name)
		if strings.Contains(strings.ToLower(name), trans) {
			transId = id
		}
	}
	if transId == "" {
		err := fmt.Errorf("Invalid Transition '%s', Available: %s", trans, strings.Join(found, ", "))
		log.Error("%s", err)
		return err
	}
	

	payload := map[string]interface{}{
		"transition": map[string]interface{}{
			"id": transId,
		},
	}

	if comment, ok := c.opts["comment"]; ok {
		payload["update"] = map[string]interface{}{
			"comment": []interface{}{
				map[string]interface{}{
					"add": map[string]interface{}{
						"body": comment,
					},
				},
			},
		}
	}

	json, err := jsonEncode(payload); if err != nil {
		return err
	}
	
	uri = fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.endpoint, issue)
	resp, err := c.post(uri, json); if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
	} else {
		logBuffer := bytes.NewBuffer(make([]byte,0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From POST")
		log.Error("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

func (c *Cli) CmdComment(issue string) error {
	log.Debug("comment called")
	
	handlePost := func(json string) error {
		log.Debug("JSON: %s", json)
		uri := fmt.Sprintf("%s/rest/api/2/issue/%s/comment", c.endpoint, issue)
		resp, err := c.post(uri, json); if err != nil {
			return err
		}
		
		if resp.StatusCode == 201 {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
			return nil
		} else {
			logBuffer := bytes.NewBuffer(make([]byte,0))
			resp.Write(logBuffer)
			err := fmt.Errorf("Unexpected Response From PUT")
			log.Error("%s:\n%s", err, logBuffer)
			return err
		}
	}

	if comment, ok := c.opts["comment"]; ok {
		json, err := jsonEncode(map[string]interface{}{
			"body": comment,
		}); if err != nil {
			return err
		}
		return handlePost(json)
	} else {
		return c.editTemplate(
			c.getTemplate(".jira.d/templates/comment", default_comment_template),
			fmt.Sprintf("%s-create-", issue),
			map[string]interface{}{},
			handlePost,
		)
	}
	return nil
}
