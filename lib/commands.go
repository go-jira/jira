package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/howeyc/gopass"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/data"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	// "github.com/kr/pretty"
)

// CmdLogin will attempt to login into jira server
func (c *Cli) CmdLogin() error {
	uri := fmt.Sprintf("%s/rest/auth/1/session", c.endpoint)
	for {
		req, _ := http.NewRequest("GET", uri, nil)
		user, _ := c.opts["user"].(string)

		fmt.Printf("Jira Password [%s]: ", user)
		pw, err := gopass.GetPasswdMasked()
		if err != nil {
			return err
		}
		passwd := string(pw)

		req.SetBasicAuth(user, passwd)

		resp, err := c.makeRequest(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 403 {
			// probably got this, need to redirect the user to login manually
			// X-Authentication-Denied-Reason: CAPTCHA_CHALLENGE; login-url=https://jira/login.jsp
			if reason := resp.Header.Get("X-Authentication-Denied-Reason"); reason != "" {
				err := fmt.Errorf("Authenticaion Failed: %s", reason)
				log.Errorf("%s", err)
				return err
			}
			err := fmt.Errorf("Authentication Failed: Unknown Reason")
			log.Errorf("%s", err)
			return err

		} else if resp.StatusCode == 200 {
			// https://confluence.atlassian.com/display/JIRA043/JIRA+REST+API+%28Alpha%29+Tutorial#JIRARESTAPI%28Alpha%29Tutorial-CAPTCHAs
			// probably bad password, try again
			if reason := resp.Header.Get("X-Seraph-Loginreason"); reason == "AUTHENTICATION_DENIED" {
				log.Warning("Authentication Failed: %s", reason)
				continue
			}
			break
		} else {
			log.Warning("Login failed")
			continue
		}
	}
	return nil
}

// CmdLogout will close any active sessions
func (c *Cli) CmdLogout() error {
	uri := fmt.Sprintf("%s/rest/auth/1/session", c.endpoint)
	req, _ := http.NewRequest("DELETE", uri, nil)
	resp, err := c.makeRequest(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == 401 || resp.StatusCode == 204 {
		// 401 == no active session
		// 204 == successfully logged out
	} else {
		err := fmt.Errorf("Failed to Logout: %s", err)
		return err
	}
	log.Notice("OK")
	return nil
}

// CmdFields will send data from /rest/api/2/field API to "fields" template
func (c *Cli) CmdFields() error {
	log.Debugf("fields called")
	uri := fmt.Sprintf("%s/rest/api/2/field", c.endpoint)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("fields"), data, nil)
}

// CmdList will query jira and send data to "list" template
func (c *Cli) CmdList() error {
	log.Debugf("list called")
	data, err := c.FindIssues(&SearchOptions{})
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("list"), data, nil)
}

func (c *Cli) List(sp SearchProvider) error {
	data, err := c.FindIssues(sp)
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("list"), data, nil)
}

// CmdView will get issue data and send to "view" template
func (c *Cli) CmdView(issue string) error {
	log.Debugf("view called")
	c.Browse(issue)
	data, err := c.ViewIssue(issue)
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("view"), data, nil)
}

// CmdWorklogs will get worklog data for given issue and sent to the "worklogs" template
func (c *Cli) CmdWorklogs(issue string) error {
	log.Debugf("worklogs called")
	c.Browse(issue)
	data, err := c.ViewIssueWorkLogs(issue)
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("worklogs"), data, nil)
}

// CmdWorklog will attempt to add (action=add) a worklog to the given issue.
// It will spawn the editor (unless --noedit isused) and post edited YAML
// content as JSON to the worklog endpoint
func (c *Cli) CmdWorklog(action string, issue string) error {
	log.Debugf("%s worklog called", action)
	c.Browse(issue)
	if action == "add" {
		uri := fmt.Sprintf("%s/rest/api/2/issue/%s/worklog", c.endpoint, issue)

		worklogData := map[string]interface{}{
			"issue":   issue,
			"comment": c.opts["comment"],
		}

		if v, ok := c.opts["time-spent"].(string); ok {
			worklogData["timeSpent"] = v
		}

		return c.editTemplate(
			c.getTemplate("worklog"),
			fmt.Sprintf("%s-worklog-", issue),
			worklogData,
			func(json string) error {
				if c.getOptBool("dryrun", false) {
					log.Debugf("POST: %s", json)
					log.Debugf("Dryrun mode, skipping POST")
					return nil
				}
				resp, err := c.post(uri, json)
				if err != nil {
					return err
				}

				if resp.StatusCode == 201 {
					c.Browse(issue)
					if !c.opts["quiet"].(bool) {
						fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
					}
					return nil
				}
				logBuffer := bytes.NewBuffer(make([]byte, 0))
				resp.Write(logBuffer)
				err = fmt.Errorf("Unexpected Response From POST")
				log.Errorf("%s:\n%s", err, logBuffer)
				return err
			},
		)
	}
	return nil
}

