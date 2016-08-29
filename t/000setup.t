#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira="../jira --user admin"

PLAN 14

# clean out any old containers
docker rm -f go-jira-test

RUNS docker build . -t go-jira-test

mkdir -p $(pwd)/.maven-cache

# start newt jira service, cache the users m2 directory to make startup faster
RUNS docker run --detach -v $(pwd)/.maven-cache:/root/.m2/repository --name go-jira-test --publish 8080:8080 go-jira-test:latest

# wait for docker service to get started
RUNS sleep 5

echo "# Waiting for jira service to be listening on port 8080"
docker exec -i go-jira-test tail -f screenlog.0

#| grep -m 1 'jira started successfully' | sed 's/^/# /'

# wait for healthchecks to pass, curl will retry 900 times over 15 min waiting
RUNS curl -q -L --retry 900 --retry-delay 1 -f -s "http://localhost:8080/rest/api/2/serverInfo?doHealthCheck=1"

# login to jira as admin user
echo "admin" | RUNS $jira login

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
