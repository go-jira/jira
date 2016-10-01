package jiracli

import (
	"os"

	"gopkg.in/Netflix-Skunkworks/go-jira.v1/lib"
	"gopkg.in/alecthomas/kingpin.v2"
)

func usage() {
	app := kingpin.New("jira", "Jira Command Line Client")
	app.Writer(os.Stdout)
	app.Version(jira.VERSION)

	opts := &JiraOptions{}

	// -b --browse         Open your browser to the Jira issue
	// -e --endpoint=URI   URI to use for jira
	// -k --insecure       disable TLS certificate verification
	// -h --help           Show this usage
	// -t --template=FILE  Template file to use for output/editing
	// -u --user=USER      Username to use for authenticaion (default: %s)
	// -v --verbose        Increase output logging

	listUsage(app, opts)
}

func listUsage(app *kingpin.Application, opts *JiraOptions) {
	list := app.Command("list", "List Jira Issues")
	list.Flag("assignee", "Username assigned the issue").Short('a').StringVar(&opts.Assignee)
	list.Flag("component", "Component to use for query").Short('c').StringVar(&opts.Component)
	list.Flag("queryfields", "Fields that are used for \"list\" template").Short('f').StringVar(&opts.QueryFields)
	list.Flag("limit", "Maximum number of results to return in query").Short('l').IntVar(&opts.MaxResults)
	list.Flag("project", "Project to use for query").Short('p').StringVar(&opts.Project)
	list.Flag("query", "Jira Query Language expression for the search").Short('q').StringVar(&opts.Query)
	list.Flag("reporter", "Reporter to use in query").Short('r').StringVar(&opts.Reporter)
	list.Flag("sort", "Sort order used in query").Short('s').StringVar(&opts.Sort)
	list.Flag("watcher", "Watcher to use in query").Short('w').StringVar(&opts.Watcher)

	// -a --assignee=USER        Username assigned the issue
	// -c --component=COMPONENT  Component to Search for
	// -f --queryfields=FIELDS   Fields that are used in "list" template: (default: %s)
	// -i --issuetype=ISSUETYPE  The Issue Type
	// -l --limit=VAL            Maximum number of results to return in query (default: %d)
	// -p --project=PROJECT      Project to Search for
	// -q --query=JQL            Jira Query Language expression for the search
	// -r --reporter=USER        Reporter to search for
	// -s --sort=ORDER           For list operations, sort issues (default: %s)
	// -w --watcher=USER         Watcher to add to issue (default: %s)
	//                           or Watcher to search for

}
