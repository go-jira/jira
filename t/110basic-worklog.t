#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
export JIRA_LOG_FORMAT="%{level:-5s} %{message}"

ENDPOINT="http://localhost:8080"
if [ -n "$JIRACLOUD" ]; then
    ENDPOINT="https://go-jira.atlassian.net"
fi

PLAN 8

# reset login
RUNS $jira logout
RUNS $jira login

# cleanup from previous failed test executions
($jira ls --project BASIC | awk -F: '{print $1}' | while read issue; do ../jira done $issue; done) | sed 's/^/# CLEANUP: /g'

###############################################################################
## Create an issue
###############################################################################
RUNS $jira create --project BASIC -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)

DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

###############################################################################
## Add a worklog to an issue
###############################################################################
RUNS $jira worklog add $issue --comment "work is hard" --time-spent "1h 12m" --noedit
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

###############################################################################
## Verify worklog got added to issue
###############################################################################
RUNS $jira worklog $issue
DIFF <<EOF
- # gojira, a minute ago
  comment: work is hard
  timeSpent: 1h 12m

EOF
