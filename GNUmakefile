PKG_NAME         ?= sakuracloud

AUTHOR          ?="terraform-provider-sakuracloud authors"
COPYRIGHT_YEAR  ?="2016-2019"
COPYRIGHT_FILES ?=$$(find . \( -name "*.dockerfile" -or -name "*.go" -or -name "*.sh" -or -name "*.pl" -or -name "*.bash" \) -print | grep -v "/vendor/")

UNIT_TEST_UA ?= (Unit Test)
ACC_TEST_UA ?= (Acceptance Test)

export GO111MODULE=on

default: fmt goimports lint tflint

clean:
	rm -Rf $(CURDIR)/bin/*

.PHONY: tools
tools:
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get -u github.com/motemen/gobump/cmd/gobump
	GO111MODULE=off go get github.com/sacloud/addlicense
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	GO111MODULE=off go get -u github.com/bflad/tfproviderlint/cmd/tfproviderlint

build-envs:
	$(eval CURRENT_VERSION := $(shell gobump show -r sakuracloud/))
	$(eval BUILD_LDFLAGS := "-s -w \
           -X github.com/sacloud/terraform-provider-sakuracloud/sakuracloud.Revision=`git rev-parse --short HEAD` \
           -X github.com/sacloud/terraform-provider-sakuracloud/sakuracloud.Version=$(CURRENT_VERSION)")

build: build-envs
	OS="`go env GOOS`" ARCH="`go env GOARCH`" ARCHIVE= BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

build-x: build-envs build-darwin build-windows build-linux shasum

build-darwin: build-envs bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-amd64.zip

build-windows: build-envs bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-amd64.zip

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

shasum:
	(cd bin/; shasum -a 256 * > terraform-provider-sakuracloud_$(CURRENT_VERSION)_SHA256SUMS)

test:
	TF_ACC= SAKURACLOUD_APPEND_USER_AGENT="$(UNIT_TEST_UA)" go test -mod=vendor -v $(TESTARGS) -timeout=30s ./...

testacc:
	TF_ACC=1 SAKURACLOUD_APPEND_USER_AGENT="$(ACC_TEST_UA)" go test -mod=vendor -v $(TESTARGS) -timeout 240m ./...

testfake:
	FAKE_MODE=1 TF_ACC=1 SAKURACLOUD_APPEND_USER_AGENT="$(ACC_TEST_UA)" go test -mod=vendor -v $(TESTARGS) -timeout 240m ./...

lint:
	golangci-lint run ./...

tflint:
	@tfproviderlint \
        -AT001 -AT002 -AT003 -AT004\
        -R001 -R002 -R004\
        -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 -S011 -S012 -S013 -S014 -S015 -S016 -S017 -S018 -S019\
        ./$(PKG_NAME)

goimports:
	goimports -l -w $(PKG_NAME)/

fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

.PHONY: build-docs serve-docs lint-docs
build-docs:
	sh -c "'$(CURDIR)/scripts/build_docs.sh'"

serve-docs:
	sh -c "'$(CURDIR)/scripts/serve_docs.sh'"

lint-docs:
	sh -c "'$(CURDIR)/scripts/lint_docs.sh'"

serve-english-docs:
	sh -c "'$(CURDIR)/scripts/serve_english_docs.sh'"

docker-test: 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'test'"

docker-testacc: 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'testacc'"

docker-testacc-resource:
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'testacc-resource'"

docker-build: clean 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'build-x'"


.PHONY: default test vet testacc fmt fmtcheck

set-license:
	@addlicense -c $(AUTHOR) -y $(COPYRIGHT_YEAR) $(COPYRIGHT_FILES)
