# Changelog

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
