#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira=../jira
. env.sh

SKIP test -n "$JIRACLOUD" # using Jira Cloud at go-jira.atlassian.net

PLAN 7

###############################################################################
## Verify logout works, we expect when we call the session api
## that we will get a 401 and prompt user for password
################################################################################
RUNS $jira logout

NRUNS $jira req /rest/auth/1/session </dev/null
ODIFF <<EOF
Jira Password [gojira]: 
EOF

###############################################################################
## Verify login works (password read from stdin) and verify that the
## sesion api no longer prompts
###############################################################################
echo "gojira123" | RUNS $jira login

RUNS $jira req /rest/auth/1/session </dev/null
GREP '"name": "gojira"'
GREP "\"self\": \"$ENDPOINT/rest/api/latest/user?username=gojira\""


