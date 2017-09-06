// +build !windows

package jiracli

import "github.com/tmc/keyring"

func keyringGet(user string) (string, error) {
	password, err := keyring.Get("go-jira", user)
	if err != nil && err != keyring.ErrNotFound {
		return password, err
	}
	return password, nil
}

func keyringSet(user, passwd string) error {
	return keyring.Set("go-jira", user, passwd)
}
