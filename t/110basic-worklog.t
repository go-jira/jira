#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira --project BASIC"
export JIRA_LOG_FORMAT="%{level:-5s} %{message}"

PLAN 8

# reset login
RUNS $jira logout
echo "gojira123" | RUNS $jira login

# cleanup from previous failed test executions
($jira ls | awk -F: '{print $1}' | while read issue; do ../jira done $issue; done) | sed 's/^/# CLEANUP: /g'

###############################################################################
## Create an issue
###############################################################################
RUNS $jira create -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)

DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

###############################################################################
## Add a worklog to an issue
###############################################################################
RUNS $jira add worklog $issue --comment "work is hard" --time-spent "1h 12m" --noedit
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
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
