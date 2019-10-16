# Makefile magic mostly from Jessie Frazelle:
# https://github.com/jessfraz
SHELL := /bin/bash

PREFIX?=$(shell pwd)

BUILDDIR := ${PREFIX}/bin/debug
RELEASEDIR := ${PREFIX}/bin/release

VERSION := $(shell cat VERSION.txt)
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif
ifeq ($(GITCOMMIT),)
    GITCOMMIT := ${GITHUB_SHA}
endif

CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

# Setup the go compiler
GO := go

GOOSARCHES = $(shell cat .goosarch)

define buildrelease
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
	 -o $(RELEASEDIR)/$(NAME)-$(1)-$(2) \
	 -a -tags "$(BUILDTAGS) static_build netgo" \
	 -installsuffix netgo ${GO_LDFLAGS_STATIC} $(3);
endef

define build
CGO_ENABLED=$(CGO_ENABLED) $(GO) build \
	-o $(BUILDDIR)/$(NAME) \
	-a -tags "$(BUILDTAGS) static_build netgo" \
	-installsuffix netgo ${GO_LDFLAGS_STATIC} $(1);
endef

define crossrelease
	$(foreach GOOSARCH,$(GOOSARCHES), $(call buildrelease,$(subst /,,$(dir $(GOOSARCH))),$(notdir $(GOOSARCH)),$(1)))
endef


define install
	$(GO) install -a -tags "$(BUILDTAGS)" ${GO_LDFLAGS} $(1)
endef


validate: clean fmt lint test vet extra_validation


.PHONY: release
release: cmd/$(NAME) VERSION.txt prebuild validate
	@echo "+ $@"
	$(call crossrelease,./cmd/$(NAME))


.PHONY: build
build: cmd/$(NAME) VERSION.txt prebuild
	@echo "+$@"
	$(call build,./cmd/$(NAME))


.PHONY: fmt
fmt:
	@echo "+ $@"
	@if [[ ! -z "$(shell gofmt -s -l . | grep -v '.pb.go:' | grep -v '.twirp.go:' | grep -v vendor | tee /dev/stderr)" ]]; then \
		exit 1; \
	fi


.PHONY: lint
lint: ## Verifies `golint` passes.
	@echo "+ $@"
	@if [[ ! -z "$(shell golint ./... | grep -v '.pb.go:' | grep -v '.twirp.go:' | grep -v vendor | tee /dev/stderr)" ]]; then \
		exit 1; \
	fi


.PHONY: test
test: prebuild
	@echo "+ $@"
	@$(GO) test -v -race -cover -tags "$(BUILDTAGS) cgo" $(shell $(GO) list ./... | grep -v vendor)


.PHONY: vet
vet: ## Verifies `go vet` passes.
	@echo "+ $@"
	@if [[ ! -z "$(shell $(GO) vet $(shell $(GO) list ./... | grep -v vendor) | tee /dev/stderr)" ]]; then \
		exit 1; \
	fi


.PHONY: clean
clean:
	@echo "+ $@"
	$(RM) -r $(BUILDDIR)
	$(RM) -r $(RELEASEDIR)


.PHONY: install
install: prebuild validate
	@echo "+ $@"
	$(call install,./cmd/$(NAME))

