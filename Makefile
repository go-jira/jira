PLATFORMS= \
	freebsd/amd64 \
	linux/386 \
	linux/amd64 \
	windows/386 \
	windows/amd64 \
	darwin/amd64 \
	$(NULL)

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

GOPATH ?= $(CWD)
export GOPATH

DIST=$(CWD)$(SEP)dist

GOBIN ?= $(CWD)

CURVER ?= $(patsubst v%,%,$(shell [ -d .git ] && git describe --abbrev=0 --tags || grep ^\#\# CHANGELOG.md | awk '{print $$2; exit}'))
LDFLAGS := -w

PACKAGE=github.com/Netflix-Skunkworks/go-jira

# use 'make debug' and you can get a debuggable golang binary
# see https://golang.org/doc/gdb
# note on mac's you will need to codesign the gdb binary before you can use it:
#     codesign -fs gdb-cert /usr/local/bin/gdb
ifneq ($(DEBUG),)
	GOBUILD=go build -ldflags "-s" -gcflags "-N -l"
else
	GOBUILD=go build -ldflags "$(LDFLAGS) -s"
endif

build: $(GOPATH)/src/$(PACKAGE)
	cd $(GOPATH)/src/$(PACKAGE) && $(GOBUILD) -o $(BIN) main.go

debug:
	$(MAKE) DEBUG=1 build

$(GOPATH)/src/%:
	mkdir -p $(@D)
	test -L $@ || ln -sf ../../.. $@
	glide install -v

vet:
	@go vet *.go lib/*.go data/*.go

lint:
	@go get github.com/golang/lint/golint
	@$(GOPATH)/bin/golint .
	@$(GOPATH)/bin/golint ./data
	@$(GOPATH)/bin/golint ./lib

test: $(GOPATH)/src/$(PACKAGE)
	cd $(GOPATH)/src/$(SUBPACKAGE) && go test -v

cross-setup:
	for p in $(PLATFORMS); do \
        echo Building for $$p"; \
		cd $(GOROOT)/src && sudo GOROOT_BOOTSTRAP=$(GOROOT) GOOS=$${p/-*/} GOARCH=$${p/*-/} bash ./make.bash --no-clean; \
   done

all: $(GOPATH)/src/$(PACKAGE)
	docker pull karalabe/xgo-latest
	rm -rf dist
	mkdir -p dist
	docker run --rm -e EXT_GOPATH=/gopath -v $(GOPATH):/gopath -e TARGETS="$(PLATFORMS)" -v $$(pwd)/dist:/build karalabe/xgo-latest $(PACKAGE)
	cd $(DIST) && for x in go-jira-*; do mv $$x $$(echo $$x | cut -c 4-); done

fmt:
	gofmt -s -w main.go lib/*.go data/*.go

install:
	${MAKE} GOBIN=$(shell echo ~)/bin build

NEWVER ?= $(shell echo $(CURVER) | awk -F. '{print $$1"."$$2"."$$3+1}')
TODAY  := $(shell date +%Y-%m-%d)

changes:
	@git log --pretty=format:"* %s [%cn] [%h]" --no-merges ^v$(CURVER) HEAD *.go lib/*.go data/*.go | grep -vE 'gofmt|go fmt'

update-changelog:
	@echo "# Changelog" > CHANGELOG.md.new; \
	echo >> CHANGELOG.md.new; \
	echo "## $(NEWVER) - $(TODAY)" >> CHANGELOG.md.new; \
	echo >> CHANGELOG.md.new; \
	$(MAKE) --no-print-directory --silent changes | \
	perl -pe 's{\[([a-f0-9]+)\]}{[[$$1](https://github.com/Netflix-Skunkworks/go-jira/commit/$$1)]}g' | \
	perl -pe 's{\#(\d+)}{[#$$1](https://github.com/Netflix-Skunkworks/go-jira/issues/$$1)}g' >> CHANGELOG.md.new; \
	tail -n +2 CHANGELOG.md >> CHANGELOG.md.new; \
	mv CHANGELOG.md.new CHANGELOG.md; \
	git commit -m "Updated Changelog" CHANGELOG.md; \
	perl -pi -e 's{VERSION = "$(CURVER)"}{VERSION = "$(NEWVER)"}' lib/cli.go; \
	git commit -m "version bump" lib/cli.go; \
	git tag v$(NEWVER)

version:
	@echo $(CURVER)

clean:
	rm -rf pkg dist bin src ./$(NAME)