// CmdEdit will populate "edit" template with issue data and issue "editmeta" data.
// Then will parse yaml template and submit data to jira.
func (c *Cli) CmdEdit(issue string) error {
	log.Debugf("edit called")

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/editmeta", c.endpoint, issue)
	editmeta, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}

	uri = fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}
	issueData := data.(map[string]interface{})
	issueData["meta"] = editmeta.(map[string]interface{})
	issueData["overrides"] = c.opts

	return c.editTemplate(
		c.getTemplate("edit"),
		fmt.Sprintf("%s-edit-", issue),
		issueData,
		func(json string) error {
			if c.getOptBool("dryrun", false) {
				log.Debugf("PUT: %s", json)
				log.Debugf("Dryrun mode, skipping PUT")
				return nil
			}
			resp, err := c.put(uri, json)
			if err != nil {
				return err
			}

			if resp.StatusCode == 204 {
				c.Browse(issueData["key"].(string))
				if !c.opts["quiet"].(bool) {
					fmt.Printf("OK %s %s/browse/%s\n", issueData["key"], c.endpoint, issueData["key"])
				}
				return nil
			}
			logBuffer := bytes.NewBuffer(make([]byte, 0))
			resp.Write(logBuffer)
			err = fmt.Errorf("Unexpected Response From PUT")
			log.Errorf("%s:\n%s", err, logBuffer)
			return err
		},
	)
}

// CmdEditMeta will send issue 'edit' metadata to the "editmeta" template
func (c *Cli) CmdEditMeta(issue string) error {
	log.Debugf("editMeta called")
	c.Browse(issue)
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/editmeta", c.endpoint, issue)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("editmeta"), data, nil)
}

// CmdTransitionMeta will send available transition metadata to the "transmeta" template
func (c *Cli) CmdTransitionMeta(issue string) error {
	log.Debugf("tranisionMeta called")
	c.Browse(issue)
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions?expand=transitions.fields", c.endpoint, issue)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("transmeta"), data, nil)
}

// CmdIssueTypes will send issue 'create' metadata to the 'issuetypes'
func (c *Cli) CmdIssueTypes() error {
	project := c.opts["project"].(string)
	log.Debugf("issueTypes called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s", c.endpoint, project)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}

	return runTemplate(c.getTemplate("issuetypes"), data, nil)
}

func (c *Cli) defaultIssueType() string {
	project := c.opts["project"].(string)
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s", c.endpoint, project)
	data, _ := responseToJSON(c.get(uri))
	issueTypeNames := make(map[string]bool)

	if data, ok := data.(map[string]interface{}); ok {
		if projects, ok := data["projects"].([]interface{}); ok {
			for _, project := range projects {
				if project, ok := project.(map[string]interface{}); ok {
					if issuetypes, ok := project["issuetypes"].([]interface{}); ok {
						if len(issuetypes) > 0 {
							for _, issuetype := range issuetypes {
								issueTypeNames[issuetype.(map[string]interface{})["name"].(string)] = true
							}
						}
					}
				}
			}
		}
	}
	if _, ok := issueTypeNames["Bug"]; ok {
		return "Bug"
	} else if _, ok := issueTypeNames["Task"]; ok {
		return "Task"
	}
	return ""
}

// CmdCreateMeta sends the 'create' metadata for the given project & issuetype and sends it to the 'createmeta' template
func (c *Cli) CmdCreateMeta() error {
	project := c.opts["project"].(string)
	issuetype := c.getOptString("issuetype", "")
	if issuetype == "" {
		issuetype = c.defaultIssueType()
	}

	log.Debugf("createMeta called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&issuetypeNames=%s&expand=projects.issuetypes.fields", c.endpoint, project, url.QueryEscape(issuetype))
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}

	if val, ok := data.(map[string]interface{})["projects"]; ok {
		if len(val.([]interface{})) == 0 {
			err = fmt.Errorf("Project '%s' or issuetype '%s' unknown.  Unable to createmeta.", project, issuetype)
			log.Errorf("%s", err)
			return err
		}
		if val, ok = val.([]interface{})[0].(map[string]interface{})["issuetypes"]; ok {
			data = val.([]interface{})[0]
		}
	}

	return runTemplate(c.getTemplate("createmeta"), data, nil)
}

