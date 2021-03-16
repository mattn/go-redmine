YARN_LOCK=$(WORKDIR)/yarn.lock

.PHONY: yarn-install
yarn-install: $(YARN_TARGET)

ifeq ($(ENVIRONMENT), ci)

$(YARN_TARGET): $(YARN_LOCK)
	@echo "Yarn install on CI server"
	@yarn install

else

$(YARN_TARGET): $(YARN_LOCK) $(PASSWD)
	@echo "Executing yarn..."
	@docker run --rm \
	  -u "$(UID_NR):$(GID_NR)" \
	  -v $(PASSWD):/etc/passwd:ro \
	  -v $(WORKDIR):$(WORKDIR) \
	  -w $(WORKDIR) \
	  node:8 \
	  yarn install
	@touch $@

endif
