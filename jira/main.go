package main

import (
	"github.com/docopt/docopt-go"
	"github.com/op/go-logging"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/Netflix-Skunkworks/go-jira/jira/cli"
)

var log = logging.MustGetLogger("jira")
var format = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"

func main() {
	user := os.Getenv("USER")
	usage := fmt.Sprintf(`
Usage:
  jira [-v ...] [-u USER] [-e URI] [-t FILE] fields
  jira [-v ...] [-u USER] [-e URI] [-t FILE] login
  jira [-v ...] [-u USER] [-e URI] [-t FILE] ls [-q JQL]
  jira [-v ...] [-u USER] [-e URI] [-t FILE] view ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] editmeta ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] edit ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] issuetypes [-p PROJECT] 
  jira [-v ...] [-u USER] [-e URI] [-t FILE] createmeta [-p PROJECT] [-i ISSUETYPE] 
  jira [-v ...] [-u USER] [-e URI] [-t FILE] transitions ISSUE

  jira TODO [-v ...] [-u USER] [-e URI] [-t FILE] create [-p PROJECT] [-i ISSUETYPE]
  jira TODO [-v ...] [-u USER] [-e URI] DUPLICATE dups ISSUE
  jira TODO [-v ...] [-u USER] [-e URI] BLOCKER blocks ISSUE
  jira TODO [-v ...] [-u USER] [-e URI] close ISSUE [-m COMMENT]
  jira TODO [-v ...] [-u USER] [-e URI] resolve ISSUE [-m COMMENT]
  jira TODO [-v ...] [-u USER] [-e URI] comment ISSUE [-m COMMENT]
  jira TODO [-v ...] [-u USER] [-e URI] take ISSUE
  jira TODO [-v ...] [-u USER] [-e URI] assign ISSUE ASSIGNEE

General Options:
  -h --help           Show this usage
  --version           Show this version
  -v --verbose        Increase output logging
  -u --user=USER      Username to use for authenticaion (default: %s)
  -e --endpoint=URI   URI to use for jira (default: https://jira)
  -t --template=FILE  Template file to use for output

List Options:
  -q --query=JQL      Jira Query Language expression for the search

Create Options:
  -p --project=PROJECT      Jira Project Name
  -i --issuetype=ISSUETYPE  Jira Issue Type (default: Bug)
`, user)
	
	args, _ := docopt.Parse(usage, nil, true, "0.0.1", false, false)
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logging.SetBackend(
		logging.NewBackendFormatter(
			logBackend, 
			logging.MustStringFormatter(format),
		),
	)
	logging.SetLevel(logging.NOTICE, "")
	if verbose, ok := args["--verbose"]; ok {
		if verbose.(int) > 1 {
			logging.SetLevel(logging.DEBUG, "")
		} else if verbose.(int) > 0 { 
			logging.SetLevel(logging.INFO, "")
		}
	}

	log.Info("Args: %v", args)


	opts := make(map[string]string)
	loadConfigs(opts)

	// strip the "--" off the command line options
	// and populate the opts that we pass to the cli ctor
	for key,val := range args {
		if val != nil && strings.HasPrefix(key, "--") {
			opt := key[2:]
			switch v := val.(type) {
				// only deal with string opts, ignore
				// other types, like int (for now) since
				// they are only used for --verbose
			case string:
				opts[opt] = v
			}
		}
	}
	
	// cant use proper [default:x] syntax in docopt
	// because only want to default if the option is not
	// already specified in some .jira.d/config.yml file
	if _, ok := opts["endpoint"]; !ok {
		opts["endpoint"] = "https://jira"
	}
	if _, ok := opts["user"]; !ok {
		opts["user"] = user
	}
	if _, ok := opts["issuetype"]; !ok {
		opts["issuetype"] = "Bug"
	}
	
	c := cli.New(opts)

	log.Debug("opts: %s", opts);
	
	var err error
	if val, ok := args["login"]; ok && val.(bool) {
		err = c.CmdLogin()
	} else if val, ok := args["fields"]; ok && val.(bool) {
		err = c.CmdFields()
	} else if val, ok := args["ls"]; ok && val.(bool) {
		err = c.CmdList()
	} else if val, ok := args["edit"]; ok && val.(bool) {
		issue, _ := args["ISSUE"]
		err = c.CmdEdit(issue.(string))
	} else if val, ok := args["editmeta"]; ok && val.(bool) {
		issue, _ := args["ISSUE"]
		err = c.CmdEditMeta(issue.(string))
	} else if val, ok := args["issuetypes"]; ok && val.(bool) {
		var project interface{}
		if project, ok = opts["project"]; !ok {
			log.Error("missing PROJECT argument or \"project\" property in the config file")
			os.Exit(1)
		}
		err = c.CmdIssueTypes(project.(string))
	} else if val, ok := args["createmeta"]; ok && val.(bool) {
		var project interface{}
		if project, ok = opts["project"]; !ok {
			log.Error("missing PROJECT argument or \"project\" property in the config file")
			os.Exit(1)
		}
		var issuetype interface{}
		if issuetype, ok = opts["issuetype"]; !ok {
			issuetype = "Bug"
		}
		err = c.CmdCreateMeta(project.(string), issuetype.(string))
	} else if val, ok := args["transitions"]; ok && val.(bool) {
		issue, _ := args["ISSUE"]
		err = c.CmdTransitions(issue.(string))
	} else if val, ok := args["ISSUE"]; ok {
		err = c.CmdView(val.(string))
	}

	if err != nil { 
		os.Exit(1)
	}
	os.Exit(0)
}

func parseYaml(file string, opts map[string]string) {
	if fh, err := ioutil.ReadFile(file); err == nil {
		log.Debug("Found Config file: %s", file)
		yaml.Unmarshal(fh, &opts)
	}
}

func loadConfigs(opts map[string]string) {
	paths := cli.FindParentPaths(".jira.d/config.yml")
	// prepend
	paths = append([]string{"/etc/jira-cli.yml"}, paths...)

	for _, file := range(paths) {
		parseYaml(file, opts)
	}
}

