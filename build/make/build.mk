ADDITIONAL_LDFLAGS?=-extldflags -static
LDFLAGS?=-ldflags "$(ADDITIONAL_LDFLAGS) -X main.Version=$(VERSION) -X main.CommitID=$(COMMIT_ID)"
GOIMAGE?=golang
GOTAG?=1.14.13
GOOS?=linux
GOARCH?=amd64
PRE_COMPILE?=
GO_ENV_VARS?=
CUSTOM_GO_MOUNT?=-v /tmp:/tmp

.PHONY: compile
compile: $(BINARY)

compile-ci:
	@echo "Compiling (CI)..."
	make compile-generic

compile-generic:
	@echo "Compiling..."
# here is go called without mod capabilities because of error "go: error loading module requirements"
# see https://github.com/golang/go/issues/30868#issuecomment-474199640
	@$(GO_ENV_VARS) go build -mod=vendor -a -tags netgo $(LDFLAGS) -installsuffix cgo -o $(BINARY)


ifeq ($(ENVIRONMENT), ci)

$(BINARY): $(SRC) vendor $(PRE_COMPILE)
	@echo "Built on CI server"
	@make compile-generic

else

$(BINARY): $(SRC) vendor $(PASSWD) $(ETCGROUP) $(HOME_DIR) $(PRE_COMPILE)
	@echo "Building locally (in Docker)"
	@docker run --rm \
		-e GOOS=$(GOOS) \
		-e GOARCH=$(GOARCH) \
		-u "$(UID_NR):$(GID_NR)" \
		-v $(PASSWD):/etc/passwd:ro \
		-v $(ETCGROUP):/etc/group:ro \
		-v $(HOME_DIR):/home/$(USER) \
		-v $(WORKDIR):/go/src/github.com/cloudogu/$(ARTIFACT_ID) \
		$(CUSTOM_GO_MOUNT) \
		-w /go/src/github.com/cloudogu/$(ARTIFACT_ID) \
		$(GOIMAGE):$(GOTAG) \
  make compile-generic

endif
