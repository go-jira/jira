package keyring

import (
	"fmt"
	"os/exec"
	"regexp"
	"syscall"
)

type osxProvider struct {
}

var pwRe = regexp.MustCompile("password: \"(.+)\"")

func (p osxProvider) Get(Service, Username string) (string, error) {
	args := []string{"find-generic-password",
		"-s", Service,
		"-a", Username,
		"-g"}
	c := exec.Command("/usr/bin/security", args...)
	o, err := c.CombinedOutput()
	if err != nil {
		exitCode := c.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
		// check particular exit code
		if exitCode == 44 {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("/usr/bin/security: %s", err)
	}
	matches := pwRe.FindStringSubmatch(string(o))
	if len(matches) != 2 {
		return "", ErrNotFound
	}
	return matches[1], nil
}

func (p osxProvider) Set(Service, Username, Password string) error {
	args := []string{"add-generic-password",
		"-s", Service,
		"-a", Username,
		"-w", Password,
		"-U"}
	c := exec.Command("/usr/bin/security", args...)
	err := c.Run()
	if err != nil {
		o, _ := c.CombinedOutput()
		return fmt.Errorf(string(o))
	}
	return nil
}

func initializeProvider() (provider, error) {
	return osxProvider{}, nil
}
