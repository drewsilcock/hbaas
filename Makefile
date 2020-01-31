-include .env

SHELL := /bin/bash

PROJECTNAME := hbaas-server

VERSION?=$(shell git describe --tags --always)
BUILDTIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go related variables.
GOBASE := $(shell pwd)
GOPATH := $(GOBASE)/.go-pkg:$(GOBASE)
GOBIN := $(GOBASE)/.go-bin
GOFILES := $(wildcard *.go)

# Directory in which to put the output executable.
OUTBINDIR?=$(GOBASE)

PACKAGENAME=$(shell go list)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X '$(PACKAGENAME)/version.Version=$(VERSION)' -X '$(PACKAGENAME)/version.BuildTime=$(BUILDTIME)'"

LINUXBUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

IS_INTERACTIVE:=$(shell [ -t 0 ] && echo 1)

ifdef IS_INTERACTIVE
LOG_INFO := $(shell tput setaf 12)
LOG_ERROR := $(shell tput setaf 9)
LOG_END := $(shell tput sgr0)
endif

define log
echo -e "$(LOG_INFO)⇛ $(1)$(LOG_END)"
endef

define log-error
echo -e "$(LOG_ERROR)⇛ $(1)$(LOG_END)"
endef

default: help

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get-build-deps go-mod-download

## start: Start in development mode. Auto-starts when code changes.
start: clean compile start-server

start-server:
	@test -e $(OUTBINDIR)/$(PROJECTNAME) || (\
	    $(call log-error,Unable to find $(PROJECTNAME) executable.) \
	    && false\
	)
	@$(call log,Starting $(PROJECTNAME)...)
	@-$(OUTBINDIR)/$(PROJECTNAME) run-server

## run: Run server with specified args, e.g. `make run args=version`
run:
	@$(OUTBINDIR)/$(PROJECTNAME) $(args)

## compile: Compile the binary.
compile: go-get-build-deps go-mod-download code-gen go-build

## compile-linux: Compile the binary for Linux.
compile-linux: go-get-build-deps go-mod-download code-gen go-build-linux

## exec: Run given command, wrapped with custom GOPATH. e.g; make exec run="go test ./..."
exec:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(run)

## build-image: Build Docker image for project.
build-image:
	@$(call log,Building Docker image...)
	docker build \
	    --build-arg version=$(VERSION) \
	    --build-arg build=$(BUILD) \
	    -t $(PROJECTNAME):latest . || \
	(\
	    $(call log-error,Unable to build Docker image.) \
	    && false \
	)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(OUTBINDIR)/$(PROJECTNAME) 2> /dev/null
	@-$(MAKE) go-clean
	@-rm ./**/bindata.go 2> /dev/null
	@-rm ./**/gen-*.go 2> /dev/null

code-gen:
	@$(call log,Running pre-build code generation...)
	@cd migrations && $(GOBIN)/go-bindata -pkg migrations . || (\
	    $(call log-error,Unable to generate binary migration data.) \
	    && false \
	)
	@cd seeddata && $(GOBIN)/go-bindata -pkg seeddata . || (\
	    $(call log-error,Unable to generate binary database seed data.) \
	    && false \
	)
	# If there's no docs package, the docs generation fails. Because of this
	# bug, we need to put a stub package in so that Swag can replace it.
	@mkdir -p docs
	@echo "package docs" > docs/docs.go
	@$(GOBIN)/swag init --parseDependency || (\
	    $(call log-error,Unable to generate Swagger docs.) \
	    && false \
	)

go-build:
	@$(call log,Building binary...)

	# Uncomment line below for debugging build arguments.
	#@echo Building $(PROJECTNAME): GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(OUTBINDIR)/$(PROJECTNAME) $(GOFILES)

	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(OUTBINDIR)/$(PROJECTNAME) $(GOFILES) || (\
	    $(call log-error,Failed to build $(PROJECTNAME).) \
	    && false \
	)

go-build-linux:
	@$(call log,Building binary...)

	# Uncomment line below for debugging build arguments.
	#@echo Linux build command: "GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(LINUXBUILD) go build $(LDFLAGS) -o $(OUTBINDIR)/$(PROJECTNAME) $(GOFILES)"

	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) $(LINUXBUILD) go build $(LDFLAGS) -o $(OUTBINDIR)/$(PROJECTNAME) $(GOFILES) || (\
	    $(call log-error,Failed to build $(PROJECTNAME).) \
	    && false \
	)

go-generate:
	@$(call log,Generating dependency files...)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

go-get-build-deps:
	@$(call log,Checking if there are any missing build-time dependencies...)
	@make go-get bin=go-bindata pkg=github.com/kevinburke/go-bindata/...
	@make go-get bin=swag pkg=github.com/swaggo/swag/cmd/swag

go-get:
	@test -e $(GOBIN)/$(bin) || \
	    GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(pkg)

go-mod-download:
	@$(call log,Checking if there are any missing modules...)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go mod download

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@$(call log,Cleaning build cache)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

.PHONY: help
all: help
help: Makefile
	@echo
	@echo "Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
