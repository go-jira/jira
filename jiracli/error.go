package jiracli

import "github.com/pkg/errors"

type Error struct {
	error
}

func CliError(cause error) error {
	return &Error{
		errors.WithStack(cause),
	}
}
