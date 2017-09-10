#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
. env.sh

PLAN 10

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
## Testing the example custom commands, env
###############################################################################

RUNS $jira env
DIFF <<'EOF'
JIRACLOUD=1
JIRA_CUSTOM_COMMANDS=[{"name":"env","script":"env | sort | grep JIRA","help":"print the JIRA environment variables available to custom commands"},{"name":"print-project","script":"echo $JIRA_PROJECT","help":"print the name of the configured project"},{"name":"mine","script":"if [ -n \"$JIRA_PROJECT\" ]; then\n    # if `project: ...` configured just list the issues for current project\n    jira list --template table --query \"resolution = unresolved and assignee=currentuser() and project = $JIRA_PROJECT ORDER BY priority asc, created\"\nelse\n    # otherwise list issues for all project\n    jira list --template table --query \"resolution = unresolved and assignee=currentuser() ORDER BY priority asc, created\"\nfi","help":"display issues assigned to me"}]
JIRA_ENDPOINT=https://go-jira.atlassian.net
JIRA_LOG_FORMAT=%{level:-5s} %{message}
JIRA_PASSWORD_SOURCE=pass
JIRA_PROJECT=BASIC
JIRA_USER=gojira
EOF

###############################################################################
## Use the "mine" alias to list issues assigned to self
###############################################################################

RUNS $jira mine
DIFF <<EOF
+----------------+---------------------------------------------------------+--------------+--------------+------------+--------------+--------------+
| Issue          | Summary                                                 | Priority     | Status       | Age        | Reporter     | Assignee     |
+----------------+---------------------------------------------------------+--------------+--------------+------------+--------------+--------------+
| $(printf %-14s $issue) | summary                                                 | Medium       | To Do        | a minute   | gojira       | gojira       |
+----------------+---------------------------------------------------------+--------------+--------------+------------+--------------+--------------+
EOF
