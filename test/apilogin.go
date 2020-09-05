package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// withApiLogin is a hack to provide an api token on every command, this means keyring
// and gpg is not necessary to run the testing suite.
//
// a buffer containing stdout will be returned to the caller if no error is encountered.
// this still expects a config file is in a parent where the test runs for project
// and endpoint details.
func withApiLogin(login string, token string, cmd *exec.Cmd) (bytes.Buffer, error) {
	var buf bytes.Buffer

	cmd.Args = append(cmd.Args, "--login", login)

	diag := fmt.Sprintf("--- running command: %+v ---\n", cmd.Args)
	io.WriteString(os.Stdout, diag)

	// write to stdout and also to our buffer
	out := io.MultiWriter(&buf, os.Stdout)
	cmd.Stdout = out

	e := io.MultiWriter(&buf, os.Stderr)
	cmd.Stderr = e

	in, err := cmd.StdinPipe()
	if err != nil {
		return buf, err
	}

	err = cmd.Start()
	if err != nil {
		return buf, err
	}

	_, err = io.WriteString(in, token)
	if err != nil {
		return buf, err
	}
	in.Close()

	err = cmd.Wait()
	if err != nil {
		return buf, err
	}

	diag = fmt.Sprintf("--- finished command: %+v ---\n\n", cmd.Args)
	io.WriteString(os.Stdout, diag)
	return buf, nil
}
