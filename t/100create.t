#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira=../jira

PLAN 24

RUNS $jira create --project BASIC -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira view $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
priority: Medium
description: |
  description
EOF

RUNS $jira create --project SCRUM -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira view $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: To Do
summary: summary
project: SCRUM
issuetype: Bug
assignee: gojira
reporter: gojira
priority: Medium
description: |
  description
EOF

RUNS $jira create --project KANBAN -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira view $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Backlog
summary: summary
project: KANBAN
issuetype: Bug
assignee: gojira
reporter: gojira
priority: Medium
description: |
  description
EOF

RUNS $jira create --project PROJECT -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira view $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: To Do
summary: summary
project: PROJECT
issuetype: Task
assignee: gojira
reporter: gojira
priority: Medium
description: |
  description
EOF

RUNS $jira create --project PROCESS -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira view $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
priority: Medium
description: |
  description
EOF

RUNS $jira create --project TASK -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira view $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: To Do
summary: summary
project: TASK
issuetype: Task
assignee: gojira
reporter: gojira
priority: Medium
description: |
  description
EOF
