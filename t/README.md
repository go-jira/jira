## Tests

The test are written using the `osht` bash testing framework.  Please read the [documentation](https://github.com/coryb/osht/blob/master/README.md) for `osht`.

## Setup
These tests assume there is a jira service running at 127.0.0.1:8080 with user "gojira" and password "gojira123".
There should also be a poweruser "admin" with password "admin123"

The test Jira was setup following the instructions [here](https://github.com/cptactionhank/docker-atlassian-jira).


### build base docker image
```
docker run --rm -i -v $(pwd):/root:ro coryb/dfpp:1.0.2 Dockerfile.pre | docker build -t go-jira-base:latest - 
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

