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


cross-setup:
	for p in $(PLATFORMS); do \
        echo "Building for $$p"; \
		cd $(GOROOT)/src && sudo GOOS=$${p/-*/} GOARCH=$${p/*-/} bash ./make.bash --no-clean; \
   done

all:
	rm -rf $(DIST); \
	mkdir -p $(DIST); \
	cd src/github.com/Netflix-Skunkworks/go-jira/jira; \
	go get -d; \
	for p in $(PLATFORMS); do \
        echo "Building for $$p"; \
		GOOS=$${p/-*/} GOARCH=$${p/*-/} go build -v -ldflags -s -o $(DIST)/jira-$$p; \
   done

fmt:
	gofmt -s -w jira

install:
	export GOBIN=~/bin && ${MAKE} build

CURVER ?= $(shell git fetch --tags && git tag | tail -1)
NEWVER ?= $(shell echo $(CURVER) | awk -F. '{print $$1"."$$2"."$$3+1}')
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
	mv CHANGELOG.md.new CHANGELOG.md; \
	git commit -m "Updated Changelog" CHANGELOG.md; \
	git tag $(NEWVER)

clean:
	rm -rf pkg dist bin && find src \! -path \*/go-jira\* -delete
