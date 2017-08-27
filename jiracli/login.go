package jiracli

import (
	"fmt"
	"net/http"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"
	"github.com/mgutz/ansi"
	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func CmdLoginRegistry(fig *figtree.FigTree, o *oreo.Client) *CommandRegistryEntry {
	opts := GlobalOptions{}
	return &CommandRegistryEntry{
		"Attempt to login into jira server",
		func() error {
			return CmdLogin(o, &opts)
		},
		func(cmd *kingpin.CmdClause) error {
			LoadConfigs(cmd, fig, &opts)
			return GlobalUsage(cmd, &opts)
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
func CmdLogin(o *oreo.Client, opts *GlobalOptions) error {
	if session, err := jira.GetSession(o, opts.Endpoint.Value); err != nil {
		ua := o.WithoutRedirect().WithRetries(0).WithPostCallback(authCallback)
		// No active session so try to create a new one
		_, err := jira.NewSession(ua, opts.Endpoint.Value, opts)
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
