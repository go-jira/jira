// +build !windows

package jiracli

import "github.com/tmc/keyring"

func keyringGet(user string) (string, error) {
	return keyring.Get("go-jira", user)
}

func keyringSet(user, passwd string) error {
	return keyring.Set("go-jira", user, passwd)
}
