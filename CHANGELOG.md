# Changelog

## 0.1.15 - 2017-08-25

* merge in edit template changes from v1 branch [#105](https://github.com/Netflix-Skunkworks/go-jira/issues/105) [Cory Bennett] [[21dec2d](https://github.com/Netflix-Skunkworks/go-jira/commit/21dec2d)]
* Handle keyring.ErrNotFound error instead of panicing. [Will Rouesnel] [[338cc4d](https://github.com/Netflix-Skunkworks/go-jira/commit/338cc4d)]
* Use name field for fix version name value. Fixes [#89](https://github.com/Netflix-Skunkworks/go-jira/issues/89) [Bryan Baugher] [[fde82b2](https://github.com/Netflix-Skunkworks/go-jira/commit/fde82b2)]

## 0.1.14 - 2017-05-10

* fix unsafe casting for --quiet flag [Cory Bennett] [[6f29f43](https://github.com/Netflix-Skunkworks/go-jira/commit/6f29f43)]
* [[#80](https://github.com/Netflix-Skunkworks/go-jira/issues/80)] add `jira unassign` and `jira give ISSUE --default` commands [Cory Bennett] [[03d8633](https://github.com/Netflix-Skunkworks/go-jira/commit/03d8633)]

## 0.1.13 - 2017-04-24

* work around `github.com/tmc/keyring` compile error for windows [Cory Bennett] [[85298e9](https://github.com/Netflix-Skunkworks/go-jira/commit/85298e9)]
* Added generic issuelink command [David Reuss] [[cc54d11](https://github.com/Netflix-Skunkworks/go-jira/commit/cc54d11)]
* Added --start parameter for pagination on results [David Reuss] [[9b94d9e](https://github.com/Netflix-Skunkworks/go-jira/commit/9b94d9e)]

## 0.1.12 - 2017-03-22

* Implement "browse" subcommand on Windows [Claus Brod] [[ca333d8](https://github.com/Netflix-Skunkworks/go-jira/commit/ca333d8)]

## 0.1.11 - 2017-02-26

* [[#69](https://github.com/Netflix-Skunkworks/go-jira/issues/69)] add subtask command [Cory Bennett] [[21a2ed5](https://github.com/Netflix-Skunkworks/go-jira/commit/21a2ed5)]

## 0.1.10 - 2017-02-08

* set GPG_TTY in .bashrc [Cory Bennett] [[b1e552f](https://github.com/Netflix-Skunkworks/go-jira/commit/b1e552f)]
* force password in case password already exists [Cory Bennett] [[d5a2c3b](https://github.com/Netflix-Skunkworks/go-jira/commit/d5a2c3b)]
* refactor password source, allow for "pass" to be used, update tests to use `password-source: pass` [Cory Bennett] [[5a71939](https://github.com/Netflix-Skunkworks/go-jira/commit/5a71939)]

## 0.1.9 - 2016-12-18

* only warn about needing login when not already running the login command [Cory Bennett] [[6c24e55](https://github.com/Netflix-Skunkworks/go-jira/commit/6c24e55)]
* fix(http): Add proxy transport [William Hearn] [[4bd740b](https://github.com/Netflix-Skunkworks/go-jira/commit/4bd740b)] [[2dff6c9](https://github.com/Netflix-Skunkworks/go-jira/commit/2dff6c9)]

## 0.1.8 - 2016-11-24

* [[#12](https://github.com/Netflix-Skunkworks/go-jira/issues/12)] integrate with keyring for password storage and provide http basic auth credentials for every request since most jira services have websudo enabled with does not allow cookie based authentication [Cory Bennett] [[b8a6e57](https://github.com/Netflix-Skunkworks/go-jira/commit/b8a6e57)]
* Cleaning up usage [Jay Shirley] [[8add52b](https://github.com/Netflix-Skunkworks/go-jira/commit/8add52b)]
* Update usage [Jay Shirley] [[b56e32a](https://github.com/Netflix-Skunkworks/go-jira/commit/b56e32a)]
* use gopkg.in for links to maintain version compatibility [Cory Bennett] [[1414d1f](https://github.com/Netflix-Skunkworks/go-jira/commit/1414d1f)]
* golint [Cory Bennett] [[44cdebf](https://github.com/Netflix-Skunkworks/go-jira/commit/44cdebf)]
* add "rank" command allow ordering backlog issues in agile projects [Cory Bennett] [[e4cc9c6](https://github.com/Netflix-Skunkworks/go-jira/commit/e4cc9c6)]
* Adding a unixproxy mechanism [Jay Shirley] [[5b9c0dd](https://github.com/Netflix-Skunkworks/go-jira/commit/5b9c0dd)]

## 0.1.7 - 2016-08-24

* Prefer transition names which match exactly [Don Brower] [[e40f9c1](https://github.com/Netflix-Skunkworks/go-jira/commit/e40f9c1)]
* update tempates to make them more readable with space trimming added to go-1.6 [Cory Bennett] [[693b3e4](https://github.com/Netflix-Skunkworks/go-jira/commit/693b3e4)]

## 0.1.6 - 2016-08-21

* make "worklogs" command print output through template allow "add worklog" command to open edit template [Cory Bennett] [[cc3fbee](https://github.com/Netflix-Skunkworks/go-jira/commit/cc3fbee)]
* remove extra newline at end of worklogs template [Cory Bennett] [[d08ef15](https://github.com/Netflix-Skunkworks/go-jira/commit/d08ef15)]
* adding worklog related templates [Cory Bennett] [[ab1cd27](https://github.com/Netflix-Skunkworks/go-jira/commit/ab1cd27)]

## 0.1.5 - 2016-08-21

* update for golint [Cory Bennett] [[5a4e17c](https://github.com/Netflix-Skunkworks/go-jira/commit/5a4e17c)]
* fix for go vet [Cory Bennett] [[355fb42](https://github.com/Netflix-Skunkworks/go-jira/commit/355fb42)]

## 0.1.4 - 2016-08-12

* when running "dups" on a Process Management Project type, you have to start/stop the task to resolve it [Cory Bennett] [[2c91905](https://github.com/Netflix-Skunkworks/go-jira/commit/2c91905)]
* allow for defaultResolution option for transition command [Cory Bennett] [[a328c2d](https://github.com/Netflix-Skunkworks/go-jira/commit/a328c2d)]
* add "backlog" command for Kanban related Issues [Cory Bennett] [[5d39b23](https://github.com/Netflix-Skunkworks/go-jira/commit/5d39b23)]
* fix --noedit flag with "dups" command [Cory Bennett] [[37c07fa](https://github.com/Netflix-Skunkworks/go-jira/commit/37c07fa)]
* add "votes" and "labels" to default view template [Cory Bennett] [[6f73b8c](https://github.com/Netflix-Skunkworks/go-jira/commit/6f73b8c)]
* add "blockerType" config param, for issueLinkType use for "blocks" command [Cory Bennett] [[30fd301](https://github.com/Netflix-Skunkworks/go-jira/commit/30fd301)]
* update gitter room [Cory Bennett] [[4b822b1](https://github.com/Netflix-Skunkworks/go-jira/commit/4b822b1)]
* default issuetype to "Bug" for project that have Bug, otherwise try "Task" [Cory Bennett] [[0c807b4](https://github.com/Netflix-Skunkworks/go-jira/commit/0c807b4)]
* make view template only show fields that have values [Cory Bennett] [[8238fe8](https://github.com/Netflix-Skunkworks/go-jira/commit/8238fe8)]
* make default create template only display fields if they are valid fields for the project [Cory Bennett] [[adc2ace](https://github.com/Netflix-Skunkworks/go-jira/commit/adc2ace)]
* ignore empty json fields when processing templates [Cory Bennett] [[f5f3e28](https://github.com/Netflix-Skunkworks/go-jira/commit/f5f3e28)]
* allow JIRA_LOG_FORMAT env variable to control log output format [Cory Bennett] [[469def0](https://github.com/Netflix-Skunkworks/go-jira/commit/469def0)]
* remove extraneous debug [Cory Bennett] [[752a94d](https://github.com/Netflix-Skunkworks/go-jira/commit/752a94d)]
* add logout command modify password prompt to echo masked password [Cory Bennett] [[8ad91be](https://github.com/Netflix-Skunkworks/go-jira/commit/8ad91be)]
* tweak cookies to store hostname dump all http request/response with --verbose [Cory Bennett] [[f93fe79](https://github.com/Netflix-Skunkworks/go-jira/commit/f93fe79)]
* load configs in order of closest to cwd (/etc/go-jira.yml is last) [Cory Bennett] [[f54267b](https://github.com/Netflix-Skunkworks/go-jira/commit/f54267b)]

## 0.1.3 - 2016-07-30

* [[#43](https://github.com/Netflix-Skunkworks/go-jira/issues/43)] add support for jira done|todo|prog commands [Cory Bennett] [[dd7d1cc](https://github.com/Netflix-Skunkworks/go-jira/commit/dd7d1cc)]
* Reporter is not generally editable. [Mike Pountney] [[a637b43](https://github.com/Netflix-Skunkworks/go-jira/commit/a637b43)]

## 0.1.2 - 2016-06-29

* [[#44](https://github.com/Netflix-Skunkworks/go-jira/issues/44)] Close tmpfile before rename to work around "The process cannot access the file because it is being used by another process" error on windows. [Cory Bennett] [[0980f8e](https://github.com/Netflix-Skunkworks/go-jira/commit/0980f8e)]

## 0.1.1 - 2016-06-28

* use USERPROFILE instead of HOME for windows, rework paths to use filepath.Join for better cross platform support [Cory Bennett] [[adcedc4](https://github.com/Netflix-Skunkworks/go-jira/commit/adcedc4)]
* Include templates from a system path [Mike Pountney] [[cf10f53](https://github.com/Netflix-Skunkworks/go-jira/commit/cf10f53)]
* Added support for the ```expand``` option for Issues [tobyjoe] [[fb4afc9](https://github.com/Netflix-Skunkworks/go-jira/commit/fb4afc9)]
* change for api changes to go-logging [Cory Bennett] [[7bfc6e8](https://github.com/Netflix-Skunkworks/go-jira/commit/7bfc6e8)]
* Fix issuetype calls adding URL escaping [Jonathan Wright] [[e4a25e2](https://github.com/Netflix-Skunkworks/go-jira/commit/e4a25e2)]

## 0.1.0 - 2016-01-29

* Fixes [#32](https://github.com/Netflix-Skunkworks/go-jira/issues/32) - make path to cookieFile if it's not present [Mike Pountney] [[6644579](https://github.com/Netflix-Skunkworks/go-jira/commit/6644579)]
* Add component/components support: add and list for now. [Mike Pountney] [[d7b3226](https://github.com/Netflix-Skunkworks/go-jira/commit/d7b3226)]
* Tweak the CmdWatch contract and add watcher remove support [Mike Pountney] [[383847a](https://github.com/Netflix-Skunkworks/go-jira/commit/383847a)]
* Amend vote/unvote to be vote/vote --down [Mike Pountney] [[797edef](https://github.com/Netflix-Skunkworks/go-jira/commit/797edef)]
* Add 'vote' and 'unvote' [Mike Pountney] [[c95e66e](https://github.com/Netflix-Skunkworks/go-jira/commit/c95e66e)]

## 0.0.20 - 2016-01-21

* [issue [#28](https://github.com/Netflix-Skunkworks/go-jira/issues/28)] check to make sure we got back issuetypes for create metadata [Cory Bennett] [[ee0e780](https://github.com/Netflix-Skunkworks/go-jira/commit/ee0e780)]
* Add insecure option for TLS endpoints [Brian Lalor] [[6a88bb9](https://github.com/Netflix-Skunkworks/go-jira/commit/6a88bb9)]
* Correct naming of parameter: set/add/remove are actions. [Mike Pountney] [[303784f](https://github.com/Netflix-Skunkworks/go-jira/commit/303784f)]
* Tweak CmdLabels args so that magic happens with CLI [Mike Pountney] [[40a7c65](https://github.com/Netflix-Skunkworks/go-jira/commit/40a7c65)]
* Expose ViewTicket as per FindIssues [Mike Pountney] [[8977f3d](https://github.com/Netflix-Skunkworks/go-jira/commit/8977f3d)]
* Add exposed versions of getTemplate and runTemplate [Mike Pountney] [[da6cbd5](https://github.com/Netflix-Skunkworks/go-jira/commit/da6cbd5)]
* Add 'labels' command to set/add/remove labels [Mike Pountney] [[230b52d](https://github.com/Netflix-Skunkworks/go-jira/commit/230b52d)]
* Add a 'join' func to the template engine [Mike Pountney] [[a7820fe](https://github.com/Netflix-Skunkworks/go-jira/commit/a7820fe)]
* make "jira" golang package, move code from jira/cli to root, move jira/main.go to main/main.go [Cory Bennett] [[7268b9e](https://github.com/Netflix-Skunkworks/go-jira/commit/7268b9e)]

## 0.0.19 - 2015-12-09

* fix jira trans TRANS ISSUE (case sensitivity issue), also go fmt [Cory Bennett] [[3c30f3b](https://github.com/Netflix-Skunkworks/go-jira/commit/3c30f3b)]

## 0.0.18 - 2015-12-03

* need to default "quiet" to false [Cory Bennett] [[4f4a89b](https://github.com/Netflix-Skunkworks/go-jira/commit/4f4a89b)]

## 0.0.17 - 2015-12-03

* add --quiet command to not print the OK .. add --saveFile option to print the issue/link to a file on create command [Cory Bennett] [[c9ac162](https://github.com/Netflix-Skunkworks/go-jira/commit/c9ac162)]
* fix overrides [Cory Bennett] [[eaddfe6](https://github.com/Netflix-Skunkworks/go-jira/commit/eaddfe6)]
* add abstract request wrapper to allow you to access/process random apis supported by Jira but not yet supported by go-jira [Cory Bennett] [[90ef56a](https://github.com/Netflix-Skunkworks/go-jira/commit/90ef56a)]

## 0.0.16 - 2015-11-23

* jira edit should not require one arguemnt (allow for --query) [Cory Bennett] [[a1eb4a1](https://github.com/Netflix-Skunkworks/go-jira/commit/a1eb4a1)]

## 0.0.15 - 2015-11-23

* [[#17](https://github.com/Netflix-Skunkworks/go-jira/issues/17)] print usage on missing arguments [Cory Bennett] [[fd2a2fe](https://github.com/Netflix-Skunkworks/go-jira/commit/fd2a2fe)]

## 0.0.14 - 2015-11-17

* s/enpoint/endpoint/g [Oliver Schrenk] [[c5d251d](https://github.com/Netflix-Skunkworks/go-jira/commit/c5d251d)]
* Implement dateFormat template command [Mike Pountney] [[68d3bae](https://github.com/Netflix-Skunkworks/go-jira/commit/68d3bae)]
* Add 'updated' field to default queryfields. [Mike Pountney] [[91e2475](https://github.com/Netflix-Skunkworks/go-jira/commit/91e2475)]
* Fix export-templates option (typo) [Mike Pountney] [[4d7fdb8](https://github.com/Netflix-Skunkworks/go-jira/commit/4d7fdb8)]
* when yaml element resolves to "\n" strip it out so we dont post it to jira [Cory Bennett] [[47ced2f](https://github.com/Netflix-Skunkworks/go-jira/commit/47ced2f)]
* print PUT/POST data when using --dryrun to help debug [Cory Bennett] [[618f245](https://github.com/Netflix-Skunkworks/go-jira/commit/618f245)]

## 0.0.13 - 2015-09-19

* replace dead/deprecated code.google.com/p/gopass with golang.org/x/crypto/ssh/terminal for reading password from stdin [Cory Bennett] [[909eb06](https://github.com/Netflix-Skunkworks/go-jira/commit/909eb06)]

## 0.0.12 - 2015-09-18

* fix exception from "jira create" [Cory Bennett] [[9348a4b](https://github.com/Netflix-Skunkworks/go-jira/commit/9348a4b)]
* add some debug messages to help diagnose login failures [Cory Bennett] [[1c08a7d](https://github.com/Netflix-Skunkworks/go-jira/commit/1c08a7d)]

## 0.0.11 - 2015-09-16

* add --version [Cory Bennett] [[8385ee2](https://github.com/Netflix-Skunkworks/go-jira/commit/8385ee2)]
* fix command line parser broken in 0.0.10 [Cory Bennett] [[15ae929](https://github.com/Netflix-Skunkworks/go-jira/commit/15ae929)]

## 0.0.10 - 2015-09-15

* allow for command aliasing in conjunction with executable config files. Issue #5 [Cory Bennett] [[23590d4](https://github.com/Netflix-Skunkworks/go-jira/commit/23590d4)]
* update usage [Cory Bennett] [[ef7a57e](https://github.com/Netflix-Skunkworks/go-jira/commit/ef7a57e)]

## 0.0.9 - 2015-09-15

* use forked yaml.v2 so as to not lose line terminations present in jira fields [Cory Bennett] [[f84e77f](https://github.com/Netflix-Skunkworks/go-jira/commit/f84e77f)]
* adding a |~ literal yaml syntax to just chomp a single newline (again to preserve existing formatting in jira fields) [Cory Bennett] [[f84e77f](https://github.com/Netflix-Skunkworks/go-jira/commit/f84e77f)]
* for indent/comment allow for unicode line termination characters that yaml will use for parsing [Cory Bennett] [[f84e77f](https://github.com/Netflix-Skunkworks/go-jira/commit/f84e77f)]
* fix "edit" default option, change how defaults are dealt with for filters [Cory Bennett] [[4265913](https://github.com/Netflix-Skunkworks/go-jira/commit/4265913)]
* for edit template add issue id as comment, also add "comments" as comment so you can review the comment details while editing [Cory Bennett] [[968a9df](https://github.com/Netflix-Skunkworks/go-jira/commit/968a9df)]
* add "comment" template filter to comment out multiline statements [Cory Bennett] [[d664868](https://github.com/Netflix-Skunkworks/go-jira/commit/d664868)]
* add getOpt wrappers to get options with defaults [Cory Bennett] [[c0070cf](https://github.com/Netflix-Skunkworks/go-jira/commit/c0070cf)]
* make --dryrun work [Cory Bennett] [[d229ac1](https://github.com/Netflix-Skunkworks/go-jira/commit/d229ac1)]
* refactor config/option loading so command options override settings in config files [Cory Bennett] [[d229ac1](https://github.com/Netflix-Skunkworks/go-jira/commit/d229ac1)]
* allow query options to be used on the "edit" command to iterate editing [Cory Bennett] [[d229ac1](https://github.com/Netflix-Skunkworks/go-jira/commit/d229ac1)]
* remove duplication for defaults [Cory Bennett] [[f8c8ddf](https://github.com/Netflix-Skunkworks/go-jira/commit/f8c8ddf)]
* use optigo for option parsing, drop docopt [Cory Bennett] [[7bbd571](https://github.com/Netflix-Skunkworks/go-jira/commit/7bbd571)]
* allow "abort: true" to be set while editing to cancel the edit operation [Cory Bennett] [[ea67a77](https://github.com/Netflix-Skunkworks/go-jira/commit/ea67a77)]
* if no changes are made on edit templates then abort edit [Cory Bennett] [[e69b65c](https://github.com/Netflix-Skunkworks/go-jira/commit/e69b65c)]

## 0.0.8 - 2015-07-31

* Add --max_results option for 'ls' [Mike Pountney] [[e06ff0c](https://github.com/Netflix-Skunkworks/go-jira/commit/e06ff0c)]

## 0.0.7 - 2015-07-01

* fix "take" command not honouring user option [Andrew Haigh] [[8f1d2b9](https://github.com/Netflix-Skunkworks/go-jira/commit/8f1d2b9)]
* fix typo [Cory Bennett] [[06f57fe](https://github.com/Netflix-Skunkworks/go-jira/commit/06f57fe)]

## 0.0.6 - 2015-02-27

* allow --sort= to disable sort override [Cory Bennett] [[701f091](https://github.com/Netflix-Skunkworks/go-jira/commit/701f091)]
* fix default JIRA_OPERATION env variable [Cory Bennett] [[82fd9b9](https://github.com/Netflix-Skunkworks/go-jira/commit/82fd9b9)]
* automatically close duplicate issues with "Duplicate" resolution [Cory Bennett] [[ebf1700](https://github.com/Netflix-Skunkworks/go-jira/commit/ebf1700)]
* set JIRA_OPERATION to "view" when no operation used (ie: jira GOJIRA-123) [Cory Bennett] [[050848a](https://github.com/Netflix-Skunkworks/go-jira/commit/050848a)]
* add --sort option to "list" command [Cory Bennett] [[f359030](https://github.com/Netflix-Skunkworks/go-jira/commit/f359030)]

## 0.0.5 - 2015-02-21

* handle editor having arguments [Cory Bennett] [[7186fb3](https://github.com/Netflix-Skunkworks/go-jira/commit/7186fb3)]
* add more template error handling [Cory Bennett] [[3e6f2b3](https://github.com/Netflix-Skunkworks/go-jira/commit/3e6f2b3)]
* allow create template to specify defalt watchers with -o watchers=... [Cory Bennett] [[4db2e4e](https://github.com/Netflix-Skunkworks/go-jira/commit/4db2e4e)]
* if config files are executable then run them and parse the output [Cory Bennett] [[7a2f7f5](https://github.com/Netflix-Skunkworks/go-jira/commit/7a2f7f5)]

## 0.0.4 - 2015-02-19

* add --template option to export-templates to export a single template [Cory Bennett] [[343fbb6](https://github.com/Netflix-Skunkworks/go-jira/commit/343fbb6)]
* add "table" template to be used with "list" command [Cory Bennett] [[8954ec1](https://github.com/Netflix-Skunkworks/go-jira/commit/8954ec1)]

## 0.0.3 - 2015-02-19

* [issue [#8](https://github.com/Netflix-Skunkworks/go-jira/issues/8)] detect X-Seraph-Loginreason: AUTHENTICATION_DENIED header to catch login failures [Cory Bennett] [[2dcf665](https://github.com/Netflix-Skunkworks/go-jira/commit/2dcf665)]
* project should always be uppercase [Jay Buffington] [[1b69d12](https://github.com/Netflix-Skunkworks/go-jira/commit/1b69d12)]
* if response is 400, check json for errorMessages and log them [Jay Buffington] [[4924dfa](https://github.com/Netflix-Skunkworks/go-jira/commit/4924dfa)]
* validate project [Jay Buffington] [[dc5ae42](https://github.com/Netflix-Skunkworks/go-jira/commit/dc5ae42)]

## 0.0.2 - 2015-02-18

* add missing --override options on transition command
* add browse command

## 0.0.1 - 2015-02-18

* Initial experimental release
