#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira --user admin"

PLAN 14

# clean out any old containers
docker rm -f go-jira-test

# start newt jira service
RUNS docker run --detach --name go-jira-test --publish 8080:8080 go-jira-test:latest

# wait a few seconds for it to bind to port 8080
RUNS sleep 10

# wait for healthchecks to pass, curl will retry 60 times over 5 min waiting
RUNS curl -q -L --retry 360 --retry-delay 1 -f -s "http://localhost:8080/rest/api/2/serverInfo?doHealthCheck=1"

# login to jira as admin user
echo "admin123" | RUNS $jira login

# create gojira user
RUNS $jira req -M POST /rest/api/2/user '{"name":"gojira","password":"gojira123","emailAddress":"gojira@example.com","displayName":"Go Jira"}'

# create mojira user (need secondary user for voting)
RUNS $jira req -M POST /rest/api/2/user '{"name":"mojira","password":"mojira123","emailAddress":"mojira@example.com","displayName":"Mo Jira"}'

# create SCRUM softwareproject
RUNS $jira req -M POST /rest/api/2/project '{"key":"SCRUM","name":"Scrum","projectTypeKey":"software","projectTemplateKey":"com.pyxis.greenhopper.jira:gh-scrum-template","lead":"gojira"}'

# create KANBAN software project
RUNS $jira req -M POST /rest/api/2/project '{"key":"KANBAN","name":"Kanban","projectTypeKey":"software","projectTemplateKey":"com.pyxis.greenhopper.jira:gh-kanban-template","lead":"gojira"}'

# create BAISC software project
RUNS $jira req -M POST /rest/api/2/project '{"key":"BASIC","name":"Basic","projectTypeKey":"software","projectTemplateKey":"com.pyxis.greenhopper.jira:basic-software-development-template","lead":"gojira"}'

# create PROJECT business project
RUNS $jira req -M POST /rest/api/2/project '{"key":"PROJECT","name":"Project","projectTypeKey":"business","projectTemplateKey":"com.atlassian.jira-core-project-templates:jira-core-project-management","lead":"gojira"}'

# create PROCESS business project
RUNS $jira req -M POST /rest/api/2/project '{"key":"PROCESS","name":"Process","projectTypeKey":"business","projectTemplateKey":"com.atlassian.jira-core-project-templates:jira-core-process-management","lead":"gojira"}'

# create TASK business project
RUNS $jira req -M POST /rest/api/2/project '{"key":"TASK","name":"Task","projectTypeKey":"business","projectTemplateKey":"com.atlassian.jira-core-project-templates:jira-core-task-management","lead":"gojira"}'

RUNS $jira logout

# export new templates so we are always using whatever is latest
# and not whatever is in the test-runners homedir
RUNS $jira export-templates -d .jira.d/templates
