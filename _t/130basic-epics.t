#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
. env.sh

PLAN 22

# reset login
RUNS $jira logout
RUNS $jira login

# cleanup from previous failed test executions
($jira ls --project BASIC | awk -F: '{print $1}' | while read issue; do ../jira done $issue; done) | sed 's/^/# CLEANUP: /g'

###############################################################################
## Create an epic
###############################################################################
RUNS $jira epic create --project BASIC -o summary="Totally Epic" -o description=description --epic-name "Basic Epic" --noedit --saveFile issue.props
epic=$(awk '/issue/{print $2}' issue.props)

DIFF <<EOF
OK $epic $ENDPOINT/browse/$epic
EOF

###############################################################################
## Create issues we can assign to epic
###############################################################################
RUNS $jira create --project BASIC -o summary="summary" -o description=description --noedit --saveFile issue.props
issue1=$(awk '/issue/{print $2}' issue.props)

DIFF <<EOF
OK $issue1 $ENDPOINT/browse/$issue1
EOF

RUNS $jira create --project BASIC -o summary="summary" -o description=description --noedit --saveFile issue.props
issue2=$(awk '/issue/{print $2}' issue.props)

DIFF <<EOF
OK $issue2 $ENDPOINT/browse/$issue2
EOF

###############################################################################
## List the issues for the epic
###############################################################################
RUNS $jira epic list $epic

DIFF<<EOF
+-------+---------+------+----------+--------+-----+----------+----------+
| Issue | Summary | Type | Priority | Status | Age | Reporter | Assignee |
+-------+---------+------+----------+--------+-----+----------+----------+
+-------+---------+------+----------+--------+-----+----------+----------+
EOF

###############################################################################
## Add issues to an epic
###############################################################################
RUNS $jira epic add $epic $issue1 $issue2

DIFF<<EOF
OK $epic $ENDPOINT/browse/$epic
OK $issue1 $ENDPOINT/browse/$issue1
OK $issue2 $ENDPOINT/browse/$issue2
EOF

###############################################################################
## List the issues for the epic
###############################################################################
RUNS $jira epic list $epic

DIFF<<EOF
+------------+---------+------+----------+--------+----------+----------+----------+
|   Issue    | Summary | Type | Priority | Status |   Age    | Reporter | Assignee |
+------------+---------+------+----------+--------+----------+----------+----------+
| $issue1 | summary | Bug  | Medium   | To Do  | a minute | GoJira   | GoJira   |
| $issue2 | summary | Bug  | Medium   | To Do  | a minute | GoJira   | GoJira   |
+------------+---------+------+----------+--------+----------+----------+----------+
EOF

###############################################################################
## Remove an issue from an Epic
###############################################################################
RUNS $jira epic remove $issue1

DIFF<<EOF
OK $issue1 $ENDPOINT/browse/$issue1
EOF

###############################################################################
## List the issues for the epic
###############################################################################
RUNS $jira epic list $epic

DIFF<<EOF
+------------+---------+------+----------+--------+----------+----------+----------+
|   Issue    | Summary | Type | Priority | Status |   Age    | Reporter | Assignee |
+------------+---------+------+----------+--------+----------+----------+----------+
| $issue2 | summary | Bug  | Medium   | To Do  | a minute | GoJira   | GoJira   |
+------------+---------+------+----------+--------+----------+----------+----------+
EOF

###############################################################################
## Remove last issue from an Epic
###############################################################################
RUNS $jira epic remove $issue2

DIFF<<EOF
OK $issue2 $ENDPOINT/browse/$issue2
EOF

###############################################################################
## List the issues for the epic
###############################################################################
RUNS $jira epic list $epic

DIFF<<EOF
+-------+---------+------+----------+--------+-----+----------+----------+
| Issue | Summary | Type | Priority | Status | Age | Reporter | Assignee |
+-------+---------+------+----------+--------+-----+----------+----------+
+-------+---------+------+----------+--------+-----+----------+----------+
EOF
