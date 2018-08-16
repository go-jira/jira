[![Join the chat at https://gitter.im/go-jira-cli/help](https://badges.gitter.im/go-jira-cli/help.svg)](https://gitter.im/go-jira-cli/help?utm_source=badge&utm_medium=badge&utm_content=badge)
[![Build Status](https://travis-ci.org/Netflix-Skunkworks/go-jira.svg?branch=master)](https://travis-ci.org/Netflix-Skunkworks/go-jira)
[![GoDoc](https://godoc.org/gopkg.in/Netflix-Skunkworks/go-jira.v1?status.svg)](https://godoc.org/gopkg.in/Netflix-Skunkworks/go-jira.v1)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

### Table of Contents

   * [Summary](#go-jira)
   * [Install](#install)
      * [Download](#download)
      * [Build](#build)
   * [Usage](#usage)
      * [TAB completion](#setting-up-tab-completion)
   * [v1 vs v0 changes](#v1-vs-v0-changes)
      * [<strong>Golang library import</strong>](#golang-library-import)
      * [<strong>Configs per command</strong>](#configs-per-command)
      * [<strong>Custom Commands</strong>](#custom-commands)
      * [<strong>Incompatible command changes</strong>](#incompatible-command-changes)
      * [<strong>Login process change</strong>](#login-process-change)
   * [Configuration](#configuration)
      * [Dynamic Configuration](#dynamic-configuration)
      * [Custom Commands](#custom-commands-1)
         * [Commands](#commands)
         * [Options](#options)
         * [Arguments](#arguments)
         * [Script Template](#script-template)
         * [Examples](#examples)
      * [Editing](#editing)
      * [Templates](#templates)
         * [Writing/Editing Templates](#writingediting-templates)
      * [Authentication](#authentication)
         * [user vs login](#user-vs-login)
         * [keyring password source](#keyring-password-source)
         * [pass password source](#pass-password-source)

# go-jira
simple command line client for Atlassian's Jira service written in Go

## Install

### Download

You can download one of the pre-built binaries for **go-jira** [here](https://github.com/Netflix-Skunkworks/go-jira/releases).

### Build

You can build and install the official repository with [Go](https://golang.org/dl/):

```
go get gopkg.in/Netflix-Skunkworks/go-jira.v1/cmd/jira
```
This will checkout this repository into `$GOPATH/src/gopkg.in/Netflix-Skunkworks/go-jira.v1`, build, and install it.

Because golang likes fully qualified import paths, forking and contributing can be a bit tricky.

If you want to tinker or hack on go-jira, the [easiest way to do so](http://code.openark.org/blog/development/forking-golang-repositories-on-github-and-managing-the-import-path) is to fork the repository and clone directly into the official path like this:

`git clone https://github.com/YOUR_USER_NAME_HERE/go-jira $GOPATH/src/gopkg.in/Netflix-Skunkworks/go-jira.v1`

From within that source dir you can build and install modifications from within that directory like:

`go install ./...`

## Usage


#### Setting up TAB completion

Since go-jira is build with the "kingpin" golang command line library we supports bash/zsh shell completion automatically:

 * <https://github.com/alecthomas/kingpin/tree/v2.2.5#bashzsh-shell-completion>

For example, in bash, adding something along the lines of:

  `eval "$(jira --completion-script-bash)"`

to your bashrc, or .profile (assuming go-jira binary is already in your path) will cause jira to offer tab completion behavior.

```
usage: jira [<flags>] <command> [<args> ...]

Jira Command Line Interface

Global flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
  -v, --verbose ...            Increase verbosity for debugging
  -e, --endpoint=ENDPOINT      Base URI to use for Jira
  -k, --insecure               Disable TLS certificate verification
  -Q, --quiet                  Suppress output to console
      --unixproxy=UNIXPROXY    Path for a unix-socket proxy
      --socksproxy=SOCKSPROXY  Address for a socks proxy
  -u, --user=USER              user name used within the Jira service
      --login=LOGIN            login name that corresponds to the user used for authentication

Commands:
  help:                Show help.
  version:             Prints version
  acknowledge:         Transition issue to acknowledge state
  assign:              Assign user to issue
  attach create:       Attach file to issue
  attach get:          Fetch attachment
  attach list:         Prints attachment details for issue
  attach remove:       Delete attachment
  backlog:             Transition issue to Backlog state
  block:               Mark issues as blocker
  browse:              Open issue in browser
  close:               Transition issue to close state
  comment:             Add comment to issue
  component add:       Add component
  components:          Show components for a project
  create:              Create issue
  createmeta:          View 'create' metadata
  done:                Transition issue to Done state
  dup:                 Mark issues as duplicate
  edit:                Edit issue details
  editmeta:            View 'edit' metadata
  epic add:            Add issues to Epic
  epic create:         Create Epic
  epic list:           Prints list of issues for an epic with optional search criteria
  epic remove:         Remove issues from Epic
  export-templates:    Export templates for customizations
  fields:              Prints all fields, both System and Custom
  in-progress:         Transition issue to Progress state
  issuelink:           Link two issues
  issuelinktypes:      Show the issue link types
  issuetypes:          Show issue types for a project
  labels add:          Add labels to an issue
  labels remove:       Remove labels from an issue
  labels set:          Set labels on an issue
  list:                Prints list of issues for given search criteria
  login:               Attempt to login into jira server
  logout:              Deactivate session with Jira server
  rank:                Mark issues as blocker
  reopen:              Transition issue to reopen state
  request:             Open issue in requestr
  resolve:             Transition issue to resolve state
  start:               Transition issue to start state
  stop:                Transition issue to stop state
  subtask:             Subtask issue
  take:                Assign issue to yourself
  todo:                Transition issue to To Do state
  transition:          Transition issue to given state
  transitions:         List valid issue transitions
  transmeta:           List valid issue transitions
  unassign:            Unassign an issue
  unexport-templates:  Remove unmodified exported templates
  view:                Prints issue details
  vote:                Vote up/down an issue
  watch:               Add/Remove watcher to issue
  worklog add:         Add a worklog to an issue
  worklog list:        Prints the worklog data for given issue
  session:             Attempt to login into jira server

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
Instead of requiring a executable template to get configs for a given command now you can create a config to be applied to a command.  So if you want to use `template: table` by default for your `jira list` you can now do:
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
```yaml
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
Unfortunately during the rewrite between v0 and v1 there were some necessary changes that broke backwards compatibility with existing commands.  Specifically the `dups`, `blocks`, `add worklog` and `add|remove|set labels` commands have had the command word swapped around:
  * `jira DUPLICATE dups ISSUE` => `jira dup DUPLICATE ISSUE`
  * `jira BLOCKER blocks ISSUE` => `jira block BLOCKER ISSUE`
  * `jira add worklog` => `jira worklog add`
  * `jira add labels` => `jira labels add`
  * `jira remove labels` => `jira labels remove`
  * `jira set labels` => `jira labels set`

###### **Login process change**
We have, once again, changed how login happens for Jira.  When authenticating against Atlassian Cloud Jira [API Tokens are now required](https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-basic-auth-and-cookie-based-auth/).  Please read [the Authentication section](#authentication) below for more information.

If you use a privately hosted Jira service, you can chose to use the API Token method or continue using the session login api.  Please read [the Authentication section](#authentication) below for more information.

Previously `jira` used attempt to get a `JSESSION` cookies by authenticating with the webservice standard GUI login process.  This has been especially problematic as users need to authenticate with various credential providers (google auth, etc).  We now attempt to authenticate via the [session login api](https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-login).  This may be problematic for users if admins have locked down the session-login api, so we might have to bring back the error-prone Basic-Auth approach.  For users that are unable to authenticate via `jira` hopefully someone in your organization can provide me with details on a process for you to authenticate and we can try to update `jira`.

## Configuration

**go-jira** uses a configuration hierarchy.  When loading the configuration from disk it will recursively look through all parent directories in your current path looking for a **.jira.d** directory.  If your current directory is not a child directory of your homedir, then your homedir will also be inspected for a **.jira.d** directory.  From all of **.jira.d** directories discovered **go-jira** will load a **&lt;command&gt;.yml** file (ie for `jira list` it will load `.jira.d/list.yml`) then it will merge in any properties from the **config.yml** if found.  The configuration properties found in a file closest to your current working directory will have precedence.  Properties overridden with command line options will have final precedence.

The complicated configuration hierarchy is used because **go-jira** attempts to be context aware.  For example, if you are working on a "foo" project and you `cd` into your project workspace, wouldn't it be nice if `jira ls` automatically knew to list only issues related to the "foo" project?  Likewise when you `cd` to the "bar" project then `jira ls` should only list issues related to "bar" project.  You can do this with by creating a configuration under your project workspace at **./.jira.d/config.yml** that looks like:

```yaml
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
```yaml
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
* `help: string` The help message displayed in usage for the option.
* `type: string`:  The type of the option, can be one of these values: `BOOL`, `COUNTER`, `ENUM`, `FLOAT32`, `FLOAT64`, `INT8`, `INT16`, `INT32`, `INT64`, `INT`, `STRING`, `STRINGMAP`, `UINT8`, `UINT16`, `UINT32`,  `UINT64` and `UINT`.  Most of these are primitive data types an should be self-explanatory.  The default type is `STRING`. There are some special types:
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
* `name: string` [**required**] Name of the option, so `name: ISSUE` will show in the usage as `jira <command> ISSUE`.  This also represents the name of the argument to be used in the script template, so `{{args.ISSUE}}`.
* `help: string` The help message displayed in usage for the argument.
* `type: string`:  The type of the argument, can be one of these values: `BOOL`, `COUNTER`, `ENUM`, `FLOAT32`, `FLOAT64`, `INT8`, `INT16`, `INT32`, `INT64`, `INT`, `STRING`, `STRINGMAP`, `UINT8`, `UINT16`, `UINT32`,  `UINT64` and `UINT`.  Most of these are primitive data types an should be self-explanatory.  The default type is `STRING`.  There are some special types:
  * `COUNTER` will be an integer type that increments each the argument is provided  So something like `jira <command> ISSUE-12 ISSUE-23` will results in `{{args.ISSUE}}` of `2`.
  * `ENUM` type is used with the `enum` property.  The raw type is a string and **must** be one of the values listed in the `enum` property.
  * `STRINGMAP` is a `string => string` map with the format of `KEY=VALUE`.  So `jira <command> foo=bar bin=baz` along with a `name: OVERRIDE` property will allow for `{{args.OVERRIDE.foo}}` to be `bar` and `{{args.OVERRIDE.bin}}` to be `baz`.
* `required: bool` Indicate that this argument must be provided on the command line.  Conflicts with the `default` property.
* `default: any` Specify the default value for the argument.  Conflicts with the `required` property.
* `repeat: bool` Indicate that this argument can be repeated.  Not applicable for `COUNTER` and `STRINGMAP` types.  This will turn the template value into an array that you can iterate over.  So `jira <command> ISSUE-12 ISSUE-23` can be used like `{{range args.ISSUE}}Issue: {{.}}{{end}}`
* `enum: string list` Used with the `type: ENUM` property, it is a list of strings values that represent the set of possible values for the argument.

##### Script Template
The `script` property is a template that would produce `/bin/sh` compatible syntax after the template has been processed.  There are 2 key template functions `{{args}}` and `{{options}}` that return the parsed arguments and option flags as a map.

To demonstrate how you might use args and options here is a `custom-test` command:
```yaml
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
```yaml
custom-commands:
  - name: env
    help: print the JIRA environment variables available to custom commands
    script: |
      env | grep JIRA
 ```

You could use the environment variables automatically, so if your `.jira.d/config.yml` looks something like this:
```yaml
project: PROJECT
custom-commands:
  - name: print-project
    help: print the name of the configured project
    script: "echo $JIRA_PROJECT"
```

##### Examples

* `jira mine` for listing issues assigned to you
```yaml
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
```yaml
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
closes **go-jira** will submit the completed form.  The order which **go-jira** attempts to determine your preferred editor is:

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

First the basic templating functionality is defined by the Go language 'text/template' library.  The library reference documentation can be found [here](https://golang.org/pkg/text/template/), and there is a good primer document [here](https://gohugo.io/templates/go-templates/).  `go-jira` also provides a few extra helper functions to make it a bit easier to format the data, those functions are defined [here](https://github.com/Netflix-Skunkworks/go-jira/blob/master/jiracli/templates.go#L64).

Knowing what data and fields are available to any given template is not obvious. The easiest approach to determine what is available is to use the `debug` template on any given operation.  For example to find out what is available to the "view" templates, you can use:
```
jira view GOJIRA-321 -t debug
```

This will print out the data in JSON format that is available to the template.  You can do this for any other operation, like "list":
```
jira list -t debug
```

### Authentication

For Atlassian Cloud hosted Jira [API Tokens are now required](https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-basic-auth-and-cookie-based-auth/).  You will automatically be prompted for an API Token if your jira endpoint ends in `.atlassian.net`.  If you are using a private Jira service, you can force `jira` to use an api-token by setting the `authentication-method: api-token` property in your `$HOME/.jira.d/config.yml` file.  The API Token needs to be presented to the Jira service on every request, so it is recommended to store this API Token security within your OS's keyring, or using the `pass` service as documented below so that it can be programmatically accessed via `jira` and not prompt you every time.  For a less-secure option you can also provide the API token via a `JIRA_API_TOKEN` environment variable.  If you are unable to use an api-token for an Atlassian Cloud hosted Jira then you can still force `jira` to use the old session based authentication (until it the hosted system stops accepting it) by setting `authentication-method: session`.

If your Jira service still allows you to use the Session based authentication method then `jira` will prompt for a password automatically when get a response header from the Jira service that indicates you do not have an active session (ie the `X-Ausername` header is set to `anonymous`).  Then after authentication we cache the `cloud.session.token` cookie returned by the service [session login api](https://docs.atlassian.com/jira/REST/cloud/#auth/1/session-login) and reuse that on subsequent requests.  Typically this cookie will be valid for several hours (depending on the service configuration).  To automatically securely store your password for easy reuse by jira You can enable a `password-source` via `.jira.d/config.yml` with possible values of `keyring` or `pass`.

#### User vs Login
The Jira service has sometimes differing opinions about how a user is identified.  In other words the ID you login with might not be ID that the jira system recognized you as.  This matters when trying to identify a user via various Jira REST APIs (like issue assignment).  This is especially relevant when trying to authenticate with an API Token where the authentication user is usually an email address, but within the Jira system the user is identified by a user name.  To accommodate this `jira` now supports two different properties in the config file.  So when authentication using the API Tokens you will likely want something like this in your `$HOME/.jira.d/config.yml` file:
```yaml
user: person
login: person@example.com
```

You can also override these values on the command line with `jira --user person --login person@example.com`.  The `login` value will be used only for authentication purposes, the `user` value will be used when a user name is required for any Jira service API calls.

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

You probably want to setup gpg-agent so that you don't have to type in your gpg passphrase all the time.  You can get `gpg-agent` to automatically start by adding something like this to your `$HOME/.bashrc`
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