// CmdComponents sends component data for given project and sends to the "components" template
func (c *Cli) CmdComponents(project string) error {
	log.Debugf("Components called")
	uri := fmt.Sprintf("%s/rest/api/2/project/%s/components", c.endpoint, project)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("components"), data, nil)
}

// ValidTransitions will return a list of valid transitions for given issue.
func (c *Cli) ValidTransitions(issue string) (jiradata.Transitions, error) {
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions?expand=transitions.fields", c.endpoint, issue)
	resp, err := c.get(uri)
	if err != nil {
		return nil, err
	}

	transMeta := &jiradata.TransitionsMeta{}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(content, transMeta)
	if err != nil {
		return nil, err
	}

	return transMeta.Transitions, nil
}

// CmdTransitions sends valid transtions for given issue to the "transitions" template
func (c *Cli) CmdTransitions(issue string) error {
	log.Debugf("Transitions called")
	// FIXME this should just call ValidTransitions then pass that data to templates
	c.Browse(issue)
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.endpoint, issue)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("transitions"), data, nil)
}

// CmdCreate sends the create-metadata to the "create" template for editing, then
// will parse the edited document as YAML and submit the document to jira.
func (c *Cli) CmdCreate() error {
	log.Debugf("create called")
	project := c.opts["project"].(string)
	issuetype := c.getOptString("issuetype", "")
	if issuetype == "" {
		issuetype = c.defaultIssueType()
	}

	uri := fmt.Sprintf("%s/rest/api/2/issue/createmeta?projectKeys=%s&issuetypeNames=%s&expand=projects.issuetypes.fields", c.endpoint, project, url.QueryEscape(issuetype))
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}

	issueData := make(map[string]interface{})
	issueData["overrides"] = c.opts
	issueData["overrides"].(map[string]interface{})["issuetype"] = issuetype

	if val, ok := data.(map[string]interface{})["projects"]; ok {
		if len(val.([]interface{})) == 0 {
			err = fmt.Errorf("Project '%s' or issuetype '%s' unknown.  Unable to create issue.", project, issuetype)
			log.Errorf("%s", err)
			return err
		}
		if val, ok = val.([]interface{})[0].(map[string]interface{})["issuetypes"]; ok {
			if len(val.([]interface{})) == 0 {
				err = fmt.Errorf("Project '%s' does not support issuetype '%s'.  Unable to create issue.", project, issuetype)
				log.Errorf("%s", err)
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
			uri := fmt.Sprintf("%s/rest/api/2/issue", c.endpoint)
			if c.getOptBool("dryrun", false) {
				log.Debugf("POST: %s", json)
				log.Debugf("Dryrun mode, skipping POST")
				return nil
			}
			resp, err := c.post(uri, json)
			if err != nil {
				return err
			}

			if resp.StatusCode == 201 {
				// response: {"id":"410836","key":"PROJ-238","self":"https://jira/rest/api/2/issue/410836"}
				json, err := responseToJSON(resp, nil)
				if err != nil {
					return err
				}
				key := json.(map[string]interface{})["key"].(string)
				link := fmt.Sprintf("%s/browse/%s", c.endpoint, key)
				c.Browse(key)
				c.SaveData(map[string]string{
					"issue": key,
					"link":  link,
				})
				if !c.opts["quiet"].(bool) {
					fmt.Printf("OK %s %s\n", key, link)
				}
				return nil
			}
			logBuffer := bytes.NewBuffer(make([]byte, 0))
			resp.Write(logBuffer)
			err = fmt.Errorf("Unexpected Response From POST")
			log.Errorf("%s:\n%s", err, logBuffer)
			return err
		},
	)
}

// CmdIssueLinkTypes will send the issue link type data to the "issuelinktypes" template.
func (c *Cli) CmdIssueLinkTypes() error {
	log.Debugf("Transitions called")
	uri := fmt.Sprintf("%s/rest/api/2/issueLinkType", c.endpoint)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("issuelinktypes"), data, nil)
}

