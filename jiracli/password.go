package jiracli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

func (o *GlobalOptions) ProvideAuthParams() *jiradata.AuthParams {
	return &jiradata.AuthParams{
		Username: o.Login.Value,
		Password: o.GetPass(),
	}
}

func (o *GlobalOptions) keyName() string {
	user := o.Login.Value
	if o.AuthMethod() == "api-token" {
		user = "api-token:" + user
	}

	if o.PasswordSource.Value == "pass" {
		if o.PasswordName.Value != "" {
			return o.PasswordName.Value
		}
		return fmt.Sprintf("GoJira/%s", user)
	}
	return user
}

func (o *GlobalOptions) GetPass() string {
	passwd := ""
	if o.PasswordSource.Value != "" {
		if o.PasswordSource.Value == "keyring" {
			var err error
			passwd, err = keyringGet(o.keyName())
			if err != nil {
				panic(err)
			}
		} else if o.PasswordSource.Value == "pass" {
			if o.PasswordDirectory.Value != "" {
				orig := os.Getenv("PASSWORD_STORE_DIR")
				os.Setenv("PASSWORD_STORE_DIR", o.PasswordDirectory.Value)
				defer os.Setenv("PASSWORD_STORE_DIR", orig)
			}
			if bin, err := exec.LookPath("pass"); err == nil {
				buf := bytes.NewBufferString("")
				cmd := exec.Command(bin, o.keyName())
				cmd.Stdout = buf
				cmd.Stderr = buf
				if err := cmd.Run(); err == nil {
					passwd = strings.TrimSpace(buf.String())
				}
			}
		} else {
			log.Warningf("Unknown password-source: %s", o.PasswordSource)
		}
	}

	if passwd != "" {
		return passwd
	}

	if passwd = os.Getenv("JIRA_API_TOKEN"); passwd != "" && o.AuthMethod() == "api-token" {
		return passwd
	}

	prompt := fmt.Sprintf("Jira Password [%s]: ", o.Login)
	help := ""

	if o.AuthMethod() == "api-token" {
		prompt = fmt.Sprintf("Jira API-Token [%s]: ", o.Login)
		help = "API Tokens may be required by your Jira service endpoint: https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-basic-auth-and-cookie-based-auth/"
	}

	err := survey.AskOne(
		&survey.Password{
			Message: prompt,
			Help:    help,
		},
		&passwd,
		nil,
	)
	if err != nil {
		log.Errorf("%s", err)
		panic(Exit{Code: 1})
	}
	o.SetPass(passwd)
	return passwd
}

func (o *GlobalOptions) SetPass(passwd string) error {
	if o.PasswordSource.Value == "keyring" {
		// save password in keychain so that it can be used for subsequent http requests
		err := keyringSet(o.keyName(), passwd)
		if err != nil {
			log.Errorf("Failed to set password in keyring: %s", err)
			return err
		}
	} else if o.PasswordSource.Value == "pass" {
		if o.PasswordDirectory.Value != "" {
			orig := os.Getenv("PASSWORD_STORE_DIR")
			os.Setenv("PASSWORD_STORE_DIR", o.PasswordDirectory.Value)
			defer os.Setenv("PASSWORD_STORE_DIR", orig)
		}
		if bin, err := exec.LookPath("pass"); err == nil {
			log.Debugf("using %s", bin)
			passName := o.keyName()
			if passwd != "" {
				in := bytes.NewBufferString(fmt.Sprintf("%s\n%s\n", passwd, passwd))
				out := bytes.NewBufferString("")
				cmd := exec.Command(bin, "insert", "--force", passName)
				cmd.Stdin = in
				cmd.Stdout = out
				cmd.Stderr = out
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("Failed to insert password: %s", out.String())
				}
			} else {
				// clear the `pass` entry on empty password
				if err := exec.Command(bin, "rm", "--force", passName).Run(); err != nil {
					return fmt.Errorf("Failed to clear password for %s", passName)
				}
			}
		}
	} else if o.PasswordSource.Value != "" {
		return fmt.Errorf("Unknown password-source: %s", o.PasswordSource)
	}
	return nil
}
