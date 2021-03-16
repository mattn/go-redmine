STATIC_ANALYSIS_DIR=$(TARGET_DIR)/static-analysis
GOIMAGE?=golang
GOTAG?=1.14.13
CUSTOM_GO_MOUNT?=-v /tmp:/tmp

REVIEW_DOG=$(TMP_DIR)/bin/reviewdog
LINT=$(TMP_DIR)/bin/golangci-lint
# ignore tests and mocks
LINTFLAGS=--tests=false --skip-files="^.*_mock.go$$" --skip-files="^.*/mock.*.go$$"

.PHONY: static-analysis
static-analysis: static-analysis-$(ENVIRONMENT)

.PHONY: static-analysis-ci
static-analysis-ci:
	@make $(STATIC_ANALYSIS_DIR)/static-analysis-cs.log $(STATIC_ANALYSIS_DIR)/static-analysis.log static-analysis-ci-report-pr

static-analysis-ci-report-pr: $(REVIEW_DOG)
	@if [ X"$(CI_PULL_REQUEST)" != X"" -a X"$(CI_PULL_REQUEST)" != X"null" ] ; then \
  		cat $(STATIC_ANALYSIS_DIR)/static-analysis-cs.log | CI_COMMIT=$(COMMIT_ID) $(REVIEW_DOG) -f=checkstyle -reporter="github-pr-review"; \
  	fi

.PHONY: static-analysis-local
static-analysis-local: $(PASSWD) $(ETCGROUP) $(HOME_DIR)
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
			make $(STATIC_ANALYSIS_DIR)/static-analysis-cs.log $(STATIC_ANALYSIS_DIR)/static-analysis.log static-analysis-ci-report-local

$(STATIC_ANALYSIS_DIR)/static-analysis.log: $(STATIC_ANALYSIS_DIR)
	@echo ""
	@echo "complete static analysis:"
	@echo ""
	@$(LINT) $(LINTFLAGS) run ./... | tee $@

$(STATIC_ANALYSIS_DIR)/static-analysis-cs.log: $(STATIC_ANALYSIS_DIR)
	@echo "run static analysis with export to checkstyle format"
	@$(LINT) $(LINTFLAGS) run --out-format=checkstyle ./... > $@ | true

$(STATIC_ANALYSIS_DIR): $(LINT)
	@mkdir -p $(STATIC_ANALYSIS_DIR)

static-analysis-ci-report-local: $(STATIC_ANALYSIS_DIR)/static-analysis-cs.log $(REVIEW_DOG)
	@echo ""
	@echo "differences to develop branch:"
	@echo ""
	@cat $(STATIC_ANALYSIS_DIR)/static-analysis-cs.log | $(REVIEW_DOG) -f checkstyle -diff "git diff develop"

$(LINT): $(TMP_DIR)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(TMP_DIR)/bin v1.33.0

$(REVIEW_DOG): $(TMP_DIR)
	@curl -sfL https://raw.githubusercontent.com/reviewdog/reviewdog/master/install.sh| sh -s -- -b $(TMP_DIR)/bin