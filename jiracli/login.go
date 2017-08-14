package jiracli

import (
	"fmt"
	"net/http"

	"github.com/mgutz/ansi"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func (jc *JiraCli) CmdLoginRegistry() *CommandRegistryEntry {
	opts := GlobalOptions{}
	return &CommandRegistryEntry{
		"Attempt to login into jira server",
		func() error {
			return jc.CmdLogin(&opts)
		},
		func(cmd *kingpin.CmdClause) error {
			return jc.GlobalUsage(cmd, &opts)
		},
	}
}

func authCallback(req *http.Request, resp *http.Response) (*http.Response, error) {
	if resp.StatusCode == 403 {
		defer resp.Body.Close()
		// X-Authentication-Denied-Reason: CAPTCHA_CHALLENGE; login-url=https://jira/login.jsp
		if reason := resp.Header.Get("X-Authentication-Denied-Reason"); reason != "" {
			return resp, fmt.Errorf("Authenticaion Failed: " + reason)
		}
		return resp, fmt.Errorf("Authenticaion Failed: Unkown Reason")
	} else if resp.StatusCode == 200 {
		if reason := resp.Header.Get("X-Seraph-Loginreason"); reason == "AUTHENTICATION_DENIED" {
			defer resp.Body.Close()
			return resp, fmt.Errorf("Authentication Failed: " + reason)
		}
	}

	return resp, nil
}

// CmdLogin will attempt to login into jira server
func (jc *JiraCli) CmdLogin(opts *GlobalOptions) error {
	defer func(h jira.HttpClient) {
		log.Debugf("Client: %#v", h)
		jc.UA = h
	}(jc.UA)
	if session, err := jc.GetSession(); err != nil {
		jc.UA = jc.oreoAgent.WithoutRedirect().WithRetries(0).WithPostCallback(authCallback)
		// No active session so try to create a new one
		_, err := jc.NewSession(opts)
		if err != nil {
			// reset password on failed session
			opts.SetPass("")
			return err
		}
		fmt.Println(ansi.Color("OK", "green"), "New session for", opts.User)
	} else {
		fmt.Println(ansi.Color("OK", "green"), "Found session for", session.Name)
	}
	return nil
}

// 	uri := fmt.Sprintf("%s/rest/auth/1/session", c.endpoint)
// 	for {
// 		req, _ := http.NewRequest("GET", uri, nil)
// 		user, _ := c.opts["user"].(string)

// 		passwd := c.GetPass(user)
// 		req.SetBasicAuth(user, passwd)

// 		resp, err := c.makeRequest(req)
// 		if err != nil {
// 			return err
// 		}
// 		if resp.StatusCode == 403 {
// 			// probably got this, need to redirect the user to login manually
// 			// X-Authentication-Denied-Reason: CAPTCHA_CHALLENGE; login-url=https://jira/login.jsp
// 			if reason := resp.Header.Get("X-Authentication-Denied-Reason"); reason != "" {
// 				err := fmt.Errorf("Authenticaion Failed: %s", reason)
// 				log.Errorf("%s", err)
// 				return err
// 			}
// 			err := fmt.Errorf("Authentication Failed: Unknown Reason")
// 			log.Errorf("%s", err)
// 			return err

// 		} else if resp.StatusCode == 200 {
// 			// https://confluence.atlassian.com/display/JIRA043/JIRA+REST+API+%28Alpha%29+Tutorial#JIRARESTAPI%28Alpha%29Tutorial-CAPTCHAs
// 			// probably bad password, try again
// 			if reason := resp.Header.Get("X-Seraph-Loginreason"); reason == "AUTHENTICATION_DENIED" {
// 				log.Warning("Authentication Failed: %s", reason)
// 				continue
// 			}
// 			if _, ok := c.opts["password-source"]; ok {
// 				return c.SetPass(user, passwd)
// 			}
// 			break
// 		} else {
// 			log.Warning("Login failed")
// 			continue
// 		}
// 	}
// 	return nil
// }
