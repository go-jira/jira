package jiracli

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiradata"
)

func (o *GlobalOptions) ProvideAuthParams() *jiradata.AuthParams {
	return &jiradata.AuthParams{
		Username: o.User.Value,
		Password: o.GetPass(),
	}
}

func (o *GlobalOptions) GetPass() string {
	passwd := ""
	if o.PasswordSource.Value != "" {
		if o.PasswordSource.Value == "keyring" {
			var err error
			passwd, err = keyringGet(o.User.Value)
			if err != nil {
				panic(err)
			}
		} else if o.PasswordSource.Value == "pass" {
			if bin, err := exec.LookPath("pass"); err == nil {
				buf := bytes.NewBufferString("")
				cmd := exec.Command(bin, fmt.Sprintf("GoJira/%s", o.User))
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
	survey.AskOne(
		&survey.Password{
			Message: fmt.Sprintf("Jira Password [%s]: ", o.User),
		},
		&passwd,
		nil,
	)
	o.SetPass(passwd)
	return passwd
}

func (o *GlobalOptions) SetPass(passwd string) error {
	if o.PasswordSource.Value == "keyring" {
		// save password in keychain so that it can be used for subsequent http requests
		err := keyringSet(o.User.Value, passwd)
		if err != nil {
			log.Errorf("Failed to set password in keyring: %s", err)
			return err
		}
	} else if o.PasswordSource.Value == "pass" {
		if bin, err := exec.LookPath("pass"); err == nil {
			log.Debugf("using %s", bin)
			passName := fmt.Sprintf("GoJira/%s", o.User)
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
