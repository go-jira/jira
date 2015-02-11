package cli

import (
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"
	"os"
	"code.google.com/p/gopass"
)

func (c *Cli) CmdLogin() {
	uri := fmt.Sprintf("%s/rest/auth/1/session", c.endpoint)
	resp := c.get(uri)
    for ; resp.StatusCode != 200 ; {
		req, _ := http.NewRequest("GET", uri, nil)
		user, _ := c.opts["user"]
		
		prompt := fmt.Sprintf("Enter Password for %s: ", user)
		passwd, _ := gopass.GetPass(prompt);

		req.SetBasicAuth(user, passwd)
		resp = c.makeRequest(req)
		if resp.StatusCode == 403 {
			// probably got this, need to redirect the user to login manually
			// X-Authentication-Denied-Reason: CAPTCHA_CHALLENGE; login-url=https://jira/login.jsp
			if reason := resp.Header.Get("X-Authentication-Denied-Reason"); reason != "" {
				log.Error("Authentication Failed: %s", reason)
				os.Exit(1)
			}
			log.Error("Authentication Failead: Unknown")
			os.Exit(1)
			
		}
		if resp.StatusCode != 200 {
			log.Error("Login failed")
		}
	}
}

func (c *Cli) CmdFields() {
	log.Debug("fields called")
	resp := c.get(fmt.Sprintf("%s/rest/api/2/field", c.endpoint))
	data := jsonDecode(resp.Body)

	if templateFile, err := FindClosestParentPath(".jira.d/templates/fields"); err != nil {
		runTemplate(default_fields_template, data)
	} else {
		log.Debug("Using Template: %s", templateFile)
		runTemplate(readFile(templateFile), data)
	}
}

func (c *Cli) CmdList() {
	log.Debug("list called")

	if query, ok := c.opts["query"]; !ok {
		log.Error("No query argument found, either use --query or set query attribute in .jira file")
		os.Exit(1)
	} else {
		buffer := bytes.NewBuffer(make([]byte, 0, len(query)))
		enc := json.NewEncoder(buffer)
		
		enc.Encode(map[string]string{
			"jql": query,
			"startAt": "0",
			"maxResults": "500",
		})

		resp := c.post(fmt.Sprintf("%s/rest/api/2/search", c.endpoint), buffer)
		data := jsonDecode(resp.Body)
		
		if templateFile, err := FindClosestParentPath(".jira.d/templates/list"); err != nil {
			runTemplate(default_list_template, data)
		} else {
			log.Debug("Using Template: %s", templateFile)
			runTemplate(readFile(templateFile), data)
		}
	}

}

func (c *Cli) CmdView(issue string) {
	log.Debug("view called")
	resp := c.get(fmt.Sprintf("%s/rest/api/2/issue/%s", c.endpoint, issue))
	data := jsonDecode(resp.Body)
	if templateFile, err := FindClosestParentPath(".jira.d/templates/view"); err != nil {
		runTemplate(default_view_template, data)
	} else {
		log.Debug("Using Template: %s", templateFile)
		runTemplate(readFile(templateFile), data)
	}
}

