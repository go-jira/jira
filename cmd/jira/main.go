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

type oreoLogger struct {
	logger *logging.Logger
}

var log = logging.MustGetLogger("jira")

func (ol *oreoLogger) Printf(format string, args ...interface{}) {
	ol.logger.Debugf(format, args...)
}

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

	configDir := ".jira.d"
	fig := figtree.NewFigTree(
		figtree.WithHome(jiracli.Homedir()),
		figtree.WithEnvPrefix("JIRA"),
		figtree.WithConfigDir(configDir),
	)

	if err := os.MkdirAll(filepath.Join(jiracli.Homedir(), configDir), 0755); err != nil {
		log.Errorf("%s", err)
		panic(jiracli.Exit{Code: 1})
	}

	o := oreo.New().WithCookieFile(filepath.Join(jiracli.Homedir(), configDir, "cookies.js")).WithLogger(&oreoLogger{log})

	jiracmd.RegisterAllCommands()

	app := jiracli.CommandLine(fig, o)
	jiracli.ParseCommandLine(app, os.Args[1:])
}
