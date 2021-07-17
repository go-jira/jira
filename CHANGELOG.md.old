# Changelog

## 1.0.25 - 2020-09-11

* Ensure body is NPE safe [Louis DeLosSantos] [[42e5d23](https://github.com/Netflix-Skunkworks/go-jira/commit/42e5d23)]
* Support empty responses in request commands [Louis DeLosSantos] [[b572037](https://github.com/Netflix-Skunkworks/go-jira/commit/b572037)]

## 1.0.24 - 2020-09-04

* Make -h flag show --help [Benjamin Kane] [[4bf1d03](https://github.com/Netflix-Skunkworks/go-jira/commit/4bf1d03)]
* transition: map field name to id [Louis DeLosSantos] [[3c1c4d9](https://github.com/Netflix-Skunkworks/go-jira/commit/3c1c4d9)]
* username-deprecation: use email and display names [Louis DeLosSantos] [[6a27e28](https://github.com/Netflix-Skunkworks/go-jira/commit/6a27e28)]
* Add support to get all comments for an issue [Thibault Jamet] [[a311d0d](https://github.com/Netflix-Skunkworks/go-jira/commit/a311d0d)]
* update all usage of user.name to user.accountId for privacy migration: https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-user-privacy-api-migration-guide/ [Cory Bennett] [[a26683e](https://github.com/Netflix-Skunkworks/go-jira/commit/a26683e)]
* add template functions to handle table output, fixes [#176](https://github.com/Netflix-Skunkworks/go-jira/issues/176), replaces [#296](https://github.com/Netflix-Skunkworks/go-jira/issues/296) [Cory Bennett] [[7e97463](https://github.com/Netflix-Skunkworks/go-jira/commit/7e97463)]
* use `password-source-path` to allow overriding binary and/or path to binary [Cory Bennett] [[d6173ce](https://github.com/Netflix-Skunkworks/go-jira/commit/d6173ce)]
* allow issues on command line to automatically prefix with project when defined [Cory Bennett] [[d002d7f](https://github.com/Netflix-Skunkworks/go-jira/commit/d002d7f)]
* Forgot you use TAB instead of spaces [Cory Bennett] [[789886c](https://github.com/Netflix-Skunkworks/go-jira/commit/789886c)]
* Fixed append project to view [Cory Bennett] [[8a46215](https://github.com/Netflix-Skunkworks/go-jira/commit/8a46215)]
* Added a line break removal function [Cory Bennett] [[9cbd993](https://github.com/Netflix-Skunkworks/go-jira/commit/9cbd993)]
* Pushed Readfile to .jira.d directory instead of pwd [Cory Bennett] [[db53622](https://github.com/Netflix-Skunkworks/go-jira/commit/db53622)]
* Cache password to avoid invoking password source on each API request [Patrick Decat] [[0f059a5](https://github.com/Netflix-Skunkworks/go-jira/commit/0f059a5)]
* Add support to switch out password source binary [Patrick Pichler] [[659a5c8](https://github.com/Netflix-Skunkworks/go-jira/commit/659a5c8)]
* Add error handling to pass password-source [Patrick Pichler] [[3339303](https://github.com/Netflix-Skunkworks/go-jira/commit/3339303)]
* Add gopass support [Patrick Pichler] [[3c0a62e](https://github.com/Netflix-Skunkworks/go-jira/commit/3c0a62e)]
* add sprig template functions, replaces [[#215](https://github.com/Netflix-Skunkworks/go-jira/issues/215)] http://masterminds.github.io/sprig/ [Cory Bennett] [[719f7a6](https://github.com/Netflix-Skunkworks/go-jira/commit/719f7a6)]
* [[#232](https://github.com/Netflix-Skunkworks/go-jira/issues/232)] fix api to use common Version schema also rewrote the schema fetcher to use Go [Cory Bennett] [[90f01ce](https://github.com/Netflix-Skunkworks/go-jira/commit/90f01ce)]
* fix syntax [Cory Bennett] [[94dd489](https://github.com/Netflix-Skunkworks/go-jira/commit/94dd489)]
* Address comments for direct name match [Patrick Cockwell] [[a70384b](https://github.com/Netflix-Skunkworks/go-jira/commit/a70384b)]
* Choose exact transition match if available [Patrick Cockwell] [[a646f76](https://github.com/Netflix-Skunkworks/go-jira/commit/a646f76)]
* [[#277](https://github.com/Netflix-Skunkworks/go-jira/issues/277)] update figtree to latest [Cory Bennett] [[0e520a4](https://github.com/Netflix-Skunkworks/go-jira/commit/0e520a4)]
* Switch over to using github.com/go-jira/jira, from gopkg.in [Mike Pountney] [[27f57b2](https://github.com/Netflix-Skunkworks/go-jira/commit/27f57b2)]
* Add 'pctOf' and 'fit' template functions [Adriano Caloiaro] [[027adee](https://github.com/Netflix-Skunkworks/go-jira/commit/027adee)]
* fix worklog template to allow multiline comments [Matthias Bethke] [[43e07f1](https://github.com/Netflix-Skunkworks/go-jira/commit/43e07f1)]
* Allow reading password from stdin [Justin Ko] [[225e1dc](https://github.com/Netflix-Skunkworks/go-jira/commit/225e1dc)]
* all: unindent some code [Daniel Martí] [[31c113d](https://github.com/Netflix-Skunkworks/go-jira/commit/31c113d)]
* don't use ReadAll when decoding JSON [Daniel Martí] [[9bcdcc1](https://github.com/Netflix-Skunkworks/go-jira/commit/9bcdcc1)]
* start making staticcheck happier [Daniel Martí] [[9b9186f](https://github.com/Netflix-Skunkworks/go-jira/commit/9b9186f)]
* all: convert to a Go module [Daniel Martí] [[f125ef3](https://github.com/Netflix-Skunkworks/go-jira/commit/f125ef3)]
* CI: test on Go 1.12.x, cleanup [Daniel Martí] [[664c5ca](https://github.com/Netflix-Skunkworks/go-jira/commit/664c5ca)]
* make automatic pagination on search optional, fix tests [Cory Bennett] [[36c99ce](https://github.com/Netflix-Skunkworks/go-jira/commit/36c99ce)]
* prefer defer resp.Body.Close to avoid leaks on subsequent errors [Cory Bennett] [[181bd74](https://github.com/Netflix-Skunkworks/go-jira/commit/181bd74)]
* add pagination to search [Miles Maddox] [[76dd1d8](https://github.com/Netflix-Skunkworks/go-jira/commit/76dd1d8)]
* Fix function comments based on best practices from Effective Go [CodeLingo Bot] [[23ac118](https://github.com/Netflix-Skunkworks/go-jira/commit/23ac118)]

## 1.0.23 - 2020-02-26

* update all usage of user.name to user.accountId for privacy migration: https://developer.atlassian.com/cloud/jira/platform/deprecation-notice-user-privacy-api-migration-guide/ [Cory Bennett] [[a26683e](https://github.com/Netflix-Skunkworks/go-jira/commit/a26683e)]
* add template functions to handle table output, fixes [#176](https://github.com/Netflix-Skunkworks/go-jira/issues/176), replaces [#296](https://github.com/Netflix-Skunkworks/go-jira/issues/296) [Cory Bennett] [[7e97463](https://github.com/Netflix-Skunkworks/go-jira/commit/7e97463)]
* use `password-source-path` to allow overriding binary and/or path to binary [Cory Bennett] [[d6173ce](https://github.com/Netflix-Skunkworks/go-jira/commit/d6173ce)]
* allow issues on command line to automatically prefix with project when defined [Cory Bennett] [[d002d7f](https://github.com/Netflix-Skunkworks/go-jira/commit/d002d7f)]
* Forgot you use TAB instead of spaces [Cory Bennett] [[789886c](https://github.com/Netflix-Skunkworks/go-jira/commit/789886c)]
* Fixed append project to view [Cory Bennett] [[8a46215](https://github.com/Netflix-Skunkworks/go-jira/commit/8a46215)]
* Added a line break removal function [Cory Bennett] [[9cbd993](https://github.com/Netflix-Skunkworks/go-jira/commit/9cbd993)]
* Pushed Readfile to .jira.d directory instead of pwd [Cory Bennett] [[db53622](https://github.com/Netflix-Skunkworks/go-jira/commit/db53622)]
* Cache password to avoid invoking password source on each API request [Patrick Decat] [[0f059a5](https://github.com/Netflix-Skunkworks/go-jira/commit/0f059a5)]
* Add support to switch out password source binary [Patrick Pichler] [[659a5c8](https://github.com/Netflix-Skunkworks/go-jira/commit/659a5c8)]
* Add error handling to pass password-source [Patrick Pichler] [[3339303](https://github.com/Netflix-Skunkworks/go-jira/commit/3339303)]
* Add gopass support [Patrick Pichler] [[3c0a62e](https://github.com/Netflix-Skunkworks/go-jira/commit/3c0a62e)]
* add sprig template functions, replaces [[#215](https://github.com/Netflix-Skunkworks/go-jira/issues/215)] http://masterminds.github.io/sprig/ [Cory Bennett] [[719f7a6](https://github.com/Netflix-Skunkworks/go-jira/commit/719f7a6)]
* [[#232](https://github.com/Netflix-Skunkworks/go-jira/issues/232)] fix api to use common Version schema also rewrote the schema fetcher to use Go [Cory Bennett] [[90f01ce](https://github.com/Netflix-Skunkworks/go-jira/commit/90f01ce)]

## 1.0.22 - 2019-09-30

* fix syntax [Cory Bennett] [[807ca76](https://github.com/Netflix-Skunkworks/go-jira/commit/807ca76)]
* Address comments for direct name match [Patrick Cockwell] [[a70384b](https://github.com/Netflix-Skunkworks/go-jira/commit/a70384b)]
* Choose exact transition match if available [Patrick Cockwell] [[a646f76](https://github.com/Netflix-Skunkworks/go-jira/commit/a646f76)]


## 1.0.21 - 2019-09-16

* [[#277](https://github.com/Netflix-Skunkworks/go-jira/issues/277)] update figtree to latest [Cory Bennett] [[0e520a4](https://github.com/Netflix-Skunkworks/go-jira/commit/0e520a4)]
* Switch over to using github.com/go-jira/jira, from gopkg.in [Mike Pountney] [[27f57b2](https://github.com/Netflix-Skunkworks/go-jira/commit/27f57b2)]
* fix worklog template to allow multiline comments [Matthias Bethke] [[43e07f1](https://github.com/Netflix-Skunkworks/go-jira/commit/43e07f1)]
* Allow reading password from stdin [Justin Ko] [[225e1dc](https://github.com/Netflix-Skunkworks/go-jira/commit/225e1dc)]
* all: unindent some code [Daniel Martí] [[31c113d](https://github.com/Netflix-Skunkworks/go-jira/commit/31c113d)]
* don't use ReadAll when decoding JSON [Daniel Martí] [[9bcdcc1](https://github.com/Netflix-Skunkworks/go-jira/commit/9bcdcc1)]
* start making staticcheck happier [Daniel Martí] [[9b9186f](https://github.com/Netflix-Skunkworks/go-jira/commit/9b9186f)]
* all: convert to a Go module [Daniel Martí] [[f125ef3](https://github.com/Netflix-Skunkworks/go-jira/commit/f125ef3)]
* CI: test on Go 1.12.x, cleanup [Daniel Martí] [[664c5ca](https://github.com/Netflix-Skunkworks/go-jira/commit/664c5ca)]
* make automatic pagination on search optional, fix tests [Cory Bennett] [[36c99ce](https://github.com/Netflix-Skunkworks/go-jira/commit/36c99ce)]
* prefer defer resp.Body.Close to avoid leaks on subsequent errors [Cory Bennett] [[181bd74](https://github.com/Netflix-Skunkworks/go-jira/commit/181bd74)]
* add pagination to search [Miles Maddox] [[76dd1d8](https://github.com/Netflix-Skunkworks/go-jira/commit/76dd1d8)]
* Fix function comments based on best practices from Effective Go [CodeLingo Bot] [[23ac118](https://github.com/Netflix-Skunkworks/go-jira/commit/23ac118)]
* [[#201](https://github.com/Netflix-Skunkworks/go-jira/issues/201)] update required library, failing to populate cookiejar from cookies file [Cory Bennett] [[ee69afa](https://github.com/Netflix-Skunkworks/go-jira/commit/ee69afa)]

## 1.0.20 - 2018-08-04

* [[#201](https://github.com/Netflix-Skunkworks/go-jira/issues/201)] update required library, failing to populate cookiejar from cookies file [Cory Bennett] [[ee69afa](https://github.com/Netflix-Skunkworks/go-jira/commit/ee69afa)]

## 1.0.19 - 2018-08-02

* [[#199](https://github.com/Netflix-Skunkworks/go-jira/issues/199)] [[#198](https://github.com/Netflix-Skunkworks/go-jira/issues/198)] update http client library, fix issues with internal login retries [Cory Bennett] [[0cba806](https://github.com/Netflix-Skunkworks/go-jira/commit/0cba806)]

## 1.0.18 - 2018-07-29

* [[#196](https://github.com/Netflix-Skunkworks/go-jira/issues/196)] add `jira session` command to show session information if user is authenticated [Cory Bennett] [[f592107](https://github.com/Netflix-Skunkworks/go-jira/commit/f592107)]
* [[#192](https://github.com/Netflix-Skunkworks/go-jira/issues/192)] fix template usage for 'fixVersions' in transition template [Cory Bennett] [[c9a4e30](https://github.com/Netflix-Skunkworks/go-jira/commit/c9a4e30)]
* move HandleExit to the jiracli package [Cory Bennett] [[ef9b731](https://github.com/Netflix-Skunkworks/go-jira/commit/ef9b731)]
* they broke golang.org/x/net: ERROR: vendor/golang.org/x/net/proxy/socks5.go:11:2: use of internal package not allowed [Cory Bennett] [[7191c77](https://github.com/Netflix-Skunkworks/go-jira/commit/7191c77)]
* udpate deps [Cory Bennett] [[d16bcc2](https://github.com/Netflix-Skunkworks/go-jira/commit/d16bcc2)]
* udpate for figtree api changes [Cory Bennett] [[07ba74b](https://github.com/Netflix-Skunkworks/go-jira/commit/07ba74b)]
* udpate to use latest dep, update figtree [Cory Bennett] [[462ef1c](https://github.com/Netflix-Skunkworks/go-jira/commit/462ef1c)]
* [[#171](https://github.com/Netflix-Skunkworks/go-jira/issues/171)] change proposed PasswordPath to PasswordName [Cory Bennett] [[213a7e0](https://github.com/Netflix-Skunkworks/go-jira/commit/213a7e0)]
* add pass path to config [dvogt23] [[fa01ff5](https://github.com/Netflix-Skunkworks/go-jira/commit/fa01ff5)]

## 1.0.17 - 2018-04-15

* fix IsTerminal usage for windows [Cory Bennett] [[7f9595c](https://github.com/Netflix-Skunkworks/go-jira/commit/7f9595c)]
* [[#166](https://github.com/Netflix-Skunkworks/go-jira/issues/166)] fix issue when editing templates specified with full path [Cory Bennett] [[d787ac0](https://github.com/Netflix-Skunkworks/go-jira/commit/d787ac0)]
* only prompt on logout if stdin and stdout are terminals [Cory Bennett] [[09a61c3](https://github.com/Netflix-Skunkworks/go-jira/commit/09a61c3)]
* [[#163](https://github.com/Netflix-Skunkworks/go-jira/issues/163)] fix url path join logic [Cory Bennett] [[9146346](https://github.com/Netflix-Skunkworks/go-jira/commit/9146346)]
* [[#160](https://github.com/Netflix-Skunkworks/go-jira/issues/160)] prompt when api-token is invalid to get new token [Cory Bennett] [[e639cce](https://github.com/Netflix-Skunkworks/go-jira/commit/e639cce)]
* [[#157](https://github.com/Netflix-Skunkworks/go-jira/issues/157)] add `password-directory: path` to allow overriding PASSWORD_STORE_DIR from configs [Cory Bennett] [[06b26c9](https://github.com/Netflix-Skunkworks/go-jira/commit/06b26c9)]
* [[#160](https://github.com/Netflix-Skunkworks/go-jira/issues/160)] allow `jira logout` to delete your api-token from keychain [Cory Bennett] [[bd3cf99](https://github.com/Netflix-Skunkworks/go-jira/commit/bd3cf99)]

## 1.0.16 - 2018-04-01

* [[#159](https://github.com/Netflix-Skunkworks/go-jira/issues/159)] fix `slice bounds out of range` error in `abbrev` template function [Cory Bennett] [[359bec2](https://github.com/Netflix-Skunkworks/go-jira/commit/359bec2)]
* [[#158](https://github.com/Netflix-Skunkworks/go-jira/issues/158)] always print usage to stdout [Cory Bennett] [[79c83f6](https://github.com/Netflix-Skunkworks/go-jira/commit/79c83f6)]

## 1.0.15 - 2018-03-08

* [[#147](https://github.com/Netflix-Skunkworks/go-jira/issues/147)] [[#148](https://github.com/Netflix-Skunkworks/go-jira/issues/148)] add support for api token based authentication [Cory Bennett] [[edb0662](https://github.com/Netflix-Skunkworks/go-jira/commit/edb0662)]
* refactor to simplify main [Cory Bennett] [[43ebc84](https://github.com/Netflix-Skunkworks/go-jira/commit/43ebc84)] [[0d7c1a7](https://github.com/Netflix-Skunkworks/go-jira/commit/0d7c1a7)]
* [[#145](https://github.com/Netflix-Skunkworks/go-jira/issues/145)] fix to match AuthProvider interface [Cory Bennett] [[80325a5](https://github.com/Netflix-Skunkworks/go-jira/commit/80325a5)]
* [[#141](https://github.com/Netflix-Skunkworks/go-jira/issues/141)] better handling in responseError for non-json error responses [Cory Bennett] [[20a9666](https://github.com/Netflix-Skunkworks/go-jira/commit/20a9666)]
* Update unexportTemplates.go [GitHub] [[6da9974](https://github.com/Netflix-Skunkworks/go-jira/commit/6da9974)]
* [[#139](https://github.com/Netflix-Skunkworks/go-jira/issues/139)] add shellquote and toMinJson template functions [Cory Bennett] [[8c7ca38](https://github.com/Netflix-Skunkworks/go-jira/commit/8c7ca38)]
* [[#137](https://github.com/Netflix-Skunkworks/go-jira/issues/137)] update kingpeon dep to allow access to dynamic command structure [Cory Bennett] [[425fa63](https://github.com/Netflix-Skunkworks/go-jira/commit/425fa63)]
* field name is "comment" not "comments" [Cory Bennett] [[464742c](https://github.com/Netflix-Skunkworks/go-jira/commit/464742c)]

## 1.0.14 - 2017-11-04

* [[#131](https://github.com/Netflix-Skunkworks/go-jira/issues/131)] fix parsing global options before command execution (allow unixproxy/socksproxy to be set in config.yml) [Cory Bennett] [[042bc48](https://github.com/Netflix-Skunkworks/go-jira/commit/042bc48)]
* add/update deps [Cory Bennett] [[a2e36e8](https://github.com/Netflix-Skunkworks/go-jira/commit/a2e36e8)]
* add support for using socks proxy [onionjake] [[ff985f9](https://github.com/Netflix-Skunkworks/go-jira/commit/ff985f9)]

## 1.0.13 - 2017-10-28

* fix transition command [Cory Bennett] [[9597f9b](https://github.com/Netflix-Skunkworks/go-jira/commit/9597f9b)]
* fix default values to load after parsing configs [Cory Bennett] [[c9b5054](https://github.com/Netflix-Skunkworks/go-jira/commit/c9b5054)]
* add test to make sure IssueType.Fields does not disappear on regeneration [Cory Bennett] [[3966def](https://github.com/Netflix-Skunkworks/go-jira/commit/3966def)]
* add tests for validating changes to auto-generated jiradata files [Cory Bennett] [[41d1a7c](https://github.com/Netflix-Skunkworks/go-jira/commit/41d1a7c)]
* Fix typo in 'logout' command help [Cory Bennett] [[9000777](https://github.com/Netflix-Skunkworks/go-jira/commit/9000777)]
* Add URL escaping to an additional issuetype call [Cory Bennett] [[7bfa241](https://github.com/Netflix-Skunkworks/go-jira/commit/7bfa241)]
* Add --resolution option [Cory Bennett] [[e6600cf](https://github.com/Netflix-Skunkworks/go-jira/commit/e6600cf)]
* Create Metadata Not Populated Correctly [Dillon Buchanan] [[093c510](https://github.com/Netflix-Skunkworks/go-jira/commit/093c510)]
* add regexReplace template function [Dirk Heilig] [[d3e294e](https://github.com/Netflix-Skunkworks/go-jira/commit/d3e294e)]

## 1.0.12 - 2017-10-04

* add `{{env.VARNAME}}` template support to allow use of env vars [Cory Bennett] [[4d74554](https://github.com/Netflix-Skunkworks/go-jira/commit/4d74554)]

## 1.0.11 - 2017-09-26

* [[#115](https://github.com/Netflix-Skunkworks/go-jira/issues/115)] fix transition template for description [Cory Bennett] [[986cc78](https://github.com/Netflix-Skunkworks/go-jira/commit/986cc78)]
* update edit command to set queryFields on search to match what is used in template [Cory Bennett] [[3913726](https://github.com/Netflix-Skunkworks/go-jira/commit/3913726)]
* fix edit with query loop, allow continuation when not submitting previous issue [Cory Bennett] [[0ba8aa0](https://github.com/Netflix-Skunkworks/go-jira/commit/0ba8aa0)]
* fix edit when priority is not set [Cory Bennett] [[098eb99](https://github.com/Netflix-Skunkworks/go-jira/commit/098eb99)]
* flatten CommandRegistry list to make it more readable [Cory Bennett] [[2ddaed2](https://github.com/Netflix-Skunkworks/go-jira/commit/2ddaed2)]

## 1.0.10 - 2017-09-18

* clean up usage formatting, print aliases [Cory Bennett] [[9f433ac](https://github.com/Netflix-Skunkworks/go-jira/commit/9f433ac)]
* fix edit [Cory Bennett] [[4c6b36c](https://github.com/Netflix-Skunkworks/go-jira/commit/4c6b36c)]
* fix named query template expansion [Cory Bennett] [[a8eaa97](https://github.com/Netflix-Skunkworks/go-jira/commit/a8eaa97)]
* fix usage message [Cory Bennett] [[cd3cfd8](https://github.com/Netflix-Skunkworks/go-jira/commit/cd3cfd8)]

## 1.0.9 - 2017-09-17

* need issuetype to use the default list table template now [Cory Bennett] [[3e8b9bd](https://github.com/Netflix-Skunkworks/go-jira/commit/3e8b9bd)]
* [[#102](https://github.com/Netflix-Skunkworks/go-jira/issues/102)] add issuetype into the default queryfields and add it to the default `table` list template [Cory Bennett] [[c9d8dfb](https://github.com/Netflix-Skunkworks/go-jira/commit/c9d8dfb)]

## 1.0.8 - 2017-09-17

* [[#100](https://github.com/Netflix-Skunkworks/go-jira/issues/100)] add support for posting, fetching, listing and removing attachments [Cory Bennett] [[66eb7bf](https://github.com/Netflix-Skunkworks/go-jira/commit/66eb7bf)]

## 1.0.7 - 2017-09-15

* [[#87](https://github.com/Netflix-Skunkworks/go-jira/issues/87)] add various commands for interacting with epics [Cory Bennett] [[893454f](https://github.com/Netflix-Skunkworks/go-jira/commit/893454f)]

## 1.0.6 - 2017-09-13

* tweaks for templates in named queries to work better [Cory Bennett] [[00cba79](https://github.com/Netflix-Skunkworks/go-jira/commit/00cba79)]
* [[#99](https://github.com/Netflix-Skunkworks/go-jira/issues/99)] add support for named queries to be stored in configs [Cory Bennett] [[fb43753](https://github.com/Netflix-Skunkworks/go-jira/commit/fb43753)]
* [[#98](https://github.com/Netflix-Skunkworks/go-jira/issues/98)] add `--status` option for JQL filter on status with `list` command [Cory Bennett] [[5da04c1](https://github.com/Netflix-Skunkworks/go-jira/commit/5da04c1)]

## 1.0.5 - 2017-09-11

* use --gjq for GJson Query to filter json response data [Cory Bennett] [[608e586](https://github.com/Netflix-Skunkworks/go-jira/commit/608e586)]
* fix field tag syntax [Cory Bennett] [[2c552ac](https://github.com/Netflix-Skunkworks/go-jira/commit/2c552ac)]
* add '{{jira}}' template macro to refer to path of currently running jira command [Cory Bennett] [[941824d](https://github.com/Netflix-Skunkworks/go-jira/commit/941824d)]

## 1.0.4 - 2017-09-08

* update deps for kingpeon update use os.exec instead of syscall.exec for windows [Cory Bennett] [[86b963b](https://github.com/Netflix-Skunkworks/go-jira/commit/86b963b)]

## 1.0.3 - 2017-09-06

* [[#66](https://github.com/Netflix-Skunkworks/go-jira/issues/66)] add --started option to `jira worklog add` to change the start time for worklog [Cory Bennett] [[e6faee1](https://github.com/Netflix-Skunkworks/go-jira/commit/e6faee1)]
* [[#45](https://github.com/Netflix-Skunkworks/go-jira/issues/45)] automatically add comment to issue even if transition does not support comment updates during transtion [Cory Bennett] [[c4be59c](https://github.com/Netflix-Skunkworks/go-jira/commit/c4be59c)]

## 1.0.2 - 2017-09-06

* update dependencies [Cory Bennett] [[aa876cd](https://github.com/Netflix-Skunkworks/go-jira/commit/aa876cd)]
* update for github.com/AlecAivazis/survey => gopkg.in/AlecAivazis/survey.v1 package [Cory Bennett] [[9453179](https://github.com/Netflix-Skunkworks/go-jira/commit/9453179)]
* use stdout to determin output terminal size [Cory Bennett] [[4d79af4](https://github.com/Netflix-Skunkworks/go-jira/commit/4d79af4)]

## 1.0.1 - 2017-09-06

* [[#13](https://github.com/Netflix-Skunkworks/go-jira/issues/13)] change default input syntax to not require escaping for special characters [Cory Bennett] [[1106558](https://github.com/Netflix-Skunkworks/go-jira/commit/1106558)]

## 1.0.0 - 2017-09-05

* fix build for windows [Cory Bennett] [[1b854da](https://github.com/Netflix-Skunkworks/go-jira/commit/1b854da)]
* change the default log output format [Cory Bennett] [[f1b8c64](https://github.com/Netflix-Skunkworks/go-jira/commit/f1b8c64)]
* tweak auto-login so it does not print the standard `jira login` command output [Cory Bennett] [[49f6cdc](https://github.com/Netflix-Skunkworks/go-jira/commit/49f6cdc)]
* add --quiet global option [Cory Bennett] [[c226077](https://github.com/Netflix-Skunkworks/go-jira/commit/c226077)]
* refactor to allow for --insecure and --unixproxy arguments [Cory Bennett] [[c0358eb](https://github.com/Netflix-Skunkworks/go-jira/commit/c0358eb)]
* handle html response on expired cookies (require X-Ausername header to always be present) [Cory Bennett] [[21920c5](https://github.com/Netflix-Skunkworks/go-jira/commit/21920c5)]
* allow login prompt to be interrupted [Cory Bennett] [[7ab6c22](https://github.com/Netflix-Skunkworks/go-jira/commit/7ab6c22)]
* fmt -> log typo [Cory Bennett] [[bccf09f](https://github.com/Netflix-Skunkworks/go-jira/commit/bccf09f)]
* make ~/.jira.d directory if not already present [Cory Bennett] [[e72479c](https://github.com/Netflix-Skunkworks/go-jira/commit/e72479c)]
* fix go vet [Cory Bennett] [[e04b506](https://github.com/Netflix-Skunkworks/go-jira/commit/e04b506)]
* fix tests [Cory Bennett] [[ba35f55](https://github.com/Netflix-Skunkworks/go-jira/commit/ba35f55)]
* add OK printf [Cory Bennett] [[dc02181](https://github.com/Netflix-Skunkworks/go-jira/commit/dc02181)]
* change --method to use -M for backwards compat [Cory Bennett] [[b120c0b](https://github.com/Netflix-Skunkworks/go-jira/commit/b120c0b)]
* add resolution to dup'd issues when necessary [Cory Bennett] [[2638396](https://github.com/Netflix-Skunkworks/go-jira/commit/2638396)]
* call correct function for `labels remove|set` commands [Cory Bennett] [[ad1a62a](https://github.com/Netflix-Skunkworks/go-jira/commit/ad1a62a)]
* data argument is optional (for GET and DELETE requests) [Cory Bennett] [[4b60313](https://github.com/Netflix-Skunkworks/go-jira/commit/4b60313)]
* fix usage, overrides not serialized correctly [Cory Bennett] [[84119a2](https://github.com/Netflix-Skunkworks/go-jira/commit/84119a2)]
* fix `jira ISSUE-123` command line parsing [Cory Bennett] [[fa4ac25](https://github.com/Netflix-Skunkworks/go-jira/commit/fa4ac25)]
* add logger object to jiracmd [Cory Bennett] [[aed952b](https://github.com/Netflix-Skunkworks/go-jira/commit/aed952b)]
* refactor for GlobalOptions and CommonOptions [Cory Bennett] [[979da1f](https://github.com/Netflix-Skunkworks/go-jira/commit/979da1f)]
* move commands from jiracli package to jiracmd package [Cory Bennett] [[0a5510b](https://github.com/Netflix-Skunkworks/go-jira/commit/0a5510b)]
* use jiracli.Error object to disambiguate between kingpin errors and cli errors [Cory Bennett] [[fb1bfeb](https://github.com/Netflix-Skunkworks/go-jira/commit/fb1bfeb)]
* fix stray newline for list table template [Cory Bennett] [[36c26c5](https://github.com/Netflix-Skunkworks/go-jira/commit/36c26c5)]
* fix dynamic table output when not on tty [Cory Bennett] [[3942f6f](https://github.com/Netflix-Skunkworks/go-jira/commit/3942f6f)]
* when using --verbose set the JIRA_DEBUG environment variable so custom-commands can auto enable verbose output [Cory Bennett] [[da9a2b2](https://github.com/Netflix-Skunkworks/go-jira/commit/da9a2b2)]
* make `jira ISSUE-123` usage call `jira view ISSUE-123` [Cory Bennett] [[ec0858b](https://github.com/Netflix-Skunkworks/go-jira/commit/ec0858b)]
* integrate kingpeon library to allow for custom commands via configuration [Cory Bennett] [[301a61f](https://github.com/Netflix-Skunkworks/go-jira/commit/301a61f)]
* use terminal width to adjust list table output [Cory Bennett] [[2a081dd](https://github.com/Netflix-Skunkworks/go-jira/commit/2a081dd)]
* set yaml/json tags for option structs [Cory Bennett] [[f52d2c4](https://github.com/Netflix-Skunkworks/go-jira/commit/f52d2c4)]
* update generated data files [Cory Bennett] [[c89f11d](https://github.com/Netflix-Skunkworks/go-jira/commit/c89f11d)]
* automatically login when anonymous user detected [Cory Bennett] [[21add54](https://github.com/Netflix-Skunkworks/go-jira/commit/21add54)]
* refactor trivial objects in favor of arguments to static functions [Cory Bennett] [[1f345ce](https://github.com/Netflix-Skunkworks/go-jira/commit/1f345ce)]
* set JIRA_OPERATION when parsing configs.  Use figtree config types for options to make defaulting work [Cory Bennett] [[5716a7c](https://github.com/Netflix-Skunkworks/go-jira/commit/5716a7c)]
* add better handing for usage error [Cory Bennett] [[b235dcc](https://github.com/Netflix-Skunkworks/go-jira/commit/b235dcc)]
* adding `request` command, removing dead code [Cory Bennett] [[56b1c9d](https://github.com/Netflix-Skunkworks/go-jira/commit/56b1c9d)]
* adding Do required for request language [Cory Bennett] [[a1c2849](https://github.com/Netflix-Skunkworks/go-jira/commit/a1c2849)]
* add `browse` command and implement -b option for most operations [Cory Bennett] [[a91b9d5](https://github.com/Netflix-Skunkworks/go-jira/commit/a91b9d5)]
* fix IssueAssign [Cory Bennett] [[f32cc70](https://github.com/Netflix-Skunkworks/go-jira/commit/f32cc70)]
* merge in update for upstream changes [#104](https://github.com/Netflix-Skunkworks/go-jira/issues/104) [Cory Bennett] [[19d8686](https://github.com/Netflix-Skunkworks/go-jira/commit/19d8686)]
* add `export-templates` command [Cory Bennett] [[abaad56](https://github.com/Netflix-Skunkworks/go-jira/commit/abaad56)]
* add `issuetypes` command [Cory Bennett] [[da39323](https://github.com/Netflix-Skunkworks/go-jira/commit/da39323)]
* add `components` command [Cory Bennett] [[0bd3ca2](https://github.com/Netflix-Skunkworks/go-jira/commit/0bd3ca2)]
* add `component add` command [Cory Bennett] [[cc90610](https://github.com/Netflix-Skunkworks/go-jira/commit/cc90610)]
* add `take`, `unassign` and `assign|give` commands [Cory Bennett] [[959524a](https://github.com/Netflix-Skunkworks/go-jira/commit/959524a)]
* adding `labels [add|set|remove]` commands [Cory Bennett] [[9161861](https://github.com/Netflix-Skunkworks/go-jira/commit/9161861)]
* add `comment` command [Cory Bennett] [[f0b08c5](https://github.com/Netflix-Skunkworks/go-jira/commit/f0b08c5)]
* add `watch` command [Cory Bennett] [[ec0ac3c](https://github.com/Netflix-Skunkworks/go-jira/commit/ec0ac3c)]
* add `rank ISSUE after|before ISSUE` command [Cory Bennett] [[8b863d2](https://github.com/Netflix-Skunkworks/go-jira/commit/8b863d2)]
* add `vote` command [Cory Bennett] [[a08c92f](https://github.com/Netflix-Skunkworks/go-jira/commit/a08c92f)]
* add `issuelinktypes` command [Cory Bennett] [[37f81a4](https://github.com/Netflix-Skunkworks/go-jira/commit/37f81a4)]
* add `issuelink` command [Cory Bennett] [[aacc9f4](https://github.com/Netflix-Skunkworks/go-jira/commit/aacc9f4)]
* fix closing duplicate issue on `dup` command [Cory Bennett] [[fc696c3](https://github.com/Netflix-Skunkworks/go-jira/commit/fc696c3)]
* rewrite checkpoint [Cory Bennett] [[36632a5](https://github.com/Netflix-Skunkworks/go-jira/commit/36632a5)]

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
