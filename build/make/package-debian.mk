# This Makefile holds all targets for building a debian package
# For deployment of the deb package include the deploy-debian.mk!

PREPARE_PACKAGE?=prepare-package
DEBIAN_PACKAGE_FORMAT_VERSION="2.0"
CONFFILES_FILE="$(DEBIAN_CONTENT_DIR)/control/conffiles"
CONFFILES_FILE_TMP="$(DEBIAN_CONTENT_DIR)/conffiles_"
DEBSRC:=$(shell find "${WORKDIR}/deb" -type f)

.PHONY: package
package: debian-with-binary

.PHONY: debian
debian: $(DEBIAN_PACKAGE)

.PHONY: debian-with-binary
debian-with-binary: $(BINARY) $(DEBIAN_PACKAGE)

.PHONY: prepare-package
prepare-package:
	@echo "Using default prepare-package target. To write your own, define a target and specify it in the PREPARE_PACKAGE variable, before the package-debian.mk import"

$(DEBIAN_BUILD_DIR):
	@mkdir $@

$(DEBIAN_BUILD_DIR)/debian-binary: $(DEBIAN_BUILD_DIR)
	@echo $(DEBIAN_PACKAGE_FORMAT_VERSION) > $@

$(DEBIAN_CONTENT_DIR)/control:
	@install -p -m 0755 -d $@

$(DEBIAN_CONTENT_DIR)/data:
	@install -p -m 0755 -d $@

$(DEBIAN_PACKAGE): $(TARGET_DIR) $(DEBIAN_CONTENT_DIR)/control $(DEBIAN_CONTENT_DIR)/data $(DEBIAN_BUILD_DIR)/debian-binary $(PREPARE_PACKAGE) $(DEBSRC)
	@echo "Creating .deb package..."

# populate control directory
	@sed -e "s/^Version:.*/Version: $(VERSION)/g" deb/DEBIAN/control > $(DEBIAN_CONTENT_DIR)/_control
	@install -p -m 0644 $(DEBIAN_CONTENT_DIR)/_control $(DEBIAN_CONTENT_DIR)/control/control

# populate data directory
	@for dir in $$(find deb -mindepth 1 -not -name "DEBIAN" -a -type d |sed s@"^deb/"@"$(DEBIAN_CONTENT_DIR)/data/"@) ; do \
		install -m 0755 -d $${dir} ; \
	done

	@for file in $$(find deb -mindepth 1 -type f | grep -v "DEBIAN") ; do \
		cp $${file} $(DEBIAN_CONTENT_DIR)/data/$${file#deb/} ; \
	done

# Copy binary to data/usr/sbin, if it exists
	@if [ -f $(BINARY) ]; then \
		echo "Copying binary to $(DEBIAN_CONTENT_DIR)/data/usr/sbin"; \
		install -p -m 0755 -d $(DEBIAN_CONTENT_DIR)/data/usr/sbin; \
		install -p -m 0755 $(BINARY) $(DEBIAN_CONTENT_DIR)/data/usr/sbin/; \
	fi

# create conffiles file which help to deal with config change
# in order to successfully add the conffiles file to the archive it must exist, even empty
	@touch $(CONFFILES_FILE_TMP)
	@for file in $$(find $(DEBIAN_CONTENT_DIR)/data/etc -mindepth 1 -type f | grep -v "DEBIAN") ; do \
		echo $$file | sed s@'.*\(/etc/\)@\1'@ >> $(CONFFILES_FILE_TMP) ; \
	done
	@install -p -m 0644 $(CONFFILES_FILE_TMP) $(CONFFILES_FILE)
	@rm $(CONFFILES_FILE_TMP)

# create control.tar.gz
	@tar cvfz $(DEBIAN_CONTENT_DIR)/control.tar.gz -C $(DEBIAN_CONTENT_DIR)/control $(TAR_ARGS) .

# create data.tar.gz
	@tar cvfz $(DEBIAN_CONTENT_DIR)/data.tar.gz -C $(DEBIAN_CONTENT_DIR)/data $(TAR_ARGS) .

# create package
	@ar roc $@ $(DEBIAN_BUILD_DIR)/debian-binary $(DEBIAN_CONTENT_DIR)/control.tar.gz $(DEBIAN_CONTENT_DIR)/data.tar.gz
	@echo "... deb package can be found at $@"

APTLY:=curl --silent --show-error --fail -u "${APT_API_USERNAME}":"${APT_API_PASSWORD}"
