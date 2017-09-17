[![Join the chat at https://gitter.im/go-jira-cli/help](https://badges.gitter.im/go-jira-cli/help.svg)](https://gitter.im/go-jira-cli/help?utm_source=badge&utm_medium=badge&utm_content=badge)
[![Build Status](https://travis-ci.org/Netflix-Skunkworks/go-jira.svg?branch=master)](https://travis-ci.org/Netflix-Skunkworks/go-jira)
[![GoDoc](https://godoc.org/gopkg.in/Netflix-Skunkworks/go-jira.v1?status.png)](https://godoc.org/gopkg.in/Netflix-Skunkworks/go-jira.v1)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

# go-jira
simple command line client for Atlassian's Jira service written in Go

## Install

### Download

You can download one of the pre-built binaries for **go-jira** [here](https://github.com/Netflix-Skunkworks/go-jira/releases).

### Build

You can build and install with [Go](https://golang.org/dl/):

```
go get gopkg.in/Netflix-Skunkworks/go-jira.v1/cmd/jira
```

## v1 vs v0 changes

###### **Golang library import**
For the new version of go-jira you should use:
```
import "gopkg.in/Netflix-Skunkworks/go-jira.v1"
```

If you have code that depends on the old apis, you can still use them with this import:
```
import "gopkg.in/Netflix-Skunkworks/go-jira.v0"
```

###### **Configs per command**
Instead of requiring a exectuable template to get configs for a given command now you can create a config to be applied to a command.  So if you want to use `template: table` by default for yor `jira list` you can now do:
```
$ cat $HOME/.jira.d/list.yml
template: table
```
Where previously you needed something like:
```
# cat $HOME/.jira.d/config.yml
#!/bin/sh
case $JIRA_OPERATION in 
    list)
      echo "template: table";;
esac
```

###### **Custom Commands**
Now you can create your own custom commands to do common operations with jira.  Please see the details **Custom Commands** section below for more details.  If you want to create a command `jira mine` that lists all the issues assigned to you now you can modify your `.jira.d/config.yml` file to add a `custom-commands` section like this:
```
custom-commands:
  - name: mine
    help: display issues assigned to me
    script: |-
      {{jira}} list --query "resolution = unresolved and assignee=currentuser() ORDER BY created"
```
Then the next time you run `jira help` you will see your usage:
```
$ jira mine --help
usage: jira mine

display issues assigned to me

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose ...          Increase verbosity for debugging
  -e, --endpoint=ENDPOINT    Base URI to use for Jira
  -u, --user=USER            Login name used for authentication with Jira service
      --unixproxy=UNIXPROXY  Path for a unix-socket proxy
  -k, --insecure             Disable TLS certificate verification
```

###### **Incompatible command changes**
Unfortunately during the rewrite between v0 and v1 there were some changes necessary that broke backwards compatibility with existing commands.  Specifically the `dups`, `blocks`, `add worklog` and `add|remove|set labels` commands have had the command word swapped around:
  * `jira DUPLICATE dups ISSUE` => `jira dup DUPLICATE ISSUE`
  * `jira BLOCKER blocks ISSUE` => `jira block BLOCKER ISSUE`
  * `jira add worklog` => `jira worklog add`
  * `jira add labels` => `jira labels add`
  * `jira remove labels` => `jira labels remove`
  * `jira set labels` => `jira labels set`

###### **Login process change**
Previously `jira` used attempt to get a `JSESSION` cookies by authenticating with the webservice standard GUI login process.  This has been especially problematic as users need to authenticate with various credential providers (google auth, etc).  We now attempt to authenticate via the [session login api](https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-login).  This may be problematic for users if admins have locked down the session-login api, so we might have to bring back the error-prone Basic-Auth approach.  For users that are unable to authenticate via `jira` hopefully someone in your organization can provide me with details on a process for you to authenticate and we can try to update `jira`.

## Configuration

**go-jira** uses a configuration hierarchy.  When loading the configuration from disk it will recursively look through
all parent directories in your current path looking for a **.jira.d** directory.  If your current directory is not
a child directory of your homedir, then your homedir will also be inspected for a **.jira.d** directory.  From all of **.jira.d** directories
discovered **go-jira** will load a **&lt;command&gt;.yml** file (ie for `jira list` it will load `.jira.d/list.yml`) then it will merge in any properties from the **config.yml** if found.  The configuration properties found in a file closests to your current working directory
will have precedence.  Properties overriden with command line options will have final precedence.

The complicated configuration hierarchy is used because **go-jira** attempts to be context aware.  For example, if you are working on a "foo" project and
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

### Custom Commands
You can now create custom commands for `jira` just by editing your `.jira.d/config.yml` config file.  These commands are effectively shell-scripts that can have documented options and arguments. The basic format is like:
```
custom-commands:
  - command1
  - command2
```
##### Commands
Where the individual commands are maps with these keys:
* `name: string` [**required**] This is the command name, so for `jira foobar` you would have `name: foobar`
* `help: string` This is help message displayed in the usage for the command
* `hidden: bool` This command will be hidden from users, but still executable.  Sometimes useful for constructing complex commands where one custom command might call another.
* `default: bool` Use this for compound command groups.  If you wanted to have `jira foo bar` and `jira foo baz` you would have two commands with `name: foo bar` and `name: foo baz`.  Then if you wanted `jira foo baz` to be called by default when you type `jira foo` you would set `default: true` for that custom command.
* `options: list`  This is the list of possible option flags that the command will accept
* `args: list` This is the list of command arguments (like the ISSUE) that the command will accept.
* `aliases: string list`: This is a list of alternate names that the user can provide on the command line to run the same command.  Typically used to shorten the command name or provide alternatives that users might expect. 
* `script: string` [**required**] This is the script that will be executed as the action for this command. The value will be treated as a template and substitutions for options and arguments will be made before executing.

##### Options
These are possible keys under the command `options` property:
* `name: string` [**required**] Name of the option, so `name: foobar` will result in `--foobar` option.
* `help: string` The help messsage displayed in usage for the option.
* `type: string`:  The type of the option, can be one of these values: `BOOL`, `COUNTER`, `ENUM`, `FLOAT32`, `FLOAT64`, `INT8`, `INT16`, `INT32`, `INT64`, `INT`, `STRING`, `STRINGMAP`, `UINT8`, `UINT16`, `UINT32`,  `UINT64` and `UINT`.  Most of these are primitive data types an should be self-explanitory.  The default type is `STRING`. There are some special types:
  * `COUNTER` will be an integer type that increments each time the option is used.  So something like `--count --count` will results in `{{options.count}}` of `2`.
  * `ENUM` type is used with the `enum` property.  The raw type is a string and **must** be one of the values listed in the `enum` property.
  * `STRINGMAP` is a `string => string` map with the format of `KEY=VALUE`.  So `--override foo=bar --override bin=baz` will allow for `{{options.override.foo}}` to be `bar` and `{{options.override.bin}}` to be `baz`.
* `short: char` The single character option to be used so `short: c` will allow for `-c`.
* `required: bool` Indicate that this option must be provided on the command line.  Conflicts with the `default` property.
* `default: any` Specify the default value for the option.  Conflicts with the `required` property.
* `hidden: bool` Hide the option from the usage help message, but otherwise works fine.  Sometimes useful for developer options that user should not play with.
* `repeat: bool` Indicate that this option can be repeated.  Not applicable for `COUNTER` and `STRINGMAP` types.  This will turn the option value into an array that you can iterate over.  So `--day Monday --day Thursday` can be used like `{{range options.day}}Day: {{.}}{{end}}`
* `enum: string list` Used with the `type: ENUM` property, it is a list of strings values that represent the set of possible values the option accepts.

##### Arguments
These are possible keys under the command `args` property:
* `name: string` [**required**] Name of the option, so `name: ISSUE` will show in the usasge as `jira <command> ISSUE`.  This also represents the name of the argument to be used in the script template, so `{{args.ISSUE}}`.
* `help: string` The help messsage displayed in usage for the argument.
* `type: string`:  The type of the argumemnt, can be one of these values: `BOOL`, `COUNTER`, `ENUM`, `FLOAT32`, `FLOAT64`, `INT8`, `INT16`, `INT32`, `INT64`, `INT`, `STRING`, `STRINGMAP`, `UINT8`, `UINT16`, `UINT32`,  `UINT64` and `UINT`.  Most of these are primitive data types an should be self-explanitory.  The default type is `STRING`.  There are some special types:
  * `COUNTER` will be an integer type that increments each the argument is provided  So something like `jira <command> ISSUE-12 ISSUE-23` will results in `{{args.ISSUE}}` of `2`.
  * `ENUM` type is used with the `enum` property.  The raw type is a string and **must** be one of the values listed in the `enum` property.
  * `STRINGMAP` is a `string => string` map with the format of `KEY=VALUE`.  So `jira <command> foo=bar bin=baz` along with a `name: OVERRIDE` property will allow for `{{args.OVERRIDE.foo}}` to be `bar` and `{{args.OVERRIDE.bin}}` to be `baz`.
* `required: bool` Indicate that this argument must be provided on the command line.  Conflicts with the `default` property.
* `default: any` Specify the default value for the argument.  Conflicts with the `required` property.
* `repeat: bool` Indicate that this argument can be repeated.  Not applicable for `COUNTER` and `STRINGMAP` types.  This will turn the template value into an array that you can iterate over.  So `jira <command> ISSUE-12 ISSUE-23` can be used like `{{range args.ISSUE}}Issue: {{.}}{{end}}`
* `enum: string list` Used with the `type: ENUM` property, it is a list of strings values that represent the set of possible values for the argument.

##### Script Template
The `script` property is a template that whould produce `/bin/sh` compatible syntax after the template has been processed.  There are 2 key template functions `{{args}}` and `{{options}}` that return the parsed arguments and option flags as a map.  

To demonstrate how you might use args and options here is a `custom-test` command:
```
custom-commands:
  - name: custom-test
    help: Testing the custom commands
    options:
      - name: abc
        short: a
        default: default
      - name: day
        type: ENUM
        enum:
          - Monday
          - Tuesday
          - Wednesday
          - Thursday
          - Friday
        required: true
    args:
      - name: ARG
        required: true
      - name: MORE
        repeat: true
    script: |
      echo COMMAND {{args.ARG}} --abc {{options.abc}} --day {{options.day}} {{range $more := args.MORE}}{{$more}} {{end}}
```

Then to run it:
```
$ jira custom-test
ERROR Invalid Usage: required flag --day not provided

$ jira custom-test --day Sunday
ERROR Invalid Usage: enum value must be one of Monday,Tuesday,Wednesday,Thursday,Friday, got 'Sunday'

$ jira custom-test --day Tuesday
ERROR Invalid Usage: required argument 'ARG' not provided

$ jira custom-test --day Tuesday arg1
COMMAND arg1 --abc default --day Tuesday

$ jira custom-test --day Tuesday arg1 more1 more2 more3
COMMAND arg1 --abc default --day Tuesday more1 more2 more3

$ jira custom-test --day Tuesday arg1 more1 more2 more3 --abc non-default
COMMAND arg1 --abc non-default --day Tuesday more1 more2 more3

$ jira custom-test --day Tuesday arg1 more1 more2 more3 -a short-non-default
COMMAND arg1 --abc short-non-default --day Tuesday more1 more2 more3
```

The script has access to all the environment variables that are in your current environment plus those that `jira` will set.  `jira` sets environment variables for each config property it has parsed from `.jira.d/config.yml` or the command configs at `.jira.d/<command>.yml`.  It might be useful to see all environment variables that `jira` is producing, so here is a simple custom command to list them:
```
custom-commands:
  - name: env
    help: print the JIRA environment variables available to custom commands
    script: |
      env | grep JIRA
 ```
 
You could use the environment variables automatically, so if your `.jira.d/config.yml` looks something like this:
```
project: PROJECT
custom-commands:
  - name: print-project
    help: print the name of the configured project
    script: "echo $JIRA_PROJECT"
```

##### Examples

* `jira mine` for listing issues assigned to you
```
custom-commands:
  - name: mine
    help: display issues assigned to me
    script: |-
      if [ -n "$JIRA_PROJECT" ]; then
          # if `project: ...` configured just list the issues for current project
          {{jira}} list --template table --query "resolution = unresolved and assignee=currentuser() and project = $JIRA_PROJECT ORDER BY priority asc, created"
      else
          # otherwise list issues for all project
          {{jira}} list --template table --query "resolution = unresolved and assignee=currentuser() ORDER BY priority asc, created"
      fi
```
* `jira sprint` for listing issues in your current sprint
```
custom-commands:
  - name: sprint
    help: display issues for active sprint
    script: |-
      if [ -n "$JIRA_PROJECT" ]; then
          # if `project: ...` configured just list the issues for current project
          {{jira}} list --template table --query "sprint in openSprints() and type != epic and resolution = unresolved and project=$JIRA_PROJECT ORDER BY rank asc, created"
      else
          # otherwise list issues for all project
          echo "\"project: ...\" configuration missing from .jira.d/config.yml"
      fi
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

First the basic templating functionality is defined by the Go language 'text/template' library.  The library reference documentation can be found [here](https://golang.org/pkg/text/template/), and there is a good primer document [here](https://gohugo.io/templates/go-templates/).  `go-jira` also provides a few extra helper functions to make it a bit easlier to format the data, those functions are defined [here](https://github.com/Netflix-Skunkworks/go-jira/blob/master/jiracli/templates.go#L64).

Knowing what data and fields are available to any given template is not obvious. The easiest approach to determine what is available is to use the `debug` template on any given operation.  For eample to find out what is available to the "view" templates, you can use:
```
jira view GOJIRA-321 -t debug
```

This will print out the data in JSON format that is available to the template.  You can do this for any other operation, like "list":
```
jira list -t debug
```

### Authentication

By default `go-jira` will prompt for a password automatically when get a response header from the Jira service that indicates you do not have an active session (ie the `X-Ausername` header is set to `anonymous`).  Then after authentication we cache the `cloud.session.token` cookie returned by the service [session login api](https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-login) and reuse that on subsequent requests.  Typically this cookie will be valid for several hours (depending on the service configuration).  To automatically securely store your password for easy reuse by jira You can enable a `password-source` via `.jira.d/config.yml` with possible values of `keyring` or `pass`.

#### keyring password source
On OSX and Linux there are a few keyring providers that `go-jira` can use (via this [golang module](https://github.com/tmc/keyring)).  To integrate `go-jira` with a supported keyring just add this configuration to `$HOME/.jira.d/config.yml`:
```yaml
password-source: keyring
```
After setting this and issuing a `jira login`, your credentials will be stored in your platform's backend (e.g. Keychain for Mac OS X) automatically. Subsequent operations, like a `jira ls`, should automatically login.

#### `pass` password source
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

if [ ! -f $HOME/.gpg-agent.conf ]; then
  cat <<EOM >$HOME/.gpg-agent.conf
default-cache-ttl 604800
max-cache-ttl 604800
default-cache-ttl-ssh 604800
max-cache-ttl-ssh 604800
EOM
fi

if [ -n "${GPG_AGENT_INFO}" ]; then
    nc  -U "${GPG_AGENT_INFO%%:*}" >/dev/null </dev/null
    if [ ! -S "${GPG_AGENT_INFO%%:*}" -o $? != 0 ]; then
        # set passphrase cache so I only have to type my passphrase once a day
        eval $(gpg-agent --options $HOME/.gpg-agent.conf --daemon --write-env-file $HOME/.gpg-agent-info --use-standard-socket --log-file $HOME/tmp/gpg-agent.log --verbose)
    fi
fi
export GPG_TTY=$(tty)
```

## Usage

```
usage: jira [<flags>] <command> [<args> ...]

Jira Command Line Interface

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose ...          Increase verbosity for debugging
  -e, --endpoint=ENDPOINT    Base URI to use for Jira
  -k, --insecure             Disable TLS certificate verification
  -Q, --quiet                Suppress output to console
      --unixproxy=UNIXPROXY  Path for a unix-socket proxy
  -u, --user=USER            Login name used for authentication with Jira service

Commands:
  help [<command>...]
    Show help.


  version
    Prints version


  login
    Attempt to login into jira server


  logout
    Deactivate sesssion with Jira server


  list [<flags>]
    Prints list of issues for given search criteria

    -t, --template=TEMPLATE        Template to use for output
        --gjq=GJQ                  GJSON Query to filter output, see https://goo.gl/iaYwJ5
    -a, --assignee=ASSIGNEE        User assigned the issue
    -c, --component=COMPONENT      Component to search for
    -i, --issuetype=ISSUETYPE      Issue type to search for
    -l, --limit=LIMIT              Maximum number of results to return in search
    -p, --project=PROJECT          Project to search for
    -n, --named-query=NAMED-QUERY  The name of a query in the `queries` configuration
    -q, --query=QUERY              Jira Query Language (JQL) expression for the search
    -f, --queryfields=QUERYFIELDS  Fields that are used in "list" template
    -r, --reporter=REPORTER        Reporter to search for
    -S, --status=STATUS            Filter on issue status
    -s, --sort=SORT                Sort order to return
    -w, --watcher=WATCHER          Watcher to search for

  view [<flags>] <ISSUE>
    Prints issue details

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --gjq=GJQ                GJSON Query to filter output, see https://goo.gl/iaYwJ5
        --expand=EXPAND ...      field to expand for the issue
        --field=FIELD ...        field to return for the issue
        --property=PROPERTY ...  property to return for issue

  create [<flags>]
    Create issue

    -b, --browse                 Open issue(s) in browser after operation
        --editor=EDITOR          Editor to use
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -p, --project=PROJECT        project to create issue in
    -i, --issuetype=ISSUETYPE    issuetype in to create
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property
        --saveFile=SAVEFILE      Write issue as yaml to file

  edit [<flags>] [<ISSUE>]
    Edit issue details

    -b, --browse                   Open issue(s) in browser after operation
        --editor=EDITOR            Editor to use
    -t, --template=TEMPLATE        Template to use for output
        --noedit                   Disable opening the editor
    -n, --named-query=NAMED-QUERY  The name of a query in the `queries` configuration
    -q, --query=QUERY              Jira Query Language (JQL) expression for the search to edit multiple issues
    -m, --comment=COMMENT          Comment message for issue
    -o, --override=OVERRIDE ...    Set issue property

  comment [<flags>] [<ISSUE>]
    Add comment to issue

    -b, --browse             Open issue(s) in browser after operation
        --editor=EDITOR      Editor to use
    -t, --template=TEMPLATE  Template to use for output
        --noedit             Disable opening the editor
    -m, --comment=COMMENT    Comment message for issue

  epic create [<flags>]
    Create Epic

    -b, --browse                 Open issue(s) in browser after operation
        --editor=EDITOR          Editor to use
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -p, --project=PROJECT        project to create epic in
    -n, --epic-name=EPIC-NAME    Epic Name
    -m, --comment=COMMENT        Comment message for epic
    -o, --override=OVERRIDE ...  Set epic property
        --saveFile=SAVEFILE      Write epic as yaml to file

  epic list [<flags>] <EPIC>
    Prints list of issues for an epic with optional search criteria

    -t, --template=TEMPLATE        Template to use for output
        --gjq=GJQ                  GJSON Query to filter output, see https://goo.gl/iaYwJ5
    -a, --assignee=ASSIGNEE        User assigned the issue
    -c, --component=COMPONENT      Component to search for
    -i, --issuetype=ISSUETYPE      Issue type to search for
    -l, --limit=LIMIT              Maximum number of results to return in search
    -p, --project=PROJECT          Project to search for
    -n, --named-query=NAMED-QUERY  The name of a query in the `queries` configuration
    -q, --query=QUERY              Jira Query Language (JQL) expression for the search
    -f, --queryfields=QUERYFIELDS  Fields that are used in "list" template
    -r, --reporter=REPORTER        Reporter to search for
    -S, --status=STATUS            Filter on issue status
    -s, --sort=SORT                Sort order to return
    -w, --watcher=WATCHER          Watcher to search for

  epic add <EPIC> <ISSUE>...
    Add issues to Epic


  epic remove <ISSUE>...
    Remove issues from Epic


  worklog list [<flags>] <ISSUE>
    Prints the worklog data for given issue

    -b, --browse             Open issue(s) in browser after operation
    -t, --template=TEMPLATE  Template to use for output
        --gjq=GJQ            GJSON Query to filter output, see https://goo.gl/iaYwJ5

  worklog add [<flags>] <ISSUE>
    Add a worklog to an issue

    -b, --browse                 Open issue(s) in browser after operation
        --editor=EDITOR          Editor to use
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for worklog
    -T, --time-spent=TIME-SPENT  Time spent working on issue
    -S, --started=STARTED        Time you started work

  fields [<flags>]
    Prints all fields, both System and Custom

    -t, --template=TEMPLATE  Template to use for output
        --gjq=GJQ            GJSON Query to filter output, see https://goo.gl/iaYwJ5

  createmeta [<flags>]
    View 'create' metadata

    -t, --template=TEMPLATE    Template to use for output
        --gjq=GJQ              GJSON Query to filter output, see https://goo.gl/iaYwJ5
    -p, --project=PROJECT      project to fetch create metadata
    -i, --issuetype=ISSUETYPE  issuetype in project to fetch create metadata

  editmeta [<flags>] <ISSUE>
    View 'edit' metadata

    -b, --browse             Open issue(s) in browser after operation
    -t, --template=TEMPLATE  Template to use for output
        --gjq=GJQ            GJSON Query to filter output, see https://goo.gl/iaYwJ5

  subtask [<flags>] [<ISSUE>]
    Subtask issue

    -b, --browse                 Open issue(s) in browser after operation
        --editor=EDITOR          Editor to use
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -p, --project=PROJECT        project to subtask issue in
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  dup [<flags>] <DUPLICATE> <ISSUE>
    Mark issues as duplicate

    -b, --browse             Open issue(s) in browser after operation
        --editor=EDITOR      Editor to use
    -t, --template=TEMPLATE  Template to use for output
    -m, --comment=COMMENT    Comment message when marking issue as duplicate

  block [<flags>] <BLOCKER> <ISSUE>
    Mark issues as blocker

    -b, --browse             Open issue(s) in browser after operation
        --editor=EDITOR      Editor to use
    -t, --template=TEMPLATE  Template to use for output
    -m, --comment=COMMENT    Comment message when marking issue as blocker

  issuelink [<flags>] <OUTWARDISSUE> <ISSUELINKTYPE> <INWARDISSUE>
    Link two issues

    -b, --browse             Open issue(s) in browser after operation
        --editor=EDITOR      Editor to use
    -t, --template=TEMPLATE  Template to use for output
    -m, --comment=COMMENT    Comment message when linking issue

  issuelinktypes [<flags>]
    Show the issue link types

    -t, --template=TEMPLATE  Template to use for output
        --gjq=GJQ            GJSON Query to filter output, see https://goo.gl/iaYwJ5

  transition [<flags>] <TRANSITION> <ISSUE>
    Transition issue to given state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  transitions [<flags>] <ISSUE>
    List valid issue transitions

    -b, --browse             Open issue(s) in browser after operation
    -t, --template=TEMPLATE  Template to use for output
        --gjq=GJQ            GJSON Query to filter output, see https://goo.gl/iaYwJ5

  transmeta [<flags>] <ISSUE>
    List valid issue transitions

    -b, --browse             Open issue(s) in browser after operation
    -t, --template=TEMPLATE  Template to use for output
        --gjq=GJQ            GJSON Query to filter output, see https://goo.gl/iaYwJ5

  close [<flags>] <ISSUE>
    Transition issue to close state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  acknowledge [<flags>] <ISSUE>
    Transition issue to acknowledge state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  reopen [<flags>] <ISSUE>
    Transition issue to reopen state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  resolve [<flags>] <ISSUE>
    Transition issue to resolve state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  start [<flags>] <ISSUE>
    Transition issue to start state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  stop [<flags>] <ISSUE>
    Transition issue to stop state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  todo [<flags>] <ISSUE>
    Transition issue to To Do state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  backlog [<flags>] <ISSUE>
    Transition issue to Backlog state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  done [<flags>] <ISSUE>
    Transition issue to Done state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  in-progress [<flags>] <ISSUE>
    Transition issue to Progress state

    -b, --browse                 Open issue(s) in browser after operation
    -t, --template=TEMPLATE      Template to use for output
        --noedit                 Disable opening the editor
    -m, --comment=COMMENT        Comment message for issue
    -o, --override=OVERRIDE ...  Set issue property

  vote [<flags>] [<ISSUE>]
    Vote up/down an issue

    -b, --browse  Open issue(s) in browser after operation
    -d, --down    downvote the issue

  rank [<flags>] <FIRST-ISSUE> <after|before> <SECOND-ISSUE>
    Mark issues as blocker

    -b, --browse  Open issue(s) in browser after operation

  watch [<flags>] <ISSUE> [<WATCHER>]
    Add/Remove watcher to issue

    -b, --browse  Open issue(s) in browser after operation
    -r, --remove  remove watcher from issue

  labels add [<flags>] <ISSUE> <LABEL>...
    Add labels to an issue

    -b, --browse  Open issue(s) in browser after operation

  labels set [<flags>] <ISSUE> <LABEL>...
    Set labels on an issue

    -b, --browse  Open issue(s) in browser after operation

  labels remove [<flags>] <ISSUE> <LABEL>...
    Remove labels from an issue

    -b, --browse  Open issue(s) in browser after operation

  take [<flags>] <ISSUE> [<ASSIGNEE>]
    Assign issue to yourself

    -b, --browse   Open issue(s) in browser after operation
        --default  use default user for assignee

  assign [<flags>] <ISSUE> [<ASSIGNEE>]
    Assign user to issue

    -b, --browse   Open issue(s) in browser after operation
        --default  use default user for assignee

  unassign [<flags>] <ISSUE> [<ASSIGNEE>]
    Unassign an issue

    -b, --browse   Open issue(s) in browser after operation
        --default  use default user for assignee

  component add [<flags>]
    Add component

        --editor=EDITOR            Editor to use
    -t, --template=TEMPLATE        Template to use for output
        --noedit                   Disable opening the editor
    -p, --project=PROJECT          project to create component in
    -n, --name=NAME                name of component
    -d, --description=DESCRIPTION  description of component
    -l, --lead=LEAD                person that acts as lead for component

  components [<flags>]
    Show components for a project

    -t, --template=TEMPLATE  Template to use for output
        --gjq=GJQ            GJSON Query to filter output, see https://goo.gl/iaYwJ5
    -p, --project=PROJECT    project to list components

  issuetypes [<flags>]
    Show issue types for a project

    -t, --template=TEMPLATE  Template to use for output
        --gjq=GJQ            GJSON Query to filter output, see https://goo.gl/iaYwJ5
    -p, --project=PROJECT    project to list issueTypes

  attach create [<flags>] <ISSUE> [<ATTACHMENT>]
    Attach file to issue

    -b, --browse             Open issue(s) in browser after operation
        --saveFile=SAVEFILE  Write attachment information as yaml to file
    -f, --filename=FILENAME  Filename to use for attachment

  attach list [<flags>] <ISSUE>
    Prints issue details

    -b, --browse             Open issue(s) in browser after operation
    -t, --template=TEMPLATE  Template to use for output

  attach get [<flags>] [<ATTACHMENT-ID>]
    Fetch attachment

    -o, --output=OUTPUT  Write attachment to specified file name, '-' for stdout

  attach remove [<ATTACHMENT-ID>]
    Delete attachment


  export-templates [<flags>]
    Export templates for customizations

    -t, --template=TEMPLATE  Template to export
    -d, --dir=DIR            directory to write tempates to

  unexport-templates [<flags>]
    Remove unmodified exported templates

    -t, --template=TEMPLATE  Template to export
    -d, --dir=DIR            directory to write tempates to

  browse <ISSUE>
    Open issue in browser


  request [<flags>] <API> [<JSON>]
    Open issue in requestr

    -t, --template=TEMPLATE  Template to use for output
        --gjq=GJQ            GJSON Query to filter output, see https://goo.gl/iaYwJ5
    -M, --method=METHOD      HTTP request method to use


```
