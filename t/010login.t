#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira=../jira

PLAN 7

RUNS $jira logout

NRUNS $jira req /rest/auth/1/session </dev/null
ODIFF <<EOF
Jira Password [gojira]: 
EOF

echo "gojira123" | RUNS $jira login

RUNS $jira req /rest/auth/1/session </dev/null
GREP '"name": "gojira"'
GREP '"self": "http://localhost:8080/rest/api/latest/user?username=gojira"'


