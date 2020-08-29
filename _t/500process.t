#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
. env.sh

PLAN 84

# cleanup from previous failed test executions
($jira ls --project PROCESS | awk -F: '{print $1}' | while read issue; do ../jira start $issue; done) | sed 's/^/# CLEANUP: /g'
($jira ls --project PROCESS | awk -F: '{print $1}' | while read issue; do ../jira stop $issue; done) | sed 's/^/# CLEANUP: /g'

# for any issues still remaining, they are stuck in "Under Review" status
($jira ls --project PROCESS | awk -F: '{print $1}' | while read issue; do ../jira transition --noedit -m  "approve" "Approve" $issue; done) | sed 's/^/# CLEANUP: /g'
($jira ls --project PROCESS | awk -F: '{print $1}' | while read issue; do ../jira transition --noedit -m  "done" "Done" $issue; done) | sed 's/^/# CLEANUP: /g'

# reset login
RUNS $jira logout
echo "gojira123" | RUNS $jira login

###############################################################################
## Create an issue
###############################################################################
RUNS $jira create --project PROCESS -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)

DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
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
assignee: GoJira
reporter: GoJira
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## List all issues, should be just the one we created
###############################################################################

RUNS $jira ls --project PROCESS
DIFF <<EOF
$(printf %-12s $issue:) summary
EOF

###############################################################################
## Try to close the issue, but PROCESS projects do not allow that state
###############################################################################

NRUNS $jira close $issue
EDIFF <<EOF
ERROR Invalid Transition "close" from "Open", Available: Start Progress
EOF

###############################################################################
## put the issue into Start Progress state, then Stop Progress state
## which will resolve the issue
###############################################################################

RUNS $jira start $issue
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF


RUNS $jira stop $issue
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

###############################################################################
## Verify there are no unresolved issues
###############################################################################

RUNS $jira ls --project PROCESS
DIFF <<EOF
EOF

###############################################################################
## Setup 2 more issues so we can test duping
###############################################################################

RUNS $jira create --project PROCESS -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

RUNS $jira create --project PROCESS -o summary=dup -o description=dup --noedit --saveFile issue.props
dup=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $dup $ENDPOINT/browse/$dup
EOF

###############################################################################
## Mark issue as duplicate, expect both issues to be updated and when viewing
## the main issue there should be a "depends" line showing the dup'd issue, and
## that issue should be resolved. For PROCESSS projects it has to go through
## 2 steps to resolve, one is "Start Progress" then resolved with "Stop
## Progress", so we see 3 updates in total
###############################################################################
RUNS $jira dup $dup $issue
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
OK $dup $ENDPOINT/browse/$dup
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: GoJira
reporter: GoJira
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

RUNS $jira ls --project PROCESS
DIFF <<EOF
$(printf %-12s $issue:) summary
EOF

###############################################################################
## Setup for testing blocking issues
###############################################################################

RUNS $jira create --project PROCESS -o summary=blocks -o description=blocks --noedit --saveFile issue.props
blocker=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

###############################################################################
## Set blocker and verify it shows up when viewing the main issue
###############################################################################

RUNS $jira block $blocker $issue
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
OK $blocker $ENDPOINT/browse/$blocker
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: GoJira
reporter: GoJira
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

RUNS $jira ls --project PROCESS
DIFF <<EOF
$(printf %-12s $issue:) summary
$(printf %-12s $blocker:) blocks
EOF

###############################################################################
# reset login for mothra for voting
###############################################################################

jira="$jira --user mothra --login mothra@corybennett.org"

RUNS $jira logout
echo "mothra123" | RUNS $jira login

###############################################################################
## vote for main issue, verify it shows when viewing the issue
###############################################################################

RUNS $jira vote $issue
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: GoJira
reporter: GoJira
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
OK $issue $ENDPOINT/browse/$issue
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: GoJira
reporter: GoJira
blockers: $blocker[Open]
depends: $dup[Cancelled]
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## set mothra user as watcher to issue and verify from REST api
###############################################################################

RUNS $jira watch $issue
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

# FIXME we probably need a watchers command to wrap this?
RUNS sh -c "$jira req /rest/api/2/issue/$issue/watchers | jq -r .watchers[].displayName | sort"
DIFF <<EOF
GoJira
Mothra
EOF

###############################################################################
## set issue to In Progress state, which is an invalid state for PROCESS
###############################################################################

NRUNS $jira trans "In Progress" $blocker --noedit
DIFF <<EOF
ERROR Invalid Transition "In Progress" from "Open", Available: Start Progress
EOF


###############################################################################
## Set issue to "In Review" state, which is an invalid state for PROCESS
###############################################################################

NRUNS $jira trans "review" $blocker --noedit
DIFF <<EOF
ERROR Invalid Transition "review" from "Open", Available: Start Progress
EOF

###############################################################################
## Set it to "Start Progress"
###############################################################################

RUNS $jira start $blocker
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

###############################################################################
## Set it back to "Stop Progress"
###############################################################################

RUNS $jira stop $blocker
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF


###############################################################################
## Set it to "Done"
###############################################################################

RUNS $jira reopen $blocker
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
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
assignee: GoJira
reporter: GoJira
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
OK $issue $ENDPOINT/browse/$issue
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: Open
summary: summary
project: PROCESS
issuetype: Task
assignee: GoJira
reporter: GoJira
blockers: $blocker[Open]
depends: $dup[Cancelled]
priority: Medium
votes: 0
description: |
  description

comments:
  - | # Mothra, a minute ago
    Yo, Comment

EOF

###############################################################################
## Verify we can add labels to an issue
###############################################################################

RUNS $jira labels add $blocker test-label another-label
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: GoJira
reporter: GoJira
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

RUNS $jira labels remove $blocker another-label
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: GoJira
reporter: GoJira
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

RUNS $jira labels set $blocker more-label better-label
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: GoJira
reporter: GoJira
blockers: 
depends: $issue[Open]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
EOF

###############################################################################
## Verify that "mothra" user can take the issue (reassign to self)
###############################################################################

RUNS $jira take $blocker
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: Mothra
reporter: GoJira
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
OK $blocker $ENDPOINT/browse/$blocker
EOF

RUNS $jira $blocker
DIFF <<EOF
issue: $blocker
created: a minute ago
status: Open
summary: blocks
project: PROCESS
issuetype: Task
assignee: GoJira
reporter: GoJira
blockers: 
depends: $issue[Open]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
EOF

