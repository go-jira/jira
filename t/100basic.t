#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
. env.sh

PLAN 94

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
## View the issue we just created
###############################################################################

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
votes: 0
description: |
  description
EOF

###############################################################################
## List all issues, should be just the one we created
###############################################################################

RUNS $jira ls --project BASIC
DIFF <<EOF
$(printf %-12s $issue:) summary
EOF

###############################################################################
## List issues using a named query
###############################################################################
RUNS $jira ls --project BASIC -n todo
DIFF <<EOF
$(printf %-12s $issue:) summary
EOF

###############################################################################
## List all issues, using the table template
###############################################################################

RUNS $jira ls --project BASIC --template table
DIFF <<EOF
+----------------+------------------------------------------+--------------+--------------+--------------+------------+--------------+--------------+
| Issue          | Summary                                  | Type         | Priority     | Status       | Age        | Reporter     | Assignee     |
+----------------+------------------------------------------+--------------+--------------+--------------+------------+--------------+--------------+
| $(printf %-14s $issue) | summary                                  | Bug          | Medium       | To Do        | a minute   | gojira       | gojira       |
+----------------+------------------------------------------+--------------+--------------+--------------+------------+--------------+--------------+
EOF

###############################################################################
## Edit an issue
###############################################################################
RUNS $jira edit $issue -m "edit comment" --override priority=High --noedit
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

###############################################################################
## Edit multiple issues with query
###############################################################################
RUNS $jira edit -m "bulk edit comment" --override priority=High --noedit --query "resolution = unresolved AND project = 'BASIC' AND status = 'To Do'"
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

RUNS $jira $issue
DIFF <<EOF
issue: $issue
created: a minute ago
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
priority: High
votes: 0
description: |
  description

comments:
  - | # gojira, a minute ago
    edit comment
  - | # gojira, a minute ago
    bulk edit comment

EOF

###############################################################################
## Try to close the issue, bug Basic projects do not allow that state
###############################################################################

NRUNS $jira close $issue
EDIFF <<EOF
ERROR Invalid Transition "close" from "To Do", Available: To Do, In Progress, In Review, Done
EOF

###############################################################################
## put the issue into Done state, resolving it.
###############################################################################

RUNS $jira done $issue
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

###############################################################################
## Verify there are no unresolved issues
###############################################################################

RUNS $jira ls --project BASIC
DIFF <<EOF
EOF

###############################################################################
## Setup 2 more issues so we can test duping
###############################################################################

RUNS $jira create --project BASIC -o summary=summary -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)
DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

RUNS $jira create --project BASIC -o summary=dup -o description=dup --noedit --saveFile issue.props
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
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
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

RUNS $jira ls --project BASIC
DIFF <<EOF
$(printf %-12s $issue:) summary
EOF

###############################################################################
## Setup for testing blocking issues
###############################################################################

RUNS $jira create --project BASIC -o summary=blocks -o description=blocks --noedit --saveFile issue.props
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
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
blockers: $blocker[To Do]
depends: $dup[Done]
priority: Medium
votes: 0
description: |
  description
EOF

###############################################################################
## Both issues are unresolved now
###############################################################################

RUNS $jira ls --project BASIC
DIFF <<EOF
$(printf %-12s $issue:) summary
$(printf %-12s $blocker:) blocks
EOF

###############################################################################
# reset login for mothra for voting
###############################################################################

jira="$jira --user mothra --login mothra@corybennett.org"

RUNS $jira logout
RUNS $jira login

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
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
blockers: $blocker[To Do]
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
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
blockers: $blocker[To Do]
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
RUNS sh -c "$jira req /rest/api/2/issue/$issue/watchers | jq -r .watchers[].name | sort"
DIFF <<EOF
gojira
mothra
EOF

###############################################################################
## set issue to In Progress state
###############################################################################

RUNS $jira trans "In Progress" $blocker --noedit
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

###############################################################################
## set it back to "To Do"
###############################################################################

RUNS $jira todo $blocker
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

###############################################################################
## Set issue to "In Review" state
###############################################################################

RUNS $jira trans "review" $blocker --noedit
DIFF <<EOF
OK $blocker $ENDPOINT/browse/$blocker
EOF

###############################################################################
## Set it back to "To Do"
###############################################################################

RUNS $jira todo $blocker
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
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
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
status: To Do
summary: summary
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
blockers: $blocker[Done]
depends: $dup[Done]
priority: Medium
votes: 0
description: |
  description

comments:
  - | # mothra, a minute ago
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
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
blockers: 
depends: $issue[To Do]
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
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
blockers: 
depends: $issue[To Do]
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
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
blockers: 
depends: $issue[To Do]
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
project: BASIC
issuetype: Bug
assignee: mothra
reporter: gojira
blockers: 
depends: $issue[To Do]
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
project: BASIC
issuetype: Bug
assignee: gojira
reporter: gojira
blockers: 
depends: $issue[To Do]
priority: Medium
votes: 0
labels: better-label, more-label
description: |
  blocks
EOF


###############################################################################
## List 150 closed issues, should be more than 100
###############################################################################

RUNS $jira ls --project BASIC --status Closed --limit 150
IS $(printf $0 | wc -l) -eq 150
