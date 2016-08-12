#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira --project PROCESS"
export JIRA_LOG_FORMAT="%{level:-5s} %{message}"

PLAN 84

# cleanup from previous failed test executions
($jira ls | awk -F: '{print $1}' | while read issue; do ../jira start $issue; done) | sed 's/^/# CLEANUP: /g'
($jira ls | awk -F: '{print $1}' | while read issue; do ../jira stop $issue; done) | sed 's/^/# CLEANUP: /g'

# reset login
RUNS $jira logout
echo "gojira123" | RUNS $jira login

###############################################################################
## Create an issue
###############################################################################
RUNS $jira create -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)

DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

###############################################################################
## View the issue we just created
###############################################################################

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
votes: 0
description: |
  description
EOF

###############################################################################
## List all issues, should be just the one we created
###############################################################################

RUNS $jira ls
DIFF <<EOF
$(printf %-12s $issue:) summary
EOF

###############################################################################
## Try to close the issue, but PROCESS projects do not allow that state
###############################################################################

NRUNS $jira close $issue
EDIFF <<EOF
ERROR Invalid Transition 'close', Available: Start Progress
EOF

###############################################################################
## put the issue into Start Progress state, then Stop Progress state
## which will resolve the issue
###############################################################################

RUNS $jira start $issue
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF


RUNS $jira stop $issue
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

###############################################################################
## Verify there are no unresolved issues
###############################################################################

RUNS $jira ls
DIFF <<EOF
EOF

###############################################################################
## Setup 2 more issues so we can test duping
###############################################################################

RUNS $jira create -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira create -o summary=dup -o description=dup --noedit --saveFile issue.props
dup=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $dup http://localhost:8080/browse/$dup
EOF

###############################################################################
## Mark issue as duplicate, expect both issues to be updated and when viewing
## the main issue there should be a "depends" line showing the dup'd issue, and
## that issue should be resolved. For PROCESSS projects it has to go through
## 2 steps to resolve, one is "Start Progress" then resolved with "Stop
## Progress", so we see 3 updates in total
###############################################################################

RUNS $jira $dup dups $issue --noedit
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
OK $dup http://localhost:8080/browse/$dup
OK $dup http://localhost:8080/browse/$dup
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: 
depends: $dup[Cancelled]
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## We should see only one unresolved issue, the Dup should be resolved
###############################################################################

RUNS $jira ls
DIFF <<EOF
$(printf %-12s $issue:) summary
EOF

###############################################################################
## Setup for testing blocking issues
###############################################################################

RUNS $jira create -o summary=blocks -o description=blocks --noedit --saveFile issue.props
blocker=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $blocker http://localhost:8080/browse/$blocker
EOF

###############################################################################
## Set blocker and verify it shows up when viewing the main issue
###############################################################################

RUNS $jira $blocker blocks $issue
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: $blocker[Open]
depends: $dup[Cancelled]
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## Both issues are unresolved now
###############################################################################

RUNS $jira ls
DIFF <<EOF
$(printf %-12s $issue:) summary
$(printf %-12s $blocker:) blocks
EOF

###############################################################################
# reset login for mojira for voting
###############################################################################

jira="$jira --user mojira"

RUNS $jira logout
echo "mojira123" | RUNS $jira login

###############################################################################
## vote for main issue, verify it shows when viewing the issue
###############################################################################

RUNS $jira vote $issue
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: $blocker[Open]
depends: $dup[Cancelled]
priority: Medium
votes: 1
description: |
  description
EOF

###############################################################################
## downvote the main issue, verify the vote count goes back to 0
###############################################################################

RUNS $jira vote $issue --down
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: $blocker[Open]
depends: $dup[Cancelled]
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## set mojira user as watcher to issue and verify from REST api
###############################################################################

RUNS $jira watch $issue
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

# FIXME we probably need a watchers command to wrap this?
RUNS sh -c "$jira req /rest/api/2/issue/$issue/watchers | jq -r .watchers[].name"
DIFF <<EOF
gojira
mojira
EOF

###############################################################################
## set issue to In Progress state, which is an invalid state for PROCESS
###############################################################################

NRUNS $jira trans "In Progress" $blocker --noedit
DIFF <<EOF
ERROR Invalid Transition 'In Progress', Available: Start Progress
EOF


###############################################################################
## Set issue to "In Review" state, which is an invalid state for PROCESS
###############################################################################

NRUNS $jira trans "review" $blocker --noedit
DIFF <<EOF
ERROR Invalid Transition 'review', Available: Start Progress
EOF

###############################################################################
## Set it to "Start Progress"
###############################################################################

RUNS $jira start $blocker
DIFF <<EOF
OK $blocker http://localhost:8080/browse/$blocker
EOF

###############################################################################
## Set it back to "Stop Progress"
###############################################################################

RUNS $jira stop $blocker
DIFF <<EOF
OK $blocker http://localhost:8080/browse/$blocker
EOF


###############################################################################
## Set it to "Done"
###############################################################################

RUNS $jira reopen $blocker
DIFF <<EOF
OK $blocker http://localhost:8080/browse/$blocker
EOF

###############################################################################
## Verify issue is now in Done state (the "blocker" issue is now Done)
###############################################################################

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: $blocker[Open]
depends: $dup[Cancelled]
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## Verify we can add a comment
###############################################################################

RUNS $jira comment $issue --noedit -m "Yo, Comment"
DIFF <<EOF
OK $issue http://localhost:8080/browse/$issue
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: $blocker[Open]
depends: $dup[Cancelled]
priority: Medium
votes: 0
description: |
  description

comments:
  - | # mojira, a minute ago
    Yo, Comment

EOF

###############################################################################
## Verify we can add labels to an issue
###############################################################################

RUNS $jira add labels $blocker test-label another-label
DIFF <<EOF
OK $blocker http://localhost:8080/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: 
depends: $issue[Open]
priority: Medium
votes: 0
labels: another-label, test-label
description: |
  blocks
EOF

###############################################################################
## Verify we can remove a label
###############################################################################

RUNS $jira remove labels $blocker another-label
DIFF <<EOF
OK $blocker http://localhost:8080/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: 
depends: $issue[Open]
priority: Medium
votes: 0
labels: test-label
description: |
  blocks
EOF

###############################################################################
## Verify we can replace the labels with a new set
###############################################################################

RUNS $jira set labels $blocker more-label better-label
DIFF <<EOF
OK $blocker http://localhost:8080/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: 
depends: $issue[Open]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
EOF

###############################################################################
## Verify that "mojira" user can take the issue (reassign to self)
###############################################################################

RUNS $jira take $blocker
DIFF <<EOF
OK $blocker http://localhost:8080/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: mojira
reporter: gojira
blockers: 
depends: $issue[Open]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
EOF

###############################################################################
## Verify we can give the issue back go "gojira" user
###############################################################################

RUNS $jira give $blocker gojira
DIFF <<EOF
OK $blocker http://localhost:8080/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: gojira
reporter: gojira
blockers: 
depends: $issue[Open]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
EOF

