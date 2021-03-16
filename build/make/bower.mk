BOWER_JSON=$(WORKDIR)/bower.json

.PHONY: bower-install
bower-install: $(BOWER_TARGET)

ifeq ($(ENVIRONMENT), ci)

$(BOWER_TARGET): $(BOWER_JSON) $(YARN_TARGET)
	@echo "Yarn run bower on CI server"
	@yarn run bower

else

$(BOWER_TARGET): $(BOWER_JSON) $(PASSWD) $(YARN_TARGET)
	@echo "Executing bower..."
	@docker run --rm \
	  -e HOME=/tmp \
	  -u "$(UID_NR):$(GID_NR)" \
	  -v $(PASSWD):/etc/passwd:ro \
	  -v $(WORKDIR):$(WORKDIR) \
	  -w $(WORKDIR) \
	  node:8 \
	  yarn run bower
	@touch $@

endif
