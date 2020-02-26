#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
. env.sh

PLAN 6

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

# just get the number
shortIssue=${issue#BASIC-}

###############################################################################
## View the issue we just created
###############################################################################

RUNS $jira view $shortIssue
DIFF <<EOF
issue: $issue
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: GoJira
reporter: GoJira
priority: Medium
votes: 0
description: |
  description
EOF

