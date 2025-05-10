package jiracli

import "fmt"

func keyringGet(user string) (string, error) {
	return "", fmt.Errorf("Keyring is not supported for Windows, see: https://github.com/dougEfresh/keyring")
}

func keyringSet(user, passwd string) error {
	return fmt.Errorf("Keyring is not supported for Windows, see: https://github.com/dougEfresh/keyring")
}
