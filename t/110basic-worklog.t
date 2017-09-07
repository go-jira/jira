#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
. env.sh

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
RUNS $jira worklog add $issue --comment "work is hard" --time-spent "1h 12m" -S "2017-01-29T09:17:00.000-0500" --noedit
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
  started: 2017-01-29T06:17:00.000-0800
  timeSpent: 1h 12m

EOF
