PKG_NAME     ?= sakuracloud
WEBSITE_REPO  = github.com/hashicorp/terraform-website

AUTHOR          ?="terraform-provider-sakuracloud authors"
COPYRIGHT_YEAR  ?="2016-2020"
COPYRIGHT_FILES ?=$$(find . \( -name "*.dockerfile" -or -name "*.go" -or -name "*.sh" -or -name "*.pl" -or -name "*.bash" \) -print | grep -v "/vendor/")

UNIT_TEST_UA ?= (Unit Test)
ACC_TEST_UA ?= (Acceptance Test)

export GO111MODULE=on

default: fmt goimports set-license lint tflint docscheck

clean:
	rm -Rf $(CURDIR)/bin/*

.PHONY: tools
tools:
	GO111MODULE=off go get -u github.com/motemen/gobump/cmd/gobump
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get github.com/sacloud/addlicense
	GO111MODULE=off go get github.com/tcnksm/ghr
	GO111MODULE=on go install github.com/bflad/tfproviderdocs
	GO111MODULE=on go install github.com/bflad/tfproviderlint/cmd/tfproviderlint
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint


.PHONY: build-envs
build-envs:
	$(eval CURRENT_VERSION ?= $(shell gobump show -r sakuracloud/))
	$(eval BUILD_LDFLAGS := "-s -w \
           -X github.com/sacloud/terraform-provider-sakuracloud/sakuracloud.Revision=`git rev-parse --short HEAD` \
           -X github.com/sacloud/terraform-provider-sakuracloud/sakuracloud.Version=$(CURRENT_VERSION)")

.PHONY: build
build: build-envs
	OS=$${OS:-"`go env GOOS`"} ARCH=$${ARCH:-"`go env GOARCH`"} BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

.PHONY: build-x
build-x: build-envs build-darwin build-windows build-linux shasum

.PHONY: build-darwin
build-darwin: build-envs bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-amd64.zip

.PHONY: build-windows
build-windows: build-envs bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-amd64.zip

.PHONY: build-linux
build-linux: build-envs bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-amd64.zip

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-386.zip: build-envs
	OS="darwin"  ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-amd64.zip: build-envs
	OS="darwin"  ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-386.zip: build-envs
	OS="windows" ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-amd64.zip: build-envs
	OS="windows" ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-386.zip: build-envs
	OS="linux"   ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-amd64.zip: build-envs
	OS="linux"   ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

.PHONY: shasum
shasum:
	(cd bin/; shasum -a 256 * > terraform-provider-sakuracloud_$(CURRENT_VERSION)_SHA256SUMS)

.PHONY: release
release:
	ghr ${CURRENT_VERSION} bin/

.PHONY: test
test:
	TF_ACC= SAKURACLOUD_APPEND_USER_AGENT="$(UNIT_TEST_UA)" go test -mod=vendor -v $(TESTARGS) -timeout=30s ./...

.PHONY: testacc
testacc:
	TF_ACC=1 SAKURACLOUD_APPEND_USER_AGENT="$(ACC_TEST_UA)" go test -mod=vendor -v $(TESTARGS) -timeout 240m ./...

.PHONY: testfake
testfake:
	FAKE_MODE=1 TF_ACC=1 SAKURACLOUD_APPEND_USER_AGENT="$(ACC_TEST_UA)" go test -mod=vendor -v $(TESTARGS) -timeout 240m ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: tflint
tflint:
	tfproviderlint \
        -AT001 -AT002 -AT003 -AT004 -AT005 -AT006 -AT007 \
        -R001 -R002 -R004 -R005 -R006 \
        -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 -S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019 -S020 -S021 -S022 -S023\
        -V001 \
        ./$(PKG_NAME)

.PHONY: goimports
goimports:
	goimports -l -w $(PKG_NAME)/

.PHONY: fmt
fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

.PHONY: docscheck
docscheck:
	tfproviderdocs check \
		-require-resource-subcategory

.PHONY: docker-build
docker-build: clean
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'tools' 'build-x'"

.PHONY: set-license
set-license:
	addlicense -c $(AUTHOR) -y $(COPYRIGHT_YEAR) $(COPYRIGHT_FILES)

.PHONY: website
website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
	(cd $(GOPATH)/src/$(WEBSITE_REPO); \
	  ln -s ../../../ext/providers/sakuracloud/website/sakuracloud.erb content/source/layouts/sakuracloud.erb; \
	  ln -s ../../../../ext/providers/sakuracloud/website/docs content/source/docs/providers/sakuracloud \
	)
endif
	$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: website-lint
website-lint:
	@echo "==> Checking website against linters..."
	misspell -error -source=text website/

.PHONY: website-test
website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
	(cd $(GOPATH)/src/$(WEBSITE_REPO); \
	  ln -s ../../../ext/providers/sakuracloud/website/sakuracloud.erb content/source/layouts/sakuracloud.erb; \
	  ln -s ../../../../ext/providers/sakuracloud/website/docs source/docs/providers/sakuracloud \
	)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: website-scaffold
website-scaffold:
	go run tools/tfdocgen/cmd/gen-sakuracloud-docs/main.go website-scaffold

