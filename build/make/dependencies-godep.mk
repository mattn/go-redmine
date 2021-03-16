GODEP=$(GOPATH)/bin/dep

$(GODEP):
	@curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

vendor: $(GODEP) Gopkg.toml Gopkg.lock
	@echo "Installing dependencies using go dep..."
	@dep ensure

dependencies: vendor
