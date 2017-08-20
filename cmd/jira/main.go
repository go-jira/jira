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

	// Usage:
	//   jira export-templates [-d DIR] [-t template]
	//   jira (b|browse) ISSUE
	//   jira request [-M METHOD] URI [DATA]
	//   jira ISSUE

	// General Options:
	//   -b --browse         Open your browser to the Jira issue
	//   -e --endpoint=URI   URI to use for jira
	//   -k --insecure       disable TLS certificate verification
	//   -h --help           Show this usage
	//   -t --template=FILE  Template file to use for output/editing
	//   -u --user=USER      Username to use for authenticaion (default: %s)
	//   -v --verbose        Increase output logging
	//   --unixproxy=PATH    Path for a unix-socket proxy (eg., --unixproxy /tmp/proxy.sock)
	//   --version           Print version

	// Query Options:
	//   -a --assignee=USER        Username assigned the issue
	//   -c --component=COMPONENT  Component to Search for
	//   -f --queryfields=FIELDS   Fields that are used in "list" template: (default: %s)
	//   -i --issuetype=ISSUETYPE  The Issue Type
	//   -l --limit=VAL            Maximum number of results to return in query (default: %d)
	//   --start=START             Start parameter for pagination
	//   -p --project=PROJECT      Project to Search for
	//   -q --query=JQL            Jira Query Language expression for the search
	//   -r --reporter=USER        Reporter to search for
	//   -s --sort=ORDER           For list operations, sort issues (default: %s)

	// Edit Options:
	//   -m --comment=COMMENT      Comment message for transition
	//   -o --override=KEY=VAL     Set custom key/value pairs

	// Create Options:
	//   -i --issuetype=ISSUETYPE  Jira Issue Type (default: Bug)
	//   -m --comment=COMMENT      Comment message for transition
	//   -o --override=KEY=VAL     Set custom key/value pairs

	// Worklog Options:
	//   -T --time-spent=TIMESPENT Time spent working on issue
	//   -m --comment=COMMENT      Comment message for worklog

	// Command Options:
	//   -d --directory=DIR        Directory to export templates to (default: %s)
	// `, user, defaultQueryFields, defaultMaxResults, defaultSort, user, fmt.Sprintf("%s/.jira.d/templates", home))
	// 		printer(output)
	// 	}

	// 	jiraCommands := map[string]string{
	// 		"export-templates": "export-templates",
	// 		"browse":           "browse",
	// 		"req":              "request",
	// 		"request":          "request",
	// 	}

	// 	defaults := map[string]interface{}{
	// 		"user":        user,
	// 		"queryfields": defaultQueryFields,
	// 		"directory":   fmt.Sprintf("%s/.jira.d/templates", home),
	// 		"sort":        defaultSort,
	// 		"max_results": defaultMaxResults,
	// 		"method":      "GET",
	// 		"quiet":       false,
	// 	}
	// 	opts := make(map[string]interface{})

	// 	setopt := func(name string, value interface{}) {
	// 		opts[name] = value
	// 	}

	// 	op := optigo.NewDirectAssignParser(map[string]interface{}{
	// 		"h|help": usage,
	// 		"version": func() {
	// 			fmt.Println(fmt.Sprintf("version: %s", jira.VERSION))
	// 			os.Exit(0)
	// 		},
	// 		"v|verbose+": func() {
	// 			logging.SetLevel(logging.GetLevel("")+1, "")
	// 		},
	// 		"dryrun":                setopt,
	// 		"b|browse":              setopt,
	// 		"editor=s":              setopt,
	// 		"u|user=s":              setopt,
	// 		"endpoint=s":            setopt,
	// 		"k|insecure":            setopt,
	// 		"t|template=s":          setopt,
	// 		"q|query=s":             setopt,
	// 		"p|project=s":           setopt,
	// 		"c|component=s":         setopt,
	// 		"a|assignee=s":          setopt,
	// 		"i|issuetype=s":         setopt,
	// 		"remove":                setopt,
	// 		"r|reporter=s":          setopt,
	// 		"f|queryfields=s":       setopt,
	// 		"x|expand=s":            setopt,
	// 		"s|sort=s":              setopt,
	// 		"l|limit|max_results=i": setopt,
	// 		"start|start_at=i":      setopt,
	// 		"o|override=s%":         &opts,
	// 		"noedit":                setopt,
	// 		"edit":                  setopt,
	// 		"m|comment=s":           setopt,
	// 		"d|dir|directory=s":     setopt,
	// 		"M|method=s":            setopt,
	// 		"S|saveFile=s":          setopt,
	// 		"T|time-spent=s":        setopt,
	// 		"Q|quiet":               setopt,
	// 		"unixproxy":             setopt,
	// 		"down":                  setopt,
	// 		"default":               setopt,
	// 	})

	// 	if err := op.ProcessAll(os.Args[1:]); err != nil {
	// 		log.Errorf("%s", err)
	// 		usage(false)
	// 	}
	// 	args := op.Args

	// 	var command string
	// 	if len(args) > 0 {
	// 		if alias, ok := jiraCommands[args[0]]; ok {
	// 			command = alias
	// 			args = args[1:]
	// 		} else if len(args) > 1 {
	// 			// look at second arg for "dups" and "blocks" commands
	// 			// also for 'set/add/remove' actions like 'labels'
	// 			if alias, ok := jiraCommands[args[1]]; ok {
	// 				command = alias
	// 				args = append(args[:1], args[2:]...)
	// 			}
	// 		}
	// 	}

	// 	if command == "" && len(args) > 0 {
	// 		command = args[0]
	// 		args = args[1:]
	// 	}

	// 	os.Setenv("JIRA_OPERATION", command)
	// 	loadConfigs(opts)

	// 	// check to see if it was set in the configs:
	// 	if value, ok := opts["command"].(string); ok {
	// 		command = value
	// 	} else if _, ok := jiraCommands[command]; !ok || command == "" {
	// 		if command != "" {
	// 			args = append([]string{command}, args...)
	// 		}
	// 		command = "view"
	// 	}

	// 	// apply defaults
	// 	for k, v := range defaults {
	// 		if _, ok := opts[k]; !ok {
	// 			log.Debugf("Setting %q to %#v from defaults", k, v)
	// 			opts[k] = v
	// 		}
	// 	}

	// 	log.Debugf("opts: %v", opts)
	// 	log.Debugf("args: %v", args)

	// 	if _, ok := opts["endpoint"]; !ok {
	// 		log.Errorf("endpoint option required.  Either use --endpoint or set a endpoint option in your ~/.jira.d/config.yml file")
	// 		os.Exit(1)
	// 	}

	// 	c := jira.New(opts)

	// 	log.Debugf("opts: %s", opts)

	// 	setEditing := func(dflt bool) {
	// 		log.Debugf("Default Editing: %t", dflt)
	// 		if dflt {
	// 			if val, ok := opts["noedit"].(bool); ok && val {
	// 				log.Debugf("Setting edit = false")
	// 				opts["edit"] = false
	// 			} else {
	// 				log.Debugf("Setting edit = true")
	// 				opts["edit"] = true
	// 			}
	// 		} else {
	// 			if _, ok := opts["edit"].(bool); !ok {
	// 				log.Debugf("Setting edit = %t", dflt)
	// 				opts["edit"] = dflt
	// 			}
	// 		}
	// 	}

	// 	requireArgs := func(count int) {
	// 		if len(args) < count {
	// 			log.Errorf("Not enough arguments. %d required, %d provided", count, len(args))
	// 			usage(false)
	// 		}
	// 	}

	// 	var err error
	// 	switch command {
	// 	case "browse":
	// 		requireArgs(1)
	// 		opts["browse"] = true
	// 		err = c.Browse(args[0])
	// 	case "export-templates":
	// 		err = c.CmdExportTemplates()
	// 	case "request":
	// 		requireArgs(1)
	// 		data := ""
	// 		if len(args) > 1 {
	// 			data = args[1]
	// 		}
	// 		err = c.CmdRequest(args[0], data)
	// 	default:
	// 		log.Errorf("Unknown command %s", command)
	// 		os.Exit(1)
	// 	}

	// 	if err != nil {
	// 		log.Errorf("%s", err)
	// 		os.Exit(1)
	// 	}
	// 	os.Exit(0)
}

