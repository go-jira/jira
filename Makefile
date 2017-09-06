PLATFORMS= \
	freebsd/amd64 \
	linux/386 \
	linux/amd64 \
	windows/386 \
	windows/amd64 \
	darwin/amd64 \
	$(NULL)

	# freebsd-386 \
	# freebsd-arm \
	# linux-arm \
	# openbsd-386 \
	# openbsd-amd64 \
	# darwin-386

NAME=jira

OS=$(shell uname -s)
ifeq ($(filter CYGWIN%,$(OS)),$(OS))
	export CWD=$(shell cygpath -wa .)
	export SEP=\\
	export CYGWIN=winsymlinks:native
	BIN ?= $(GOBIN)$(SEP)$(NAME).exe
else
	export CWD=$(shell pwd)
	export SEP=/
	BIN ?= $(GOBIN)$(SEP)$(NAME)
endif

DIST=$(CWD)$(SEP)dist

GOBIN ?= $(CWD)

CURVER ?= $(patsubst v%,%,$(shell [ -d .git ] && git describe --abbrev=0 --tags || grep ^\#\# CHANGELOG.md | awk '{print $$2; exit}'))
LDFLAGS:= -w

build:
	go build -gcflags="-e" -v -ldflags "$(LDFLAGS) -s" -o '$(BIN)' cmd/jira/main.go

vet:
	@go vet .
	@go vet ./jiracli
	@go vet ./jiracmd
	@go vet ./jiradata
	@go vet ./cmd/jira

lint:
	@go get github.com/golang/lint/golint
	@golint .
	@golint ./jiracli
	@golint ./jiracmd
	@golint ./jiradata
	@golint ./cmd/jira

all: 
	docker pull karalabe/xgo-latest
	rm -rf dist
	mkdir -p dist
	docker run --rm -e EXT_GOPATH=/gopath -v $$(pwd):/gopath/src/gopkg.in/Netflix-Skunkworks/go-jira.v1 -e TARGETS="$(PLATFORMS)" -v $$(pwd)/dist:/build karalabe/xgo-latest gopkg.in/Netflix-Skunkworks/go-jira.v1/cmd/jira

install:
	${MAKE} GOBIN=$$HOME/bin build

NEWVER ?= $(shell echo $(CURVER) | awk -F. '{print $$1"."$$2"."$$3+1}')
TODAY  := $(shell date +%Y-%m-%d)

changes:
	@git log --pretty=format:"* %s [%cn] [%h]" --no-merges ^v$(CURVER) HEAD *.go jiracli/*.go jiradata/*.go jiracmd/*.go cmd/*/*.go | grep -vE 'gofmt|go fmt'

update-changelog:
	@echo "# Changelog" > CHANGELOG.md.new; \
	echo >> CHANGELOG.md.new; \
	echo "## $(NEWVER) - $(TODAY)" >> CHANGELOG.md.new; \
	echo >> CHANGELOG.md.new; \
	$(MAKE) --no-print-directory --silent changes | \
	perl -pe 's{\[([a-f0-9]+)\]}{[[$$1](https://github.com/Netflix-Skunkworks/go-jira/commit/$$1)]}g' | \
	perl -pe 's{\#(\d+)}{[#$$1](https://github.com/Netflix-Skunkworks/go-jira/issues/$$1)}g' >> CHANGELOG.md.new; \
	tail -n +2 CHANGELOG.md >> CHANGELOG.md.new; \
    perl -pi -e 's{VERSION = "$(CURVER)"}{VERSION = "$(NEWVER)"}' jira.go; \
	mv CHANGELOG.md.new CHANGELOG.md; \
	$(NULL)

release:
	git commit -m "Updated Changelog" CHANGELOG.md; \
	git tag v$(NEWVER)
	git push --tags

version:
	@echo $(CURVER)

clean:
	rm -rf ./$(NAME)

export GNUPGHOME=$(CWD)/t/.gnupg
export PASSWORD_STORE_DIR=$(CWD)/t/.password-store
export JIRACLOUD=1

prove:
	chmod -R g-rwx,o-rwx $(GNUPGHOME)
	OSHT_VERBOSE=1 prove -v 

generate:
	cd schemas && ./fetch-schemas.py
	grep -h slipscheme jiradata/*.go | grep json | sort | uniq | awk -F\/\/ '{print $$2}' | while read cmd; do $$cmd; done
