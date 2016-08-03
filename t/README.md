## Tests

The test are written using the `osht` bash testing framework.  Please read the [documentation](https://github.com/coryb/osht/blob/master/README.md) for `osht`.

## Setup
These tests assume there is a jira service running at 127.0.0.1:8080 with user "gojira" and password "gojira123".
There should also be a poweruser "admin" with password "admin123"

The test Jira was setup following the instructions [here](https://github.com/cptactionhank/docker-atlassian-jira).


### build base docker image
```
docker run --rm -i -v $(pwd):/root:ro coryb/dfpp Dockerfile.pre | docker build -t go-jira-base:latest - 
```

### Initialize container
```
docker run --detach --name go-jira-test --publish 8080:8080 go-jira-base:latest
```

### create admin user
```
open http://localhost:8080
```
Then follow UI workflow to create "admin" user, skip intro and project creation.

### Create gojira user
```
jira req --user admin -M POST /rest/api/2/user '{"name":"gojira","password":"gojira123","emailAddress":"gojira@example.com","displayName":"Go Jira"}'
```

### Initialize new projects
```
jira req --user admin -M POST /rest/api/2/project '{"key":"SCRUM","name":"Scrum","projectTypeKey":"software","projectTemplateKey":"com.pyxis.greenhopper.jira:gh-scrum-template","lead":"gojira"}'
jira req --user admin -M POST /rest/api/2/project '{"key":"KANBAN","name":"Kanban","projectTypeKey":"software","projectTemplateKey":"com.pyxis.greenhopper.jira:gh-kanban-template","lead":"gojira"}'
jira req --user admin -M POST /rest/api/2/project '{"key":"BASIC","name":"Basic","projectTypeKey":"software","projectTemplateKey":"com.pyxis.greenhopper.jira:basic-software-development-template","lead":"gojira"}'

jira req --user admin -M POST /rest/api/2/project '{"key":"PROJECT","name":"Project","projectTypeKey":"business","projectTemplateKey":"com.atlassian.jira-core-project-templates:jira-core-project-management","lead":"gojira"}'
jira req --user admin -M POST /rest/api/2/project '{"key":"PROCESS","name":"Process","projectTypeKey":"business","projectTemplateKey":"com.atlassian.jira-core-project-templates:jira-core-process-management","lead":"gojira"}'
jira req --user admin -M POST /rest/api/2/project '{"key":"TASK","name":"Task","projectTypeKey":"business","projectTemplateKey":"com.atlassian.jira-core-project-templates:jira-core-task-management","lead":"gojira"}'
```

### snapshot docker container
```
docker commit go-jira-test go-jira-test:latest
```

### Destroy base container

```
docker rm -f go-jira-test
```

## Running Test:

From the top level of the project you can run:
```
# this creates a local "jira" binary
make

# this runs the integration tests in the "t" directory
prove
```

## API Documentation:
https://docs.atlassian.com/jira/REST/cloud/

## projectTempalteKey missing documentation
https://answers.atlassian.com/questions/36176301/jira-api-7.1.0-create-project

