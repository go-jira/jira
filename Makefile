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

CURVER := $(shell grep '\#\#' CHANGELOG.md  | awk '{print $$2; exit}')
NEWVER := $(shell awk -F'"' '/docopt.Parse/{print $$2}' jira/main.go)
TODAY  := $(shell date +%Y-%m-%d)

changes:
	@git log --pretty=format:"* %s [%cn] [%h]" --no-merges ^$(CURVER) HEAD jira | grep -v gofmt | grep -v "bump version"

update-changelog: 
	@echo "# Changelog" > CHANGELOG.md.new; \
	echo >> CHANGELOG.md.new; \
	echo "## $(NEWVER) - $(TODAY)" >> CHANGELOG.md.new; \
	echo >> CHANGELOG.md.new; \
	$(MAKE) changes | \
	perl -pe 's{\[([a-f0-9]+)\]}{[[$$1](https://github.com/Netflix-Skunkworks/go-jira/commit/$$1)]}g' | \
	perl -pe 's{\#(\d+)}{[#$$1](https://github.com/Netflix-Skunkworks/go-jira/issues/$$1)}g' >> CHANGELOG.md.new; \
	tail +2 CHANGELOG.md >> CHANGELOG.md.new; \
	mv CHANGELOG.md.new CHANGELOG.md

# https://github.com/Netflix-Skunkworks/go-jira/commit/d5330fd
# [#1349](https://github.com/bower/bower/issues/1349)
