[![Join the chat at https://gitter.im/go-jira-cli/help](https://badges.gitter.im/go-jira-cli/help.svg)](https://gitter.im/go-jira-cli/help?utm_source=badge&utm_medium=badge&utm_content=badge)
[![Build Status](https://travis-ci.org/Netflix-Skunkworks/go-jira.svg?branch=master)](https://travis-ci.org/Netflix-Skunkworks/go-jira)
[![GoDoc](https://godoc.org/gopkg.in/Netflix-Skunkworks/go-jira.v0?status.png)](https://godoc.org/gopkg.in/Netflix-Skunkworks/go-jira.v0)


# go-jira
simple command line client for Atlassian's Jira service written in Go

## Synopsis

```bash
jira ls -p GOJIRA                       # list all unresolved issues for project GOJRIA
jira ls -p GOJIRA -a mothra             # as above also assigned to user mothra
jira ls -p GOJIRA -w mothra             # lists GOJIRA unresolved issues watched by user mothra
jira ls -p GOJIRA -r mothra             # list GOJIRA unresolved issues reported by user mothra
jira ls -t table -p GOJIRA              # list all unresolved issues in pretty table output

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

jira reopen GOJIRA-321 -m "reopening"  # reopen issue

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
jira ls -t table                        # list all unresolved issues in pretty table output

jira create                             # create new "Bug" type issue for project GOJIRA
jira create -i Task                     # create new Task type issue

# make the table template your default "list" template:
jira export-templates -t table
mv $HOME/.jira.d/templates/table $HOME/.jira.d/templates/list
```

## Download

You can download one of the pre-built binaries for **go-jira** [here](https://github.com/Netflix-Skunkworks/go-jira/releases).

## Build

* **NOTE** You will need **`go-1.4.1`** minimum

*  To build the `jira` binary the current directory just run:
```bash
make
```
* To install the binary to you ~/bin directory you can run:
```bash
make install
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

Then use `jira login` to authenticate yourself as $USER. To change your username, use the `-u` CLI flag or set `user:` in your config.yml

### Dynamic Configuration

If the **.jira.d/config.yml** file is executable, then **go-jira** will attempt to execute the file and use the stdout for configuration.  You can use this to customize templates or other overrides depending on what type of operation you are running.  For example if you would like to use the "table" template when ever you run `jira ls`, then you can create a template like this:

```sh
#!/bin/sh

echo "endpoint: https://jira.mycompany.com"
echo "editor: emacs -nw"

case $JIRA_OPERATION in 
    list)
      echo "template: table";;
esac
```

Or if you always set the same overrides when you create an issue for your project you can do something like this:

```sh
#!/bin/sh
echo "project: GOJIRA"

case $JIRA_OPERATION in
    create)
        echo "assignee: $USER"
        echo "watchers: mothra"
        ;;
esac
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

#### Writing/Editing Templates

First the basic templating functionality is defined by the Go language 'text/template' library.  The library reference documentation can be found [here](https://golang.org/pkg/text/template/), and there is a good primer document [here](https://gohugo.io/templates/go-templates/).  `go-jira` also provides a few extra helper functions to make it a bit easlier to format the data, those functions are defined [here](https://github.com/Netflix-Skunkworks/go-jira/blob/master/util.go#L133).

Knowing what data and fields are available to any given template is not obvious. The easiest approach to determine what is available is to use the `debug` template on any given operation.  For eample to find out what is available to the "view" templates, you can use:
```
jira view GOJIRA-321 -t debug
```

This will print out the data in JSON format that is available to the template.  You can do this for any other operation, like "list":
```
jira list -t debug
```

Figuring out what is available to input templates (like for the `create` operation) is a bit more tricky, but similar.  To find the data available for a `create` template you can run:
```
jira create --dryrun -t debug --editor /bin/cat
```
This will attempt to fetch metadata for your default project (you can provide any options that you would normally specify for the `create` operation).  It uses the `--dryrun` option to prevent any actual updates being sent to Jira.  The `-t debug` is like before to cause the input to be serialized to JSON and printed for your inspection.  Finally the `--editor /bin/cat` will cause `go-jira` to just print the template rather than open up an editor and wait for you to edit/save it.

### Authentication

By default `go-jira` will prompt for a password automatically when we receive an 403 http response.  Then after authentication we cache the JSESSSION cookie returned by the service and reuse that on subsequent requests.  Typically this cookie will be valid for several hours (depending on the service configuration).  Many deployments of Jira (like the cloud services on atlassian.net) have "websudo" enabled which will prevent the cookie based authentcation from working.  On these deployments you have a few options with `go-jira`.  You can enable a `password-source` via `.jira.d/config.yml` with possible values of `keyring` or `pass`.

#### keyring password source
**Note: Version 0.1.9 required.**
On OSX and Linux there are a few keyring providers that `go-jira` can use (via this [golang module](https://github.com/tmc/keyring)).  To integrate `go-jira` with a supported keyring just add this configuration to `$HOME/.jira.d/config.yml`:
```yaml
password-source: keyring
```
After setting this and issuing a `jira login`, your credentials will be stored in your platform's backend (e.g. Keychain for Mac OS X) automatically. Subsequent operations, like a `jira ls`, should "just work" from there.

#### `pass` password source
**Note: Version 0.1.9 required.**
An alternative to the keyring password source is the `pass` tool (documentation [here](https://www.passwordstore.org/)).  This uses gpg to encrypt/decrypt passwords on demand and by using `gpg-agent` you can cache the gpg credentials for a period of time so you will not be prompted repeatedly for decrypting the passwords.  The advantage over the keyring integration is that `pass` can be used on more platforms than OSX and Linux, although it does require more setup.  To use `pass` for password storage and retrieval via `go-jira` just add this configuration to `$HOME/.jira.d/config.yml`:
```yaml
password-source: pass
```

This assumes you have already setup `pass` correctly on your system.  Specifically you will need to have created a gpg key like this:

```
$ gpg --gen-key
```

Then you will need the GPG Key ID you want associated with `pass`.  First list the available keys:
```
$ gpg --list-keys
/home/gojira/.gnupg/pubring.gpg
-------------------------------------------------
pub   2048R/A307D709 2016-12-18
uid                  Go Jira <gojira@example.com>
sub   2048R/F9A047B8 2016-12-18
```

Then initialize the `pass` tool to use the correct key:
```
$ pass init "Go Jira <gojira@example.com>"
```

You probably want to setup gpg-agent so that you dont have to type in your gpg passphrase all the time.  You can get `gpg-agent` to automatically start by adding something like this to your `$HOME/.bashrc`
```bash
if [ -f $HOME/.gpg-agent-info ]; then
    . $HOME/.gpg-agent-info
    export GPG_AGENT_INFO
fi
# verify sock file from GPG_AGENT_INFO is actually present
if [ ! -S "${GPG_AGENT_INFO%%:*}" ]; then
    # set passphrase cache so I only have to type my passphrase once a day
    eval $(gpg-agent --default-cache-ttl 604800 --daemon --write-env-file $HOME/.gpg-agent-info)
fi
export GPG_TTY=$(tty)
```


## Usage

```
Usage:
  jira (ls|list) <Query Options> 
  jira view ISSUE
  jira edit [--noedit] <Edit Options> [ISSUE | <Query Options>]
  jira create [--noedit] [-p PROJECT] <Create Options>
  jira DUPLICATE dups ISSUE
  jira BLOCKER blocks ISSUE
  jira watch ISSUE [-w WATCHER]
  jira (trans|transition) TRANSITION ISSUE [--noedit] <Edit Options>
  jira ack ISSUE [--edit] <Edit Options>
  jira close ISSUE [--edit] <Edit Options>
  jira resolve ISSUE [--edit] <Edit Options>
  jira reopen ISSUE [--edit] <Edit Options>
  jira start ISSUE [--edit] <Edit Options>
  jira stop ISSUE [--edit] <Edit Options>
  jira comment ISSUE [--noedit] <Edit Options>
  jira take ISSUE
  jira (assign|give) ISSUE ASSIGNEE
  jira fields
  jira issuelinktypes
  jira transmeta ISSUE
  jira editmeta ISSUE
  jira issuetypes [-p PROJECT] 
  jira createmeta [-p PROJECT] [-i ISSUETYPE] 
  jira transitions ISSUE
  jira export-templates [-d DIR] [-t template]
  jira (b|browse) ISSUE
  jira login
  jira ISSUE

General Options:
  -b --browse         Open your browser to the Jira issue
  -e --endpoint=URI   URI to use for jira
  -h --help           Show this usage
  -t --template=FILE  Template file to use for output/editing
  -u --user=USER      Username to use for authentication (default: $USER)
  -v --verbose        Increase output logging

Query Options:
  -a --assignee=USER        Username assigned the issue
  -c --component=COMPONENT  Component to Search for
  -f --queryfields=FIELDS   Fields that are used in "list" template: (default: summary,created,updated,priority,status,reporter,assignee)
  -i --issuetype=ISSUETYPE  The Issue Type
  -l --limit=VAL            Maximum number of results to return in query (default: 500)
  -p --project=PROJECT      Project to Search for
  -q --query=JQL            Jira Query Language expression for the search
  -r --reporter=USER        Reporter to search for
  -s --sort=ORDER           For list operations, sort issues (default: priority asc, created)
  -w --watcher=USER         Watcher to add to issue (default: $USER)
                            or Watcher to search for

Edit Options:
  -m --comment=COMMENT      Comment message for transition
  -o --override=KEY=VAL     Set custom key/value pairs

Create Options:
  -i --issuetype=ISSUETYPE  Jira Issue Type (default: Bug)
  -m --comment=COMMENT      Comment message for transition
  -o --override=KEY=VAL     Set custom key/value pairs

Command Options:
  -d --directory=DIR        Directory to export templates to (default: $HOME/.jira.d/templates)
```
