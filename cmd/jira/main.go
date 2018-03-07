package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/coryb/figtree"
	"github.com/coryb/oreo"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracmd"
	"gopkg.in/op/go-logging.v1"
)

var (
	log = logging.MustGetLogger("jira")
)

func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(jiracli.Exit); ok {
			os.Exit(exit.Code)
		} else {
			fmt.Fprintf(os.Stderr, "%s\n%s", e, debug.Stack())
			os.Exit(1)
		}
	}
}

func main() {
	defer handleExit()

	jiracli.InitLogging()

	fig := figtree.NewFigTree()
	fig.EnvPrefix = "JIRA"
	fig.ConfigDir = ".jira.d"

	if err := os.MkdirAll(filepath.Join(jiracli.Homedir(), fig.ConfigDir), 0755); err != nil {
		log.Errorf("%s", err)
		panic(jiracli.Exit{Code: 1})
	}

	o := oreo.New().WithCookieFile(filepath.Join(jiracli.Homedir(), fig.ConfigDir, "cookies.js"))

	jiracmd.RegisterAllCommands()

	app := jiracli.CommandLine(fig, o)
	jiracli.ParseCommandLine(app, os.Args[1:])
}
