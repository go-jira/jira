## Tests

The test are written using the `osht` bash testing framework.  Please read the [documentation](https://github.com/coryb/osht/blob/master/README.md) for `osht`.

## Running Test:

From the top level of the project you can run:
```
# this creates a local "jira" binary
make

# this runs the integration tests in the "t" directory
prove
```

### Running individual tests
To run a specific test you can run it directly like:
```
./100basic.t
```
There is a useful `-v` option to make the test more verbose and an `-a` option to casue the test to abort after the first failure.

The tests all require the jira service to be running from the docker container, so you will have to manually run the setup script:
```
./000setup.t
```

After than you can run the other tests over and over.  The jira service is just a test instance started for local development.  It comes with
a temporary license (I think it is 8 hours) so you will have to run the `./000setup.t` script at least once daily.

## API Documentation:
https://docs.atlassian.com/jira/REST/cloud/
https://docs.atlassian.com/jira-software/REST/cloud

## projectTempalteKey missing documentation
https://answers.atlassian.com/questions/36176301/jira-api-7.1.0-create-project

