CHECKSUM=$(TARGET_DIR)/$(ARTIFACT_ID).sha256sum

.PHONY: checksum
checksum: $(CHECKSUM)
# we have to depend on target dir, because we want to rebuild the checksum
# if one of the artefacts was changed
$(CHECKSUM): $(TARGET_DIR)
	@echo "Generating Checksums"
	@cd $(TARGET_DIR); find . -maxdepth 1 -not -type d | egrep -v ".(sha256sum|asc)$$" | xargs shasum -a 256 > $$(basename $@)

SIGNATURE=$(CHECKSUM).asc

.PHONY: signature
signature: $(SIGNATURE)
$(SIGNATURE): $(CHECKSUM)
	@echo "Generating Signature"
	@gpg --batch --yes --detach-sign --armor -o $@ $<
