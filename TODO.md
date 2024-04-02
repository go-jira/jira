# TODO

* https://confluence.atlassian.com/jiracore/createmeta-rest-endpoint-to-be-removed-975040986.html
* Jira 9.0 moved the parameters on api/2/issue/createmeta into path values
  Need to allow both versions to function
  Need a configuration item to indicate which url/method to call
  Need to write the new method to conform to the new api method format
  Need new jiradata type because the api endpoint changed the output json schema
* slipstream to create jiradata types
  slipscheme -stdout schema/IssueTypes.json > jiradata/IssueTypes.go

* Old API Doc: https://docs.atlassian.com/software/jira/docs/api/REST/7.2.7/#api/2/issue-getCreateIssueMeta
* New API Doc: https://docs.atlassian.com/software/jira/docs/api/REST/9.0.0/#issue-getCreateIssueMetaProjectIssueTypes
