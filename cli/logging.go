package jiracli

import (
	"os"

	"gopkg.in/op/go-logging.v1"
)

var (
	log           = logging.MustGetLogger("jiracli")
	defaultFormat = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
)

func InitLogging() {
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
	logging.SetLevel(logging.NOTICE, "")
}

func VerboseLogging() {
	logging.SetLevel(logging.GetLevel("")+1, "")
}
