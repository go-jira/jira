package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/coryb/figtree"
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

	fig := figtree.NewFigTree()
	fig.EnvPrefix = "JIRA"
	fig.ConfigDir = ".jira.d"

	o := oreo.New().WithCookieFile(filepath.Join(jiracli.Homedir(), fig.ConfigDir, "cookies.js"))
	o = o.WithPostCallback(
		func(req *http.Request, resp *http.Response) (*http.Response, error) {
			if resp.Header.Get("X-Ausername") == "anonymous" {
				// we are not logged in, so force login now by running the "login" command
				app.Parse([]string{"login"})
				return o.Do(req)
			}
			return resp, nil
		},
	)

	registry := []jiracli.CommandRegistry{
		jiracli.CommandRegistry{
			Command: "login",
			Entry:   jiracli.CmdLoginRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "logout",
			Entry:   jiracli.CmdLogoutRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "list",
			Aliases: []string{"ls"},
			Entry:   jiracli.CmdListRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "view",
			Entry:   jiracli.CmdViewRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "create",
			Entry:   jiracli.CmdCreateRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "edit",
			Entry:   jiracli.CmdEditRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "comment",
			Entry:   jiracli.CmdCommentRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "worklog list",
			Entry:   jiracli.CmdWorklogListRegistry(fig, o),
			Default: true,
		},
		jiracli.CommandRegistry{
			Command: "worklog add",
			Entry:   jiracli.CmdWorklogAddRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "fields",
			Entry:   jiracli.CmdFieldsRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "createmeta",
			Entry:   jiracli.CmdCreateMetaRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "editmeta",
			Entry:   jiracli.CmdEditMetaRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "subtask",
			Entry:   jiracli.CmdSubtaskRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "dup",
			Entry:   jiracli.CmdDupRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "block",
			Entry:   jiracli.CmdBlockRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "issuelink",
			Entry:   jiracli.CmdIssueLinkRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "issuelinktypes",
			Entry:   jiracli.CmdIssueLinkTypesRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "transition",
			Aliases: []string{"trans"},
			Entry:   jiracli.CmdTransitionRegistry(fig, o, ""),
		},
		jiracli.CommandRegistry{
			Command: "transitions",
			Entry:   jiracli.CmdTransitionsRegistry(fig, o, "transitions"),
		},
		jiracli.CommandRegistry{
			Command: "transmeta",
			Entry:   jiracli.CmdTransitionsRegistry(fig, o, "debug"),
		},
		jiracli.CommandRegistry{
			Command: "close",
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "close"),
		},
		jiracli.CommandRegistry{
			Command: "acknowledge",
			Aliases: []string{"ack"},
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "acknowledge"),
		},
		jiracli.CommandRegistry{
			Command: "reopen",
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "reopen"),
		},
		jiracli.CommandRegistry{
			Command: "resolve",
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "resolve"),
		},
		jiracli.CommandRegistry{
			Command: "start",
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "start"),
		},
		jiracli.CommandRegistry{
			Command: "stop",
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "stop"),
		},
		jiracli.CommandRegistry{
			Command: "todo",
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "To Do"),
		},
		jiracli.CommandRegistry{
			Command: "backlog",
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "Backlog"),
		},
		jiracli.CommandRegistry{
			Command: "done",
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "Done"),
		},
		jiracli.CommandRegistry{
			Command: "in-progress",
			Aliases: []string{"prog", "progress"},
			Entry:   jiracli.CmdTransitionRegistry(fig, o, "Progress"),
		},
		jiracli.CommandRegistry{
			Command: "vote",
			Entry:   jiracli.CmdVoteRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "rank",
			Entry:   jiracli.CmdRankRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "watch",
			Entry:   jiracli.CmdWatchRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "labels add",
			Entry:   jiracli.CmdLabelsAddRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "labels set",
			Entry:   jiracli.CmdLabelsAddRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "labels remove",
			Entry:   jiracli.CmdLabelsAddRegistry(fig, o),
			Aliases: []string{"rm"},
		},
		jiracli.CommandRegistry{
			Command: "take",
			Entry:   jiracli.CmdTakeRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "assign",
			Entry:   jiracli.CmdAssignRegistry(fig, o),
			Aliases: []string{"give"},
		},
		jiracli.CommandRegistry{
			Command: "unassign",
			Entry:   jiracli.CmdUnassignRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "component add",
			Entry:   jiracli.CmdComponentAddRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "components",
			Entry:   jiracli.CmdComponentsRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "issuetypes",
			Entry:   jiracli.CmdIssueTypesRegistry(fig, o),
		},
		jiracli.CommandRegistry{
			Command: "export-templates",
			Entry:   jiracli.CmdExportTemplatesRegistry(fig),
		},
		jiracli.CommandRegistry{
			Command: "unexport-templates",
			Entry:   jiracli.CmdUnexportTemplatesRegistry(fig),
		},
		jiracli.CommandRegistry{
			Command: "browse",
			Entry:   jiracli.CmdBrowseRegistry(fig),
			Aliases: []string{"b"},
		},
		jiracli.CommandRegistry{
			Command: "request",
			Entry:   jiracli.CmdRequestRegistry(fig, o),
			Aliases: []string{"req"},
		},
	}

	jiracli.Register(app, registry)

	app.Terminate(func(status int) {
		for _, arg := range os.Args {
			if arg == "-h" || arg == "--help" || len(os.Args) == 1 {
				panic(jiracli.Exit{Code: 0})
			}
		}
		panic(jiracli.Exit{Code: 1})
	})
	if _, err := app.Parse(os.Args[1:]); err != nil {
		ctx, _ := app.ParseContext(os.Args[1:])
		if ctx != nil {
			app.UsageForContext(ctx)
		}
		log.Errorf("Invalid Usage: %s", err)
		panic(jiracli.Exit{Code: 1})
	}
}
