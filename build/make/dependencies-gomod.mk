.PHONY: dependencies
dependencies: vendor

vendor: go.mod go.sum
	@echo "Installing dependencies using go modules..."
	${GO_CALL} mod vendor