// func parseYaml(file string, opts map[string]interface{}) {
// 	if fh, err := ioutil.ReadFile(file); err == nil {
// 		log.Debugf("Found Config file: %s", file)
// 		if err := yaml.Unmarshal(fh, &opts); err != nil {
// 			log.Errorf("Unable to parse %s: %s", file, err)
// 		}
// 	}
// }

// func populateEnv(opts map[string]interface{}) {
// 	for k, v := range opts {
// 		envName := fmt.Sprintf("JIRA_%s", strings.ToUpper(k))
// 		var val string
// 		switch t := v.(type) {
// 		case string:
// 			val = t
// 		case int, int8, int16, int32, int64:
// 			val = fmt.Sprintf("%d", t)
// 		case float32, float64:
// 			val = fmt.Sprintf("%f", t)
// 		case bool:
// 			val = fmt.Sprintf("%t", t)
// 		default:
// 			val = fmt.Sprintf("%v", t)
// 		}
// 		os.Setenv(envName, val)
// 	}
// }

// func loadConfigs(opts map[string]interface{}) {
// 	populateEnv(opts)
// 	paths := jira.FindParentPaths(".jira.d/config.yml")
// 	// prepend
// 	paths = append(paths, "/etc/go-jira.yml")

// 	// iterate paths in reverse
// 	for i := 0; i < len(paths); i++ {
// 		file := paths[i]
// 		if stat, err := os.Stat(file); err == nil {
// 			tmp := make(map[string]interface{})
// 			// check to see if config file is exectuable
// 			if stat.Mode()&0111 == 0 {
// 				parseYaml(file, tmp)
// 			} else {
// 				log.Debugf("Found Executable Config file: %s", file)
// 				// it is executable, so run it and try to parse the output
// 				cmd := exec.Command(file)
// 				stdout := bytes.NewBufferString("")
// 				cmd.Stdout = stdout
// 				cmd.Stderr = bytes.NewBufferString("")
// 				if err := cmd.Run(); err != nil {
// 					log.Errorf("%s is exectuable, but it failed to execute: %s\n%s", file, err, cmd.Stderr)
// 					os.Exit(1)
// 				}
// 				yaml.Unmarshal(stdout.Bytes(), &tmp)
// 			}
// 			for k, v := range tmp {
// 				if _, ok := opts[k]; !ok {
// 					log.Debugf("Setting %q to %#v from %s", k, v, file)
// 					opts[k] = v
// 				}
// 			}
// 			populateEnv(opts)
// 		}
// 	}
// }
