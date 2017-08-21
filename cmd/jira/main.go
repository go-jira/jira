package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/coryb/oreo"

	jira "gopkg.in/Netflix-Skunkworks/go-jira.v1"
	"gopkg.in/Netflix-Skunkworks/go-jira.v1/jiracli"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/op/go-logging.v1"
)

var (
	log           = logging.MustGetLogger("jira")
	defaultFormat = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
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
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	format := os.Getenv("JIRA_LOG_FORMAT")
	if format == "" {
		format = defaultFormat
	}
	logging.SetBackend(
		logging.NewBackendFormatter(
			logBackend,
			logging.MustStringFormatter(format),
		),
	)
	if os.Getenv("JIRA_DEBUG") == "" {
		logging.SetLevel(logging.NOTICE, "")
	} else {
		logging.SetLevel(logging.DEBUG, "")
	}

	app := kingpin.New("jira", "Jira Command Line Interface")
	app.Command("version", "Prints version").PreAction(func(*kingpin.ParseContext) error {
		fmt.Println(jira.VERSION)
		panic(jiracli.Exit{Code: 0})
	})

	app.Flag("verbose", "Increase verbosity for debugging").Short('v').PreAction(func(_ *kingpin.ParseContext) error {
		logging.SetLevel(logging.GetLevel("")+1, "")
		if logging.GetLevel("") > logging.DEBUG {
			oreo.TraceRequestBody = true
			oreo.TraceResponseBody = true
		}
		return nil
	}).Counter()

	cli := jiracli.New(".jira.d")

	registry := []jiracli.CommandRegistry{
		jiracli.CommandRegistry{
			Command: "login",
			Entry:   cli.CmdLoginRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "logout",
			Entry:   cli.CmdLogoutRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "list",
			Aliases: []string{"ls"},
			Entry:   cli.CmdListRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "view",
			Entry:   cli.CmdViewRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "create",
			Entry:   cli.CmdCreateRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "edit",
			Entry:   cli.CmdEditRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "comment",
			Entry:   cli.CmdCommentRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "worklog list",
			Entry:   cli.CmdWorklogListRegistry(),
			Default: true,
		},
		jiracli.CommandRegistry{
			Command: "worklog add",
			Entry:   cli.CmdWorklogAddRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "fields",
			Entry:   cli.CmdFieldsRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "createmeta",
			Entry:   cli.CmdCreateMetaRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "editmeta",
			Entry:   cli.CmdEditMetaRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "subtask",
			Entry:   cli.CmdSubtaskRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "dup",
			Entry:   cli.CmdDupRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "block",
			Entry:   cli.CmdBlockRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "issuelink",
			Entry:   cli.CmdIssueLinkRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "issuelinktypes",
			Entry:   cli.CmdIssueLinkTypesRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "transition",
			Aliases: []string{"trans"},
			Entry:   cli.CmdTransitionRegistry(""),
		},
		jiracli.CommandRegistry{
			Command: "transitions",
			Entry:   cli.CmdTransitionsRegistry("transitions"),
		},
		jiracli.CommandRegistry{
			Command: "transmeta",
			Entry:   cli.CmdTransitionsRegistry("debug"),
		},
		jiracli.CommandRegistry{
			Command: "close",
			Entry:   cli.CmdTransitionRegistry("close"),
		},
		jiracli.CommandRegistry{
			Command: "acknowledge",
			Aliases: []string{"ack"},
			Entry:   cli.CmdTransitionRegistry("acknowledge"),
		},
		jiracli.CommandRegistry{
			Command: "reopen",
			Entry:   cli.CmdTransitionRegistry("reopen"),
		},
		jiracli.CommandRegistry{
			Command: "resolve",
			Entry:   cli.CmdTransitionRegistry("resolve"),
		},
		jiracli.CommandRegistry{
			Command: "start",
			Entry:   cli.CmdTransitionRegistry("start"),
		},
		jiracli.CommandRegistry{
			Command: "stop",
			Entry:   cli.CmdTransitionRegistry("stop"),
		},
		jiracli.CommandRegistry{
			Command: "todo",
			Entry:   cli.CmdTransitionRegistry("To Do"),
		},
		jiracli.CommandRegistry{
			Command: "backlog",
			Entry:   cli.CmdTransitionRegistry("Backlog"),
		},
		jiracli.CommandRegistry{
			Command: "done",
			Entry:   cli.CmdTransitionRegistry("Done"),
		},
		jiracli.CommandRegistry{
			Command: "in-progress",
			Aliases: []string{"prog", "progress"},
			Entry:   cli.CmdTransitionRegistry("Progress"),
		},
		jiracli.CommandRegistry{
			Command: "vote",
			Entry:   cli.CmdVoteRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "rank",
			Entry:   cli.CmdRankRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "watch",
			Entry:   cli.CmdWatchRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "labels add",
			Entry:   cli.CmdLabelsAddRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "labels set",
			Entry:   cli.CmdLabelsAddRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "labels remove",
			Entry:   cli.CmdLabelsAddRegistry(),
			Aliases: []string{"rm"},
		},
		jiracli.CommandRegistry{
			Command: "take",
			Entry:   cli.CmdTakeRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "assign",
			Entry:   cli.CmdAssignRegistry(),
			Aliases: []string{"give"},
		},
		jiracli.CommandRegistry{
			Command: "unassign",
			Entry:   cli.CmdUnassignRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "component add",
			Entry:   cli.CmdComponentAddRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "components",
			Entry:   cli.CmdComponentsRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "issuetypes",
			Entry:   cli.CmdIssueTypesRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "export-templates",
			Entry:   cli.CmdExportTemplatesRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "unexport-templates",
			Entry:   cli.CmdUnexportTemplatesRegistry(),
		},
		jiracli.CommandRegistry{
			Command: "browse",
			Entry:   cli.CmdBrowseRegistry(),
			Aliases: []string{"b"},
		},
		jiracli.CommandRegistry{
			Command: "request",
			Entry:   cli.CmdRequestRegistry(),
			Aliases: []string{"req"},
		},
	}

	cli.Register(app, registry)

	app.Terminate(func(status int) {
		for _, arg := range os.Args {
			if arg == "-h" || arg == "--help" || len(os.Args) == 1 {
				panic(jiracli.Exit{Code: 0})
			}
		}
		panic(jiracli.Exit{Code: 1})
	})
	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("%s", err)
	}
}
