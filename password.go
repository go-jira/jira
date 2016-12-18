package jira

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/tmc/keyring"
)

func (c *Cli) GetPass(user string) string {
	passwd := ""
	if source, ok := c.opts["password-source"].(string); ok {
		if source == "keyring" {
			passwd, _ = keyring.Get("go-jira", user)
		} else if source == "pass" {
			if bin, err := exec.LookPath("pass"); err == nil {
				buf := bytes.NewBufferString("")
				cmd := exec.Command(bin, fmt.Sprintf("GoJira/%s", user))
				cmd.Stdout = buf
				cmd.Stderr = buf
				if err := cmd.Run(); err == nil {
					passwd = strings.TrimSpace(buf.String())
				}
			}
		} else {
			log.Warningf("Unknown password-source: %s", source)
		}
	}

	if passwd != "" {
		return passwd
	}
	fmt.Printf("Jira Password [%s]: ", user)
	pw, err := gopass.GetPasswdMasked()
	if err != nil {
		return ""
	}
	passwd = string(pw)
	return passwd
}

func (c *Cli) SetPass(user, passwd string) error {
	if source, ok := c.opts["password-source"].(string); ok {
		if source == "keyring" {
			// save password in keychain so that it can be used for subsequent http requests
			err := keyring.Set("go-jira", user, passwd)
			if err != nil {
				log.Errorf("Failed to set password in keyring: %s", err)
				return err
			}
		} else if source == "pass" {
			if bin, err := exec.LookPath("pass"); err == nil {
				in := bytes.NewBufferString(fmt.Sprintf("%s\n%s\n", passwd, passwd))
				out := bytes.NewBufferString("")
				cmd := exec.Command(bin, "insert", fmt.Sprintf("GoJira/%s", user))
				cmd.Stdin = in
				cmd.Stdout = out
				cmd.Stderr = out
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("Failed to insert password: %s", out.String())
				}
			}
		} else {
			return fmt.Errorf("Unknown password-source: %s", source)
		}
	}
	return nil
}
