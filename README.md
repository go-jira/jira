# go-jira
simple jira command line client in Go

## Build

* **NOTE** You will need **`go-1.4.1`** minimum

*  If you do not have a **GOPATH** setup, these are simple build steps:

```bash
git clone git@github.com:Netflix-Skunkworks/go-jira.git
cd go-jira
export GOPATH=$(pwd)
export GOBIN=$GOPATH/bin
export PATH=$GOBIN:$PATH
cd src/github.com/Netflix-Skunkworks/go-jira/jira
go get -v
```

* If you do have a **GOPATH** setup, these are the standard steps to build:

```
cd $GOPATH
git clone git@github.com:Netflix-Skunkworks/go-jira.git src/github.com/Netflix-Skunkworks/go-jira
cd src/github.com/Netflix-Skunkworks/go-jira/jira
go get -v
```

## Configuration

**go-jira** uses a configuration heirarchy.  When loading the configuration from disk it will recursively look through
all parent directories in your current  path looking for a **.jira.d** directory anywhere in your.  (If your current directory is not
a child directory of your homedir, then your homedir will also be inspected for a **.jira.d** directory.  From all of **.jira.d** directories
discovered **go-jira** will load a **config.yml** found.  The configuration properties found in a file closests to your current working directory
will have precedence.  Properties overriden with command line options will have final precedence.

The complicated configuration heirarchy is used because **go-jira** attempts to be context aware.  For example if you are working on a "foo" project and 
you `cd` into your project workspace, wouldn't it be nice if `jira ls` automatically knew to list only issues related to the "foo" project?  Likewise when you
`cd` to the "bar" project then `jira ls` should only list issues related to "bar" project.  You can do this with by creating a configuration under your project
workspace at **./.jira.d/config.yml** that looks like:

```
project: foo
```

You will need to specify your local jira endpoint first, typically in your homedir like:

```bash
mkdir ~/.jira.d

cat <<EOM >~/.jira.d/config.yml
endpoint: https://jira.mycompany.com
EOM
```

### Editing

when you run command like `jira edit` it will open up your favorite editor with the templatized output, so you can quickly edit.  When then editor
closes then **go-jira** will submit the complted form.  To specify your editor, and this is the order which **go-jira** attempts to determine your p
prefered editor:

* **editor** property in any config.yml file
* **JIRA_EDITOR** envionment variable
* **EDITOR** environment variable
* vim

### Templates

**go-jira** has the ability to customize most output (and editor input) via templates  There are default templates available for all operations,
which may or may not work for your actual jira implementation.  Jira is endlessly customizable, so it is hard to provide default templates
that will work for all issue types.

When running a command like `jira edit` it will look through the current directory hierachy trying to find a file that matches **.jira.d/templates/edit**,
if found it will use that file as the template, otherwise it will use the default **edit** template hard-coded into **go-jira**.  You can export the default
hard-coded templates with `jira export-templates` which will write then to **~/.jira.d/templtes/.**.

## Usage

```
Usage:
  jira [-v ...] [-u USER] [-e URI] [-t FILE] (ls|list) ( [-q JQL] | [-p PROJECT] [-c COMPONENT] [-a ASSIGNEE] [-i ISSUETYPE]) 
  jira [-v ...] [-u USER] [-e URI] [-t FILE] view ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] edit ISSUE [-m COMMENT] [-o KEY=VAL]...
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
  jira [-v ...] [-u USER] [-e URI] [-t FILE] fields
  jira [-v ...] [-u USER] [-e URI] [-t FILE] issuelinktypes
  jira [-v ...] [-u USER] [-e URI] [-t FILE] transmeta ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] editmeta ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] issuetypes [-p PROJECT] 
  jira [-v ...] [-u USER] [-e URI] [-t FILE] createmeta [-p PROJECT] [-i ISSUETYPE] 
  jira [-v ...] [-u USER] [-e URI] [-t FILE] transitions ISSUE
  jira [-v ...] export-templates [-d DIR]
  jira [-v ...] [-u USER] [-e URI] [-t FILE] login
  jira [-v ...] [-u USER] [-e URI] [-t FILE] ISSUE
 
General Options:
  -e --endpoint=URI   URI to use for jira
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