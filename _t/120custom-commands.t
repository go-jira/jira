#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
. env.sh

PLAN 16

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
## Testing the example custom commands, print-project
###############################################################################

RUNS $jira print-project
DIFF <<EOF
BASIC
EOF

###############################################################################
## Testing the example custom commands, jira-path
###############################################################################

RUNS $jira jira-path
DIFF <<EOF
../jira
EOF

###############################################################################
## Testing the example custom commands, env
###############################################################################

RUNS $jira env
GREP ^JIRA_PROJECT=BASIC

###############################################################################
## Testing the example custom commands, argtest
###############################################################################

RUNS $jira argtest TEST
DIFF <<EOF
TEST
EOF

###############################################################################
## Testing the example custom commands, opttest
###############################################################################

RUNS $jira opttest --OPT TEST
DIFF <<EOF
TEST
EOF

###############################################################################
## Use the "mine" alias to list issues assigned to self
###############################################################################

RUNS $jira mine
DIFF <<EOF
+------------+------------------------------------------+------------+----------------+----------------+------------+---------------------------+--------------+
| Issue      | Summary                                  | Type       | Priority       | Status         | Age        | Reporter                  | Assignee     |
+------------+------------------------------------------+------------+----------------+----------------+------------+---------------------------+--------------+
| $(printf %-10s $issue) | summary                                  | Bug        | Medium         | To Do          | a minute   | gojira                    | gojira       |
+------------+------------------------------------------+------------+----------------+----------------+------------+---------------------------+--------------+
EOF
