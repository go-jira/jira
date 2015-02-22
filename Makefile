PLATFORMS= \
	freebsd-386 \
	freebsd-amd64 \
	freebsd-arm \
	linux-386 \
	linux-amd64 \
	linux-arm \
	openbsd-386 \
	openbsd-amd64 \
	windows-386 \
	windows-amd64 \
	darwin-386 \
	darwin-amd64 \
	$(NULL)

DIST=$(shell pwd)/dist
export GOPATH=$(shell pwd)

build:
	cd src/github.com/Netflix-Skunkworks/go-jira/jira; \
	go get -v

all:
	mkdir -p $(DIST); \
	cd src/github.com/Netflix-Skunkworks/go-jira/jira; \
	go get -d; \
	for p in $(PLATFORMS); do \
        echo "Building for $$p"; \
        GOOS=$${p/-*/} GOARCH=$${p/*-/} go build -v -o $(DIST)/jira-$$p; \
   done

fmt:
	gofmt -s -w jira
