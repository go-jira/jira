package jiracli

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/go-jira/jira/jiradata"
	"gopkg.in/AlecAivazis/survey.v1"
)

func (o *GlobalOptions) ProvideAuthParams() *jiradata.AuthParams {
	return &jiradata.AuthParams{
		Username: o.Login.Value,
		Password: o.GetPass(),
	}
}

func (o *GlobalOptions) keyName() string {
	user := o.Login.Value
	if o.AuthMethodIsToken() {
		user = "api-token:" + user
	}

	if o.PasswordSource.Value == "pass" {
		if o.PasswordName.Value != "" {
			return o.PasswordName.Value
		}
		return fmt.Sprintf("GoJira/%s", user)
	}

	if o.PasswordSource.Value == "gopass" {
		if o.PasswordName.Value != "" {
			return o.PasswordName.Value
		}
		return fmt.Sprintf("GoJira/%s", user)
	}
	return user
}

func (o *GlobalOptions) GetPasswordPath() string {
	// if no password source path then just default
	// to the password source name
	if o.PasswordSourcePath.Value == "" {
		return o.PasswordSource.Value
	}
	return o.PasswordSourcePath.Value
}

func (o *GlobalOptions) GetPass() string {
	if o.cachedPassword != "" {
		return o.cachedPassword
	}
	log.Debugf("Getting Password")
	if o.PasswordSource.Value != "" {
		log.Debugf("password-source: %s", o.PasswordSource)
		if o.PasswordSource.Value == "keyring" {
			log.Info("Querying keyring password source.")
			var err error
			o.cachedPassword, err = keyringGet(o.keyName())
			if err != nil {
				panic(err)
			}
		} else if o.PasswordSource.Value == "gopass" {
			log.Debugf("Querying gopass password source.")
			binary := o.GetPasswordPath()
			if o.PasswordDirectory.Value != "" {
				orig := os.Getenv("PASSWORD_STORE_DIR")
				log.Debugf("using password-directory: %s", o.PasswordDirectory)
				os.Setenv("PASSWORD_STORE_DIR", o.PasswordDirectory.Value)
				defer os.Setenv("PASSWORD_STORE_DIR", orig)
			}
			if passDir := os.Getenv("PASSWORD_STORE_DIR"); passDir != "" {
				log.Debugf("using PASSWORD_STORE_DIR=%s", passDir)
			}
			if bin, err := exec.LookPath(binary); err == nil {
				log.Debugf("found gopass at: %s", bin)
				buf := bytes.NewBufferString("")
				cmd := exec.Command(bin, "show", "-o", o.keyName())
				cmd.Stdout = buf
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err == nil {
					o.cachedPassword = strings.TrimSpace(buf.String())
				} else {
					log.Warningf("gopass command failed with:\n%s", buf.String())
				}
			} else {
				log.Warning("Gopass binary was not found! Fallback to default password behaviour!")
			}
		} else if o.PasswordSource.Value == "pass" {
			log.Debugf("Querying pass password source.")
			binary := o.GetPasswordPath()
			if o.PasswordDirectory.Value != "" {
				orig := os.Getenv("PASSWORD_STORE_DIR")
				log.Debugf("using password-directory: %s", o.PasswordDirectory)
				os.Setenv("PASSWORD_STORE_DIR", o.PasswordDirectory.Value)
				defer os.Setenv("PASSWORD_STORE_DIR", orig)
			}
			if passDir := os.Getenv("PASSWORD_STORE_DIR"); passDir != "" {
				log.Debugf("using PASSWORD_STORE_DIR=%s", passDir)
			}
			if bin, err := exec.LookPath(binary); err == nil {
				log.Debugf("found pass at: %s", bin)
				buf := bytes.NewBufferString("")
				cmd := exec.Command(bin, o.keyName())
				cmd.Stdout = buf
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err == nil {
					o.cachedPassword = strings.TrimSpace(buf.String())
				} else {
					log.Warningf("pass command failed with:\n%s", buf.String())
				}
			} else {
				log.Warning("pass binary was not found! Fallback to default password behaviour!")
			}
		} else if o.PasswordSource.Value == "stdin" {
			log.Info("Reading password from stdin.")
			allBytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic(fmt.Sprintf("unable to read bytes from stdin: %s", err))
			}
			o.cachedPassword = string(allBytes)
		} else {
			log.Warningf("Unknown password-source: %s", o.PasswordSource)
		}
	}

	if o.cachedPassword != "" {
		log.Info("Password cached.")
		return o.cachedPassword
	}

	if o.cachedPassword = os.Getenv("JIRA_API_TOKEN"); o.cachedPassword != "" && o.AuthMethodIsToken() {
		return o.cachedPassword
	}

	prompt := fmt.Sprintf("Jira Password [%s]: ", o.Login)
	help := ""

	if o.AuthMethodIsToken() {
		prompt = fmt.Sprintf("Jira API-Token [%s]: ", o.Login)
		help = "API Tokens may be required by your Jira service endpoint: https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-basic-auth-and-cookie-based-auth/"
	}

	err := survey.AskOne(
		&survey.Password{
			Message: prompt,
			Help:    help,
		},
		&o.cachedPassword,
		nil,
	)
	if err != nil {
		log.Errorf("%s", err)
		panic(Exit{Code: 1})
	}
	o.SetPass(o.cachedPassword)
	return o.cachedPassword
}

func (o *GlobalOptions) SetPass(passwd string) error {
	// dont reset password to empty string
	if passwd == "" {
		return nil
	}

	if o.PasswordSource.Value == "keyring" {
		// save password in keychain so that it can be used for subsequent http requests
		err := keyringSet(o.keyName(), passwd)
		if err != nil {
			log.Errorf("Failed to set password in keyring: %s", err)
			return err
		}
	} else if o.PasswordSource.Value == "gopass" {
		if o.PasswordDirectory.Value != "" {
			orig := os.Getenv("PASSWORD_STORE_DIR")
			os.Setenv("PASSWORD_STORE_DIR", o.PasswordDirectory.Value)
			defer os.Setenv("PASSWORD_STORE_DIR", orig)
		}
		if bin, err := exec.LookPath("gopass"); err == nil {
			log.Debugf("using %s", bin)
			passName := o.keyName()
			if passwd != "" {
				in := bytes.NewBufferString(fmt.Sprintf("%s\n", passwd))
				out := bytes.NewBufferString("")
				cmd := exec.Command(bin, "insert", "--force", passName)
				cmd.Stdin = in
				cmd.Stdout = out
				cmd.Stderr = out
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("Failed to insert password: %s", out.String())
				}
			}
		} else {
			return fmt.Errorf("Gopass binary not found!")
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
			return fmt.Errorf("Pass binary not found!")
		}
	} else if o.PasswordSource.Value != "" {
		return fmt.Errorf("Unknown password-source: %s", o.PasswordSource)
	}
	return nil
}
