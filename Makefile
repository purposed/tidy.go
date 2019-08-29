# Makefile magic from Jessie Frazelle:
# https://github.com/jessfraz

NAME := tidy
PKG := github.com/purposed/$(NAME)
PREFIX ?=$(shell pwd)
BUILDDIR := ${PREFIX}/bin

VERSION := $(shell cat VERSION.txt)
GITCOMMIT := $(shell git rev-parse --short HEAD)
GITUNTRACKEDCHANGES := $(shell git status --porcelain --untracked-files=no)
ifneq ($(GITUNTRACKEDCHANGES),)
	GITCOMMIT := $(GITCOMMIT)-dirty
endif
ifeq ($(GITCOMMIT),)
    GITCOMMIT := ${GITHUB_SHA}
endif

GOOSARCHES = linux/amd64 darwin/amd64 windows/amd64
BUILDTAGS :=

CTIMEVAR=-X $(PKG)/version.GITCOMMIT=$(GITCOMMIT) -X $(PKG)/version.VERSION=$(VERSION)
GO_LDFLAGS=-ldflags "-w $(CTIMEVAR)"
GO_LDFLAGS_STATIC=-ldflags "-w $(CTIMEVAR) -extldflags -static"

GO := go

define buildarch
mkdir -p $(BUILDDIR)/$(1)-$(2)
GOOS=$(1) GOARCH=$(2) CGO_ENABLED=0 $(GO) build \
	-o $(BUILDDIR)/$(1)-$(2)/$(3)/$(4) \
	-a -tags "$(BUILDTAGS) static_build netgo" \
	-installsuffix netgo ${GO_LDFLAGS_STATIC} $(5);
endef

define build
CGO_ENABLED=0 $(GO) build \
	-o $(1)/$(2) \
	-a -tags "$(BUILDTAGS) static_build netgo" \
	-installsuffix netgo ${GO_LDFLAGS_STATIC} $(3);
endef

.PHONY: clean
clean: ## Cleanup any build binaries or packages.
	@echo "+ $@"
	$(RM) -r $(BUILDDIR)

.PHONY: build
build: *.go VERSION.txt
	@echo "+ $@"
	$(call build,$(BUILDDIR),$(NAME),.)

.PHONY: install
install: build
	@echo "+ @"
	(cp $(BUILDDIR)/$(NAME) $(HOME)/bin)

.PHONY: release
release: *.go VERSION.txt
	@echo "+ $@"
	$(GO) generate ./...
	$(foreach GOOSARCH,$(GOOSARCHES), \
		$(call buildarch,$(subst /,,$(dir $(GOOSARCH))),$(notdir $(GOOSARCH)),.,$(NAME),.)\
	)

.PHONY: bump-version
BUMP := patch
bump-version: ## Bump the version in the version file. Set BUMP to [ patch | major | minor ].
	@$(GO) get -u github.com/jessfraz/junk/sembump # update sembump tool
	$(eval NEW_VERSION = $(shell sembump --kind $(BUMP) $(VERSION)))
	@echo "Bumping VERSION.txt from $(VERSION) to $(NEW_VERSION)"
	echo $(NEW_VERSION) > VERSION.txt
	@echo "Updating links to download binaries in README.md"
	sed -i s/$(VERSION)/$(NEW_VERSION)/g README.md
	git add VERSION.txt README.md
	git commit -vsam "Bump version to $(NEW_VERSION)"
	@echo "Run make tag to create and push the tag for new version $(NEW_VERSION)"

.PHONY: tag
tag: ## Create a new git tag to prepare to build a release.
	git tag -a $(VERSION) -m "$(VERSION)"
	@echo "Run git push origin $(VERSION) to push your new tag to GitHub and trigger a travis build."

.PHONY: test
test:
	@echo "+ $@"
	go test -cover -race ./...
