# Changelog

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
