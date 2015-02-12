package cli

import (
	"net/http"
	"fmt"
	"code.google.com/p/gopass"
	"os"
	"bytes"
	"os/exec"
	"io/ioutil"
	"gopkg.in/yaml.v1"
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

	issueData["meta"] = editmeta.(map[string]interface{})["fields"]
	
	tmpdir := fmt.Sprintf("%s/.jira.d/tmp", os.Getenv("HOME"))
	fh, err := ioutil.TempFile(tmpdir, fmt.Sprintf("%s-edit-", issue)); if err != nil {
		log.Error("Failed to make temp file in %s: %s", tmpdir, err)
		return err
	}
	defer fh.Close()
	
	tmpFileName := fmt.Sprintf("%s.yml", fh.Name())
	if err := os.Rename(fh.Name(), tmpFileName); err != nil {
		log.Error("Failed to rename %s to %s: %s", fh.Name(), fmt.Sprintf("%s.yml", fh.Name()), err)
		return err
	}
	
	err = runTemplate(c.getTemplate(".jira.d/templates/edit", default_edit_template), issueData, fh); if err != nil {
		return err
	}
	
	fh.Close()
	
	editor, ok := c.opts["editor"]; if !ok {
		editor = os.Getenv("JIRA_EDITOR"); if editor == "" {
			editor = os.Getenv("EDITOR"); if editor == "" {
				editor = "vim"
			}
		}
	}
	for ; true ; {
		log.Debug("Running: %s %s", editor, tmpFileName)
		cmd := exec.Command(editor, tmpFileName)
		cmd.Stdout, cmd.Stderr, cmd.Stdin = os.Stdout, os.Stderr, os.Stdin
		if err := cmd.Run(); err != nil {
			log.Error("Failed to edit template with %s: %s", editor, err)
			if promptYN("edit again?", true) {
				continue
			}
			return err
		}
		
		edited := make(map[string]interface{})
		if fh, err := ioutil.ReadFile(tmpFileName); err != nil {
			log.Error("Failed to read tmpfile %s: %s", tmpFileName, err)
			if promptYN("edit again?", true) {
				continue
			}
			return err
		} else {
			if err := yaml.Unmarshal(fh, &edited); err != nil {
				log.Error("Failed to parse YAML: %s", err)
				if promptYN("edit again?", true) {
					continue
				}
				return err
			}
		}
		
		if fixed, err := yamlFixup(edited); err != nil {
			return err
		} else {
			edited = fixed.(map[string]interface{})
		}
		
		mf := editmeta.(map[string]interface{})["fields"]
		f  := edited["fields"].(map[string]interface{})
		for k, _ := range f {
			if _, ok := mf.(map[string]interface{})[k]; !ok {
				err := fmt.Errorf("Field %s is not editable", k)
				log.Error("%s", err)
				if promptYN("edit again?", true) {
					continue
				}
				return err
			}
		}

		json, err := jsonEncode(edited); if err != nil {
			return err
		}
		
		resp, err := c.put(uri, json); if err != nil {
			return err
		}

		if resp.StatusCode == 204 {
			fmt.Printf("OK %s %s", issueData["key"], issueData["self"])
			return nil
		} else {
			logBuffer := bytes.NewBuffer(make([]byte,0))
			resp.Write(logBuffer)
			err := fmt.Errorf("Unexpected Response From PUT")
			log.Error("%s:\n%s", err, logBuffer)
			return err
		}
	}
	return nil
}

func (c *Cli) CmdEditMeta(issue string) error {
	log.Debug("editMeta called")
	uri := fmt.Sprintf("%s/rest/api/2/issue/%s/editmeta", c.endpoint, issue)
	data, err := responseToJson(c.get(uri)); if err != nil {
		return err
	}
	
	return runTemplate(c.getTemplate(".jira.d/templates/editmeta", default_fields_template), data, nil)
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