// CmdBlocks will update the given issue as being "blocked" by the given blocker
func (c *Cli) CmdBlocks(blocker string, issue string) error {
	log.Debugf("blocks called")

	json, err := jsonEncode(map[string]interface{}{
		"type": map[string]string{
			"name": c.GetOptString("blockerType", "Blocks"),
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
	if c.getOptBool("dryrun", false) {
		log.Debugf("POST: %s", json)
		log.Debugf("Dryrun mode, skipping POST")
		return nil
	}
	resp, err := c.post(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		c.Browse(issue)
		if !c.opts["quiet"].(bool) {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From POST")
		log.Errorf("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

// CmdDups will update the given issue as being a duplicate by the given dup issue
// and will attempt to resolve the dup issue
func (c *Cli) CmdDups(duplicate string, issue string) error {
	log.Debugf("dups called")

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
	if c.getOptBool("dryrun", false) {
		log.Debugf("POST: %s", json)
		log.Debugf("Dryrun mode, skipping POST")
		return nil
	}
	resp, err := c.post(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		c.Browse(issue)
		if !c.opts["quiet"].(bool) {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From POST")
		log.Errorf("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

// CmdWatch will add the given watcher to the issue (or remove the watcher
// given the 'remove' flag)
func (c *Cli) CmdWatch(issue string, watcher string, remove bool) error {
	log.Debugf("watch called: watcher: %q, remove: %n", watcher, remove)

	var uri string
	json, err := jsonEncode(watcher)
	if err != nil {
		return err
	}

	if c.getOptBool("dryrun", false) {
		if !remove {
			log.Debugf("POST: %s", json)
			log.Debugf("Dryrun mode, skipping POST")
		} else {
			log.Debugf("DELETE: %s", watcher)
			log.Debugf("Dryrun mode, skipping POST")
		}
		return nil
	}

	var resp *http.Response
	if !remove {
		uri = fmt.Sprintf("%s/rest/api/2/issue/%s/watchers", c.endpoint, issue)
		resp, err = c.post(uri, json)
	} else {
		uri = fmt.Sprintf("%s/rest/api/2/issue/%s/watchers?username=%s", c.endpoint, issue, watcher)
		resp, err = c.delete(uri)
	}
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		c.Browse(issue)
		if !c.opts["quiet"].(bool) {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		if !remove {
			err = fmt.Errorf("Unexpected Response From POST")
		} else {
			err = fmt.Errorf("Unexpected Response From DELETE")
		}
		log.Errorf("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

// CmdVote will add or remove a vote on an issue
func (c *Cli) CmdVote(issue string, up bool) error {
	log.Debugf("vote called, with up: %n", up)

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/votes", c.endpoint, issue)
	if c.getOptBool("dryrun", false) {
		if up {
			log.Debugf("POST: %s", "")
			log.Debugf("Dryrun mode, skipping POST")
		} else {
			log.Debugf("DELETE: %s", "")
			log.Debugf("Dryrun mode, skipping DELETE")
		}
		return nil
	}
	var resp *http.Response
	var err error
	if up {
		resp, err = c.post(uri, "")
	} else {
		resp, err = c.delete(uri)
	}
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		c.Browse(issue)
		if !c.opts["quiet"].(bool) {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		if up {
			err = fmt.Errorf("Unexpected Response From POST")
		} else {
			err = fmt.Errorf("Unexpected Response From DELETE")
		}
		log.Errorf("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

// CmdRankAfter rank issue after target issue
func (c *Cli) CmdRankAfter(issue, after string) error {
	err := c.RankIssue(issue, after, RANKAFTER)
	if err != nil {
		return nil
	}
	if !c.opts["quiet"].(bool) {
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
	}
	return nil
}

// CmdRankBefore rank issue before target issue
func (c *Cli) CmdRankBefore(issue, before string) error {
	err := c.RankIssue(issue, before, RANKBEFORE)
	if err != nil {
		return nil
	}
	if !c.opts["quiet"].(bool) {
		fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
	}
	return nil
}

// CmdTransition will move state of the given issue to the given transtion
func (c *Cli) CmdTransition(issue string, trans string) error {
	log.Debugf("transition called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/transitions?expand=transitions.fields", c.endpoint, issue)
	data, err := responseToJSON(c.get(uri))
	if err != nil {
		return err
	}

	transitions := data.(map[string]interface{})["transitions"].([]interface{})
	var transID, transName string
	var transMeta map[string]interface{}
	found := make([]string, 0, len(transitions))
	for _, transition := range transitions {
		name := transition.(map[string]interface{})["name"].(string)
		id := transition.(map[string]interface{})["id"].(string)
		found = append(found, name)
		if strings.Contains(strings.ToLower(name), strings.ToLower(trans)) {
			transName = name
			transID = id
			transMeta = transition.(map[string]interface{})
			if strings.ToLower(name) == strings.ToLower(trans) {
				break
			}
		}
	}
	if transID == "" {
		err := fmt.Errorf("Invalid Transition '%s', Available: %s", trans, strings.Join(found, ", "))
		log.Debugf("%s", err)
		return err
	}

	handlePost := func(json string) error {
		uri = fmt.Sprintf("%s/rest/api/2/issue/%s/transitions", c.endpoint, issue)
		if c.getOptBool("dryrun", false) {
			log.Debugf("POST: %s", json)
			log.Debugf("Dryrun mode, skipping POST")
			return nil
		}
		resp, err := c.post(uri, json)
		if err != nil {
			return err
		}
		if resp.StatusCode == 204 {
			c.Browse(issue)
			if !c.opts["quiet"].(bool) {
				fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
			}
		} else {
			logBuffer := bytes.NewBuffer(make([]byte, 0))
			resp.Write(logBuffer)
			err := fmt.Errorf("Unexpected Response From POST")
			log.Errorf("%s:\n%s", err, logBuffer)
			return err
		}
		return nil
	}

	uri = fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
	data, err = responseToJSON(c.get(uri))
	if err != nil {
		return err
	}
	issueData := data.(map[string]interface{})
	issueData["meta"] = transMeta
	if c.GetOptString("defaultResolution", "") == "" {
		// .meta.fields.resolution.allowedValues
		if fields, ok := transMeta["fields"].(map[string]interface{}); ok {
			if resolution, ok := fields["resolution"].(map[string]interface{}); ok {
				if allowedValues, ok := resolution["allowedValues"].([]interface{}); ok {
					for _, allowedValueRaw := range allowedValues {
						if allowedValues, ok := allowedValueRaw.(map[string]interface{}); ok {
							if allowedValues["name"] == "Fixed" {
								c.opts["defaultResolution"] = "Fixed"
							} else if allowedValues["name"] == "Done" {
								c.opts["defaultResolution"] = "Done"
							}
						}
					}
				}
			}
		}
	}
	issueData["overrides"] = c.opts
	issueData["transition"] = map[string]interface{}{
		"name": transName,
		"id":   transID,
	}
	return c.editTemplate(
		c.getTemplate("transition"),
		fmt.Sprintf("%s-trans-%s-", issue, trans),
		issueData,
		handlePost,
	)
}

// CmdComment will open up editor with "comment" template and submit
// YAML output to jira
func (c *Cli) CmdComment(issue string) error {
	log.Debugf("comment called")

	handlePost := func(json string) error {
		uri := fmt.Sprintf("%s/rest/api/2/issue/%s/comment", c.endpoint, issue)
		if c.getOptBool("dryrun", false) {
			log.Debugf("POST: %s", json)
			log.Debugf("Dryrun mode, skipping POST")
			return nil
		}
		resp, err := c.post(uri, json)
		if err != nil {
			return err
		}

		if resp.StatusCode == 201 {
			c.Browse(issue)
			if !c.opts["quiet"].(bool) {
				fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
			}
			return nil
		}
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err = fmt.Errorf("Unexpected Response From POST")
		log.Errorf("%s:\n%s", err, logBuffer)
		return err
	}

	if comment, ok := c.opts["comment"]; ok && comment != "" {
		json, err := jsonEncode(map[string]interface{}{
			"body": comment,
		})
		if err != nil {
			return err
		}
		return handlePost(json)
	}
	return c.editTemplate(
		c.getTemplate("comment"),
		fmt.Sprintf("%s-create-", issue),
		map[string]interface{}{},
		handlePost,
	)
}

// CmdComponent will add a new component to given project
func (c *Cli) CmdComponent(action string, project string, name string, desc string, lead string) error {
	log.Debugf("component called")

	switch action {
	case "add":
	default:
		return fmt.Errorf("CmdComponent: %q is not a valid action", action)
	}

	json, err := jsonEncode(map[string]interface{}{
		"name":         name,
		"description":  desc,
		"leadUserName": lead,
		"project":      project,
	})
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/component", c.endpoint)
	if c.getOptBool("dryrun", false) {
		log.Debugf("POST: %s", json)
		log.Debugf("Dryrun mode, skipping POST")
		return nil
	}
	resp, err := c.post(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 201 {
		if !c.opts["quiet"].(bool) {
			fmt.Printf("OK %s %s\n", project, name)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From POST")
		log.Errorf("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

// CmdLabels will add, remove or set labels on a given issue
func (c *Cli) CmdLabels(action string, issue string, labels []string) error {
	log.Debugf("label called")

	if action != "add" && action != "remove" && action != "set" {
		return fmt.Errorf("action must be 'add', 'set' or 'remove': %q is invalid", action)
	}

	handlePut := func(json string) error {
		uri := fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue)
		if c.getOptBool("dryrun", false) {
			log.Debugf("PUT: %s", json)
			log.Debugf("Dryrun mode, skipping POST")
			return nil
		}
		resp, err := c.put(uri, json)
		if err != nil {
			return err
		}

		if resp.StatusCode == 204 {
			c.Browse(issue)
			if !c.opts["quiet"].(bool) {
				fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
			}
			return nil
		}
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err = fmt.Errorf("Unexpected Response From PUT")
		log.Errorf("%s:\n%s", err, logBuffer)
		return err
	}

	var labelsJSON string
	var err error
	if action == "set" {
		labelsActions := make([]map[string][]string, 1)
		labelsActions[0] = map[string][]string{
			"set": labels,
		}
		labelsJSON, err = jsonEncode(map[string]interface{}{
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
		labelsJSON, err = jsonEncode(map[string]interface{}{
			"labels": labelsActions,
		})
	}
	if err != nil {
		return err
	}
	json := fmt.Sprintf("{ \"update\": %s }", labelsJSON)
	return handlePut(json)

}

// CmdAssign will assign the given user to be the owner of the given issue
func (c *Cli) CmdAssign(issue string, user string) error {
	log.Debugf("assign called")

	json, err := jsonEncode(map[string]interface{}{
		"name": user,
	})
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/assignee", c.endpoint, issue)
	if c.getOptBool("dryrun", false) {
		log.Debugf("PUT: %s", json)
		log.Debugf("Dryrun mode, skipping PUT")
		return nil
	}
	resp, err := c.put(uri, json)
	if err != nil {
		return err
	}
	if resp.StatusCode == 204 {
		c.Browse(issue)
		if !c.opts["quiet"].(bool) {
			fmt.Printf("OK %s %s/browse/%s\n", issue, c.endpoint, issue)
		}
	} else {
		logBuffer := bytes.NewBuffer(make([]byte, 0))
		resp.Write(logBuffer)
		err := fmt.Errorf("Unexpected Response From PUT")
		log.Errorf("%s:\n%s", err, logBuffer)
		return err
	}
	return nil
}

// CmdExportTemplates will export the default templates to the template directory.
func (c *Cli) CmdExportTemplates() error {
	dir := c.opts["directory"].(string)
	if err := mkdir(dir); err != nil {
		return err
	}

	for name, template := range allTemplates {
		if wanted, ok := c.opts["template"]; ok && wanted != name {
			continue
		}
		templateFile := fmt.Sprintf("%s/%s", dir, name)
		if _, err := os.Stat(templateFile); err == nil {
			log.Warning("Skipping %s, already exists", templateFile)
			continue
		}
		fh, err := os.OpenFile(templateFile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Errorf("Failed to open %s for writing: %s", templateFile, err)
			return err
		}
		defer fh.Close()
		log.Noticef("Creating %s", templateFile)
		fh.Write([]byte(template))
	}
	return nil
}

// CmdRequest will use the given uri to make a request and potentially send provided content.
func (c *Cli) CmdRequest(uri, content string) (err error) {
	log.Debugf("request called")

	if !strings.HasPrefix(uri, "http") {
		uri = fmt.Sprintf("%s%s", c.endpoint, uri)
	}

	method := strings.ToUpper(c.opts["method"].(string))
	var data interface{}
	if method == "GET" {
		data, err = responseToJSON(c.get(uri))
	} else if method == "POST" {
		data, err = responseToJSON(c.post(uri, content))
	} else if method == "PUT" {
		data, err = responseToJSON(c.put(uri, content))
	}
	if err != nil {
		return err
	}
	return runTemplate(c.getTemplate("request"), data, nil)
}
