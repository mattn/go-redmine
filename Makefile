ARTIFACT_ID := go-redmine
VERSION := 0.1.0

GOTAG=1.14.13
CUSTOM_GO_MOUNT=-v $(WORKDIR)/resources/compileHeaders/usr/include/btrfs:/usr/include/btrfs
# overwrite ADDITIONAL_LDFLAGS to disable static compilation
# this should fix https://github.com/golang/go/issues/13470
ADDITIONAL_LDFLAGS=""
MAKEFILES_VERSION=4.3.0
.DEFAULT_GOAL:=default

default: compile signature

ADDITIONAL_CLEAN=clean_add
clean_add:
	rm -rf $(BIN) goxz

include build/make/variables.mk
PACKAGES_FOR_INTEGRATION_TEST=github.com/cloudogu/cesapp/v2/tasks github.com/cloudogu/cesapp/v2/registry github.com/cloudogu/cesapp/v2/containers

include build/make/info.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
include build/make/self-update.mk

BIN=godmine

CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-s -w -X main.Version=$(VERSION)"


GOBIN ?= $(shell go env GOPATH)/bin
export GO111MODULE=on

.PHONY: all
all: compile

.PHONY: install
install:
	go install -ldflags=$(BUILD_LDFLAGS) ./...

.PHONY: cross
cross: $(GOBIN)/goxz
	goxz -n $(ARTIFACT_ID) -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) ./cmd/$(BIN)

$(GOBIN)/goxz:
	go get github.com/Songmu/goxz/cmd/goxz

.PHONY: test
test: build
	go test -v ./...

.PHONY: lint
lint: $(GOBIN)/golint
	go vet ./...
	golint -set_exit_status ./...

$(GOBIN)/golint:
	cd && go get golang.org/x/lint/golint

.PHONY: upload
upload: $(GOBIN)/ghr
	ghr "v$(VERSION)" goxz

$(GOBIN)/ghr:
	go get github.com/tcnksm/ghr
