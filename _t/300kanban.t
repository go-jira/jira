#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
. env.sh

PLAN 86

# cleanup from previous failed test executions
($jira ls --project KANBAN | awk -F: '{print $1}' | while read issue; do ../jira done $issue; done) | sed 's/^/# CLEANUP: /g'

# reset login
RUNS $jira logout
echo "gojira123" | RUNS $jira login

###############################################################################
## Create an issue
###############################################################################
RUNS $jira create --project KANBAN -o summary=summary -o description=description --noedit --saveFile issue.props
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
status: Backlog
summary: summary
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## List all issues, we should see ours
###############################################################################

RUNS $jira ls --project KANBAN | grep $issue

###############################################################################
## Try to close the issue, bug Basic projects do not allow that state
###############################################################################

NRUNS $jira close $issue
EDIFF <<EOF
ERROR Invalid Transition "close" from "Backlog", Available: Backlog, Selected for Development, In Progress, Done
EOF

###############################################################################
## put the issue into Done state, resolving it.
###############################################################################

RUNS $jira done $issue
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

###############################################################################
## Verify our issue is not unresolved
###############################################################################

RUNS $jira ls --project KANBAN | grep $issue

###############################################################################
## Setup 2 more issues so we can test duping
###############################################################################

RUNS $jira create --project KANBAN -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

RUNS $jira create --project KANBAN -o summary=dup -o description=dup --noedit --saveFile issue.props
dup=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $dup $ENDPOINT/browse/$dup
EOF

###############################################################################
## Mark issue as duplicate, expect both issues to be updated and when viewing
## the main issue there should be a "depends" line showing the dup'd issue, and
## that issue should be resolved
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
status: Backlog
summary: summary
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: $dup[Done]
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## We should see only one unresolved issue, the Dup should be resolved
###############################################################################

RUNS $jira ls --project KANBAN | grep $issue

###############################################################################
## Setup for testing blocking issues
###############################################################################

RUNS $jira create --project KANBAN -o summary=blocks -o description=blocks --noedit --saveFile issue.props
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
status: Backlog
summary: summary
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: $blocker[Backlog]
depends: $dup[Done]
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## Both issues are unresolved now
###############################################################################

RUNS $jira ls --project KANBAN | grep $issue
RUNS $jira ls --project KANBAN | grep $blocker

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
status: Backlog
summary: summary
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: $blocker[Backlog]
depends: $dup[Done]
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
status: Backlog
summary: summary
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: $blocker[Backlog]
depends: $dup[Done]
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
## set issue to In Progress state
###############################################################################

RUNS $jira trans "In Progress" $blocker --noedit
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

###############################################################################
## set it to "To Do", which is not a valid state for KANBAN issues
###############################################################################

NRUNS $jira todo $blocker
DIFF <<EOF
ERROR Invalid Transition "To Do" from "In Progress", Available: Backlog, Selected for Development, In Progress, Done
EOF

###############################################################################
## set issue back to backlog state
###############################################################################

RUNS $jira backlog $blocker --noedit
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

###############################################################################
## Set issue to "In Review" state, which is an invalid state for KANBAN
###############################################################################

NRUNS $jira trans "review" $blocker --noedit
DIFF <<EOF
ERROR Invalid Transition "review" from "Backlog", Available: Backlog, Selected for Development, In Progress, Done
EOF

###############################################################################
## Set it back to "Backlog"
###############################################################################

RUNS $jira backlog $blocker
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

###############################################################################
## Set it to "In Progress"
###############################################################################

RUNS $jira prog $blocker
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

###############################################################################
## Set it to "Done"
###############################################################################

RUNS $jira done $blocker
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
status: Backlog
summary: summary
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: $blocker[Done]
depends: $dup[Done]
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
status: Backlog
summary: summary
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: $blocker[Done]
depends: $dup[Done]
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
status: Done
summary: blocks
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: $issue[Backlog]
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
status: Done
summary: blocks
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: $issue[Backlog]
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
status: Done
summary: blocks
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: $issue[Backlog]
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
status: Done
summary: blocks
project: KANBAN
issuetype: Bug
assignee: Mothra
reporter: GoJira
blockers: 
depends: $issue[Backlog]
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
status: Done
summary: blocks
project: KANBAN
issuetype: Bug
assignee: GoJira
reporter: GoJira
blockers: 
depends: $issue[Backlog]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
EOF
