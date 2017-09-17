#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira"
. env.sh

PLAN 43

# reset login
RUNS $jira logout
RUNS $jira login

# cleanup from previous failed test executions
($jira ls --project BASIC | awk -F: '{print $1}' | while read issue; do ../jira done $issue; done) | sed 's/^/# CLEANUP: /g'

###############################################################################
## Create an issue
###############################################################################
RUNS $jira create --project BASIC -o summary="Attach To Me" -o description=description --noedit --saveFile issue.props
issue=$(awk '/issue/{print $2}' issue.props)

DIFF <<EOF
OK $issue $ENDPOINT/browse/$issue
EOF

###############################################################################
## Attach via stdin
###############################################################################
RUNS $jira attach create $issue --filename README.md --saveFile attach.props < ./README.md
attach1=$(awk '/^id:/{print $2}' attach.props)

DIFF <<EOF
OK $attach1 $ENDPOINT/secure/attachment/$attach1/README.md
EOF

###############################################################################
## Attach binary file
###############################################################################
RUNS dd of=garbage.bin if=/dev/urandom count=1k bs=1k
RUNS $jira attach create $issue garbage.bin --saveFile attach.props
attach2=$(awk '/^id:/{print $2}' attach.props)

DIFF <<EOF
OK $attach2 $ENDPOINT/secure/attachment/$attach2/garbage.bin
EOF

###############################################################################
## Attach binary file with different name
###############################################################################
RUNS $jira attach create $issue garbage.bin --filename foobar.bin --saveFile attach.props
attach3=$(awk '/^id:/{print $2}' attach.props)

DIFF <<EOF
OK $attach3 $ENDPOINT/secure/attachment/$attach3/foobar.bin
EOF

###############################################################################
## List attachments
###############################################################################
RUNS $jira attach list $issue
DIFF <<EOF
+------------+------------------------------+------------+--------------+--------------+
| id         | filename                     | bytes      | user         | created      |
+------------+------------------------------+------------+--------------+--------------+
| $(printf %10s $attach1) | README.md                    |       1238 | gojira       | a minute     |
| $(printf %10s $attach2) | garbage.bin                  |    1048576 | gojira       | a minute     |
| $(printf %10s $attach3) | foobar.bin                   |    1048576 | gojira       | a minute     |
+------------+------------------------------+------------+--------------+--------------+
EOF

###############################################################################
## Fetch text attachment
###############################################################################
RUNS $jira attach get $attach1 -o attach1.txt
DIFF <<EOF
OK Wrote attach1.txt
EOF

# verify no diffs
RUNS diff -q README.md attach1.txt

###############################################################################
## Fetch text attachment to stdout
###############################################################################
RUNS sh -c "$jira attach get $attach1 -o- > attach1.txt"

# verify no diffs
RUNS diff -q README.md attach1.txt

###############################################################################
## Fetch text attachment as same name
###############################################################################
RUNS $jira attach get $attach1
DIFF <<EOF
OK Wrote README.md
EOF

# verify no diffs
RUNS git diff README.md

###############################################################################
## Fetch binary attachment
###############################################################################
RUNS $jira attach get $attach2 --output binary.out
DIFF <<EOF
OK Wrote binary.out
EOF

# verify no diffs
RUNS diff -q garbage.bin binary.out

###############################################################################
## Fetch binary attachment to stdout
###############################################################################
RUNS sh -c "$jira attach get $attach2 -o- > binary.out"

# verify no diffs
RUNS diff -q garbage.bin binary.out

###############################################################################
## Fetch binary attachment
###############################################################################
RUNS $jira attach get $attach3
DIFF <<EOF
OK Wrote foobar.bin
EOF

# verify no diffs
RUNS diff -q garbage.bin foobar.bin

###############################################################################
## Fetch binary attachment to stdout
###############################################################################
RUNS sh -c "$jira attach get $attach3 --output=- > binary.out"

# verify no diffs
RUNS diff -q garbage.bin binary.out

###############################################################################
## Delete attachment
###############################################################################
RUNS $jira attach remove $attach1
DIFF <<EOF
OK Deleted Attachment $attach1
EOF

RUNS $jira attach list $issue
DIFF <<EOF
+------------+------------------------------+------------+--------------+--------------+
| id         | filename                     | bytes      | user         | created      |
+------------+------------------------------+------------+--------------+--------------+
| $(printf %10s $attach2) | garbage.bin                  |    1048576 | gojira       | a minute     |
| $(printf %10s $attach3) | foobar.bin                   |    1048576 | gojira       | a minute     |
+------------+------------------------------+------------+--------------+--------------+
EOF


###############################################################################
## Delete attachment
###############################################################################
RUNS $jira attach rm $attach2
DIFF <<EOF
OK Deleted Attachment $attach2
EOF

RUNS $jira attach list $issue
DIFF <<EOF
+------------+------------------------------+------------+--------------+--------------+
| id         | filename                     | bytes      | user         | created      |
+------------+------------------------------+------------+--------------+--------------+
| $(printf %10s $attach3) | foobar.bin                   |    1048576 | gojira       | a minute     |
+------------+------------------------------+------------+--------------+--------------+
EOF

###############################################################################
## Delete last
###############################################################################
RUNS $jira attach rm $attach3
DIFF <<EOF
OK Deleted Attachment $attach3
EOF

RUNS $jira attach list $issue
DIFF <<EOF
+------------+------------------------------+------------+--------------+--------------+
| id         | filename                     | bytes      | user         | created      |
+------------+------------------------------+------------+--------------+--------------+
+------------+------------------------------+------------+--------------+--------------+
EOF
