package jiracli

import (
	"os"
	"strconv"

	logging "gopkg.in/op/go-logging.v1"
)

var (
	log = logging.MustGetLogger("jira")
)

func IncreaseLogLevel(verbosity int) {
	logging.SetLevel(logging.GetLevel("")+logging.Level(verbosity), "")
}

func InitLogging() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	format := os.Getenv("JIRA_LOG_FORMAT")
	if format == "" {
		format = "%{color}%{level:-5s}%{color:reset} %{message}"
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
		if verbosity, err := strconv.Atoi(os.Getenv("JIRA_DEBUG")); err == nil {
			IncreaseLogLevel(verbosity)
		}
	}
}
