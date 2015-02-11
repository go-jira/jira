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

func parseYaml(file string, opts map[string]string) {
	if fh, err := ioutil.ReadFile(file); err == nil {
		log.Debug("Found Config file: %s", file)
		yaml.Unmarshal(fh, &opts)
	}
}

func loadConfigs(opts map[string]string) {
	paths := cli.FindParentPaths(".jira")
	// prepend
	paths = append([]string{"/etc/jira-cli.yml"}, paths...)

	for _, file := range(paths) {
		parseYaml(file, opts)
	}
}

func main() {
	user := os.Getenv("USER")
	usage := fmt.Sprintf(`
Usage:
  jira [-v ...] [-u USER] [-e URI] [-t FILE] fields
  jira [-v ...] [-u USER] [-e URI] [-t FILE] ls [--query=JQL]
  jira [-v ...] [-u USER] [-e URI] [-t FILE] view ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] ISSUE


General Options:
  -h --help           Show this usage
  --version           Show this version
  -v --verbose        Increase output logging
  -u --user=USER      Username to use for authenticaion (default: %s)
  -e --endpoint=URI   URI to use for jira (default: https://jira)
  -t --template=FILE  Template file to use for output

List options:
  -q --query=FILE  Template to use for output
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

	for key,val := range args {
		if val != nil && strings.HasPrefix(key, "--") {
			opt := key[2:]
			switch v := val.(type) {
			case string:
				opts[opt] = v
			}
		}
	}
	
	if _, ok := opts["endpoint"]; !ok {
		opts["endpoint"] = "https://jira"
	}
	if _, ok := opts["user"]; !ok {
		opts["user"] = user
	}
	
	c := cli.New(opts)

	log.Debug("opts: %s", opts);

	c.CmdLogin()

	if val, ok := args["fields"]; ok && val.(bool) {
		c.CmdFields()
	} else if val, ok := args["ls"]; ok && val.(bool) {
		c.CmdList()
	} else if val, ok := args["ISSUE"]; ok {
		c.CmdView(val.(string))
	}

	os.Exit(0)
}
