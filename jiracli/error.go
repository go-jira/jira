package jiracli

import "github.com/pkg/errors"

type Error struct {
	error
}

func cliError(cause error) error {
	return &Error{
		errors.WithStack(cause),
	}
}
