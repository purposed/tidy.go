NAME := tidy
PKG := github.com/purposed/$(NAME)

CGO_ENABLED := 0

BUILDTAGS :=

include root.mk

.PHONY: prebuild
prebuild:

.PHONY: extra_validation
extra_validation:
