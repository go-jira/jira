# go-jira
simple jira command line client in Go

## Synopsis

```bash
jira ls -p GOJIRA                       # list all unresolved issues for project GOJRIA
jira ls -p GOJIRA -a mothra             # as above also assigned to user mothra
jira ls -p GOJIRA -w mothra             # lists GOJIRA unresolved issues watched by user mothra
jira ls -p GOJIRA -r mothra             # list GOJIRA unresolved issues reported by user mothra

jira view GOJIRA-321                    # print Issue using "view" template
jira GOJIRA-321                         # same as above

jira edit GOJIRA-321                    # open up the issue in an editor, when you exit the
                                        # editor the issue will post the updates to the server

# edit the issue, using the overirdes on the command line, skip the interactive editor:
jira edit GOJIRA-321 --noedit \
     -o assignee=mothra \
     -o comment="mothra, please take care of this." \
     -o priority=Major

jira create -p GOJIRA                   # create new "Bug" type issue for project GOJIRA
jira create -p GOJIRA -i Task           # create new Task type issue

jira trans close GOJIRA-321             # close issue, with interactive editor to set fields
jira close GOJIRA-321 --edit            # same as above

# close the issue, set the resolution, and skip interactive editor:
jira trans close GOJIRA-321 -o resolution="Won't Fix" --noedit
# same as above
jira close GOJIRA-321 -o resolution="Won't Fix"

jira repopen GOJIRA-321 -m "reopening"  # reopen issue

jira watch GOJIRA-321                   # add self as watcher to the issue

jira comment GOJIRA-321 -m "done yet?"  # add comment to the issue

jira take GOJIRA-321                    # assign issue to self

jira give GOJIRA-321 mothra             # assign issue to user mothra

# create local project config to set defaults
mkdir .jira.d
echo "project: GOJIRA" > .jira.d/config.yml

jira ls                                 # list all unresolved issues for project GOJRIA
jira ls -a mothra                       # as above also assigned to user mothra
jira ls -w mothra                       # lists GOJIRA unresolved issues watched by user mothra
jira ls -r mothra                       # list GOJIRA unresolved issues reported by user mothra

jira create                             # create new "Bug" type issue for project GOJIRA
jira create -i Task                     # create new Task type issue
```

## Download

