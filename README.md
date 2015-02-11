# go-jira
simple jira command line client in Go

## Build

```bash
git clone git@github.com:Netflix-Skunkworks/go-jira.git
cd go-jira
export GOPATH=$(pwd)
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
  jira [-v ...] [-u USER] [-e URI] [-t FILE] ls [--query=JQL]
  jira [-v ...] [-u USER] [-e URI] [-t FILE] view ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] ISSUE


General Options:
  -h --help           Show this usage
  --version           Show this version
  -v --verbose        Increase output logging
  -u --user=USER      Username to use for authenticaion (default: cbennett)
  -e --endpoint=URI   URI to use for jira (default: https://jira)
  -t --template=FILE  Template file to use for output

List options:
  -q --query=JQL      Jira Query Language expression for the search
```