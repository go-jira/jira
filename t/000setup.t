#!/bin/bash
eval "$(curl -q -s https://raw.githubusercontent.com/coryb/osht/master/osht.sh)"
cd $(dirname $0)
jira=../jira

PLAN 4

RUNS sh -c "docker rm -f go-jira-test || true"

RUNS docker run --detach --name go-jira-test --publish 8080:8080 go-jira-test:latest

RUNS sleep 10

RUNS curl -q -L --retry 60 --retry-delay 1 -f -s "http://localhost:8080/rest/api/2/serverInfo?doHealthCheck=1"