You can download one of the pre-built binaries for **go-jira** [here](https://github.com/Netflix-Skunkworks/go-jira/releases).

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

**go-jira** uses a configuration hierarchy.  When loading the configuration from disk it will recursively look through
all parent directories in your current path looking for a **.jira.d** directory.  If your current directory is not
a child directory of your homedir, then your homedir will also be inspected for a **.jira.d** directory.  From all of **.jira.d** directories
discovered **go-jira** will load a **config.yml** if found.  The configuration properties found in a file closests to your current working directory
will have precedence.  Properties overriden with command line options will have final precedence.

The complicated configuration heirarchy is used because **go-jira** attempts to be context aware.  For example, if you are working on a "foo" project and
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

When you run command like `jira edit` it will open up your favorite editor with the templatized output so you can quickly edit.  When the editor
closes **go-jira** will submit the completed form.  The order which **go-jira** attempts to determine your prefered editor is:

* **editor** property in any config.yml file
* **JIRA_EDITOR** environment variable
* **EDITOR** environment variable
* vim

### Templates

**go-jira** has the ability to customize most output (and editor input) via templates.  There are default templates available for all operations,
which may or may not work for your actual jira implementation.  Jira is endlessly customizable, so it is hard to provide default templates
that will work for all issue types.

When running a command like `jira edit` it will look through the current directory hierarchy trying to find a file that matches **.jira.d/templates/edit**,
if found it will use that file as the template, otherwise it will use the default **edit** template hard-coded into **go-jira**.  You can export the default
hard-coded templates with `jira export-templates` which will write them to **~/.jira.d/templates/**.

## Usage

```
Usage:
Usage:
  jira [-v ...] [-u USER] [-e URI] [-t FILE] (ls|list) ( [-q JQL] | [-p PROJECT] [-c COMPONENT] [-a ASSIGNEE] [-i ISSUETYPE] [-w WATCHER] [-r REPORTER]) 
  jira [-v ...] [-u USER] [-e URI] [-b] [-t FILE] view ISSUE
  jira [-v ...] [-u USER] [-e URI] [-b] [-t FILE] edit ISSUE [--noedit] [-m COMMENT] [-o KEY=VAL]... 
  jira [-v ...] [-u USER] [-e URI] [-b] [-t FILE] create [--noedit] [-p PROJECT] [-i ISSUETYPE] [-o KEY=VAL]...
  jira [-v ...] [-u USER] [-e URI] [-b] DUPLICATE dups ISSUE
  jira [-v ...] [-u USER] [-e URI] [-b] BLOCKER blocks ISSUE
  jira [-v ...] [-u USER] [-e URI] [-b] watch ISSUE [-w WATCHER]
  jira [-v ...] [-u USER] [-e URI] [-b] [-t FILE] (trans|transition) TRANSITION ISSUE [-m COMMENT] [--noedit]
  jira [-v ...] [-u USER] [-e URI] [-b] ack ISSUE [-m COMMENT] [--edit]
  jira [-v ...] [-u USER] [-e URI] [-b] close ISSUE [-m COMMENT] [--edit]
  jira [-v ...] [-u USER] [-e URI] [-b] resolve ISSUE [-m COMMENT] [--edit]
  jira [-v ...] [-u USER] [-e URI] [-b] reopen ISSUE [-m COMMENT] [--edit]
  jira [-v ...] [-u USER] [-e URI] [-b] start ISSUE [-m COMMENT] [--edit]
  jira [-v ...] [-u USER] [-e URI] [-b] stop ISSUE [-m COMMENT] [--edit]
  jira [-v ...] [-u USER] [-e URI] [-b] [-t FILE] comment ISSUE [-m COMMENT]
  jira [-v ...] [-u USER] [-e URI] [-b] take ISSUE
  jira [-v ...] [-u USER] [-e URI] [-b] (assign|give) ISSUE ASSIGNEE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] fields
  jira [-v ...] [-u USER] [-e URI] [-t FILE] issuelinktypes
  jira [-v ...] [-u USER] [-e URI] [-b][-t FILE] transmeta ISSUE
  jira [-v ...] [-u USER] [-e URI] [-b] [-t FILE] editmeta ISSUE
  jira [-v ...] [-u USER] [-e URI] [-t FILE] issuetypes [-p PROJECT] 
  jira [-v ...] [-u USER] [-e URI] [-t FILE] createmeta [-p PROJECT] [-i ISSUETYPE] 
  jira [-v ...] [-u USER] [-e URI] [-b] [-t FILE] transitions ISSUE
  jira [-v ...] export-templates [-d DIR]
  jira [-v ...] [-u USER] [-e URI] [-t FILE] login
  jira [-v ...] [-u USER] [-e URI] [-b] [-t FILE] ISSUE
 
General Options:
  -e --endpoint=URI   URI to use for jira
  -h --help           Show this usage
  -t --template=FILE  Template file to use for output/editing
  -u --user=USER      Username to use for authenticaion (default: cbennett)
  -v --verbose        Increase output logging
  --version           Show this version

Command Options:
  -a --assignee=USER        Username assigned the issue
  -b --browse               Open your browser to the Jira issue
  -c --component=COMPONENT  Component to Search for
  -d --directory=DIR        Directory to export templates to (default: /Users/cbennett/.jira.d/templates)
  -f --queryfields          Fields that are used in "list" template: (default: summary)
  -i --issuetype=ISSUETYPE  Jira Issue Type (default: Bug)
  -m --comment=COMMENT      Comment message for transition
  -o --override=KEY:VAL     Set custom key/value pairs
  -p --project=PROJECT      Project to Search for
  -q --query=JQL            Jira Query Language expression for the search
  -r --reporter=USER        Reporter to search for
  -w --watcher=USER         Watcher to add to issue (default: cbennett)
                            or Watcher to search for
```
