# go-jira
simple jira command line client in Go

## Build

```bash
git clone git@github.com:Netflix-Skunkworks/go-jira.git
cd go-jira
export GOPATH=$(pwd)
export GOBIN=$GOPATH/bin
cd src/github.com/Netflix-Skunkworks/go-jira/jira
go get -v
go install -v
```

## Simple Config file

```bash
mkdir ~/.jira.d

cat <<EOM >~/.jira.d/config.yml
endpoint: https://jira.mycompany.com
EOM
```

## Usage

```
Usage:
  jira [-v ...] [-u USER] [-e URI] [-t FILE] fields
  jira [-v ...] [-u USER] [-e URI] [-t FILE] login
  jira [-v ...] [-u USER] [-e URI] [-t FILE] (ls|list) ( [-q JQL] | [-p PROJECT] [-c COMPONENT] [-a ASSIGNEE] [-i ISSUETYPE]) 
  jira [-v ...] [-u USER] [-e URI] [-t FILE] view ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] issuelinktypes
  jira [-v ...] [-u USER] [-e URI] [-t FILE] transmeta ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] editmeta ISSUE
  jira [-v ...] export-templates [-d DIR]
  jira [-v ...] [-u USER] [-e URI] [-t FILE] ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] edit ISSUE [-m COMMENT] [-o KEY=VAL]...
  jira [-v ...] [-u USER] [-e URI] [-t FILE] issuetypes [-p PROJECT] 
  jira [-v ...] [-u USER] [-e URI] [-t FILE] createmeta [-p PROJECT] [-i ISSUETYPE] 
  jira [-v ...] [-u USER] [-e URI] [-t FILE] transitions ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] create [-p PROJECT] [-i ISSUETYPE] [-o KEY=VAL]...
  jira [-v ...] [-u USER] [-e URI] DUPLICATE dups ISSUE
  jira [-v ...] [-u USER] [-e URI] BLOCKER blocks ISSUE
  jira [-v ...] [-u USER] [-e URI] watch ISSUE [-w WATCHER]
  jira [-v ...] [-u USER] [-e URI] (trans|transition) TRANSITION ISSUE [-m COMMENT]
  jira [-v ...] [-u USER] [-e URI] ack ISSUE [-m COMMENT]
  jira [-v ...] [-u USER] [-e URI] close ISSUE [-m COMMENT]
  jira [-v ...] [-u USER] [-e URI] resolve ISSUE [-m COMMENT]
  jira [-v ...] [-u USER] [-e URI] reopen ISSUE [-m COMMENT]
  jira [-v ...] [-u USER] [-e URI] start ISSUE [-m COMMENT]
  jira [-v ...] [-u USER] [-e URI] stop ISSUE [-m COMMENT]
  jira [-v ...] [-u USER] [-e URI] [-t FILE] comment ISSUE [-m COMMENT]
  jira [-v ...] [-u USER] [-e URI] take ISSUE
  jira [-v ...] [-u USER] [-e URI] (assign|give) ISSUE ASSIGNEE

General Options:
  -e --endpoint=URI   URI to use for jira (default: https://jira)
  -h --help           Show this usage
  -t --template=FILE  Template file to use for output/editing
  -u --user=USER      Username to use for authenticaion (default: cbennett)
  -v --verbose        Increase output logging
  --version           Show this version

Command Options:
  -a --assignee=USER        Username assigned the issue
  -c --component=COMPONENT  Component to Search for
  -d --directory=DIR        Directory to export templates to (default: /Users/cbennett/.jira.d/templates)
  -i --issuetype=ISSUETYPE  Jira Issue Type (default: Bug)
  -m --comment=COMMENT      Comment message for transition
  -o --override=KEY:VAL     Set custom key/value pairs
  -p --project=PROJECT      Project to Search for
  -q --query=JQL            Jira Query Language expression for the search
  -w --watcher=USER         Watcher to add to issue (default: cbennett)
```