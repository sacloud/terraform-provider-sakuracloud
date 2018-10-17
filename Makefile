TEST1?=./
TEST2?=./sakuracloud
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOLINT_TARGETS?=$$(golint github.com/sacloud/terraform-provider-sakuracloud/sakuracloud | grep -v 'underscores in Go names' | tee /dev/stderr)
CURRENT_VERSION = $(shell git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $$1;exit}')
PROTOCOL_VERSION = $(shell go run tools/plugin-protocol-version/main.go)

BUILD_LDFLAGS = "-s -w \
	  -X github.com/sacloud/terraform-provider-sakuracloud/version.Revision=`git rev-parse --short HEAD` \
	  -X github.com/sacloud/terraform-provider-sakuracloud/version.Version=$(CURRENT_VERSION)"

default: test vet

clean:
	rm -Rf $(CURDIR)/bin/*

build: clean vet
	OS="`go env GOOS`" ARCH="`go env GOARCH`" ARCHIVE= BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) PROTOCOL_VERSION=$(PROTOCOL_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

build-x: build-darwin build-windows build-linux shasum

build-darwin: bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-amd64.zip

build-windows: bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-amd64.zip

build-linux: bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-amd64.zip

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-386.zip:
	OS="darwin"  ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) PROTOCOL_VERSION=$(PROTOCOL_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-amd64.zip:
	OS="darwin"  ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) PROTOCOL_VERSION=$(PROTOCOL_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-386.zip:
	OS="windows" ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) PROTOCOL_VERSION=$(PROTOCOL_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-amd64.zip:
	OS="windows" ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) PROTOCOL_VERSION=$(PROTOCOL_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-386.zip:
	OS="linux"   ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) PROTOCOL_VERSION=$(PROTOCOL_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-amd64.zip:
	OS="linux"   ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) PROTOCOL_VERSION=$(PROTOCOL_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

shasum:
	(cd bin/; shasum -a 256 * > terraform-provider-sakuracloud_$(CURRENT_VERSION)_SHA256SUMS)

test: vet
	TF_ACC= go test $(TEST1) -v $(TESTARGS) -timeout=30s -parallel=4 ; \
	TF_ACC= go test $(TEST2) -v $(TESTARGS) -timeout=30s -parallel=4

testacc: vet
	TF_ACC=1 go test $(TEST1) -v $(TESTARGS) -timeout 240m ; \
	TF_ACC=1 go test $(TEST2) -v $(TESTARGS) -timeout 240m

testacc-resource: vet
	TF_ACC=1 go test $(TEST1) -v $(TESTARGS) -run="^TestAccResource" -timeout 240m ; \
	TF_ACC=1 go test $(TEST2) -v $(TESTARGS) -run="^TestAccResource" -timeout 240m

vet: golint
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

golint: goimports
	test -z "$(GOLINT_TARGETS)"

goimports: fmt
	goimports -w $(GOFMT_FILES)

fmt:
	gofmt -w $(GOFMT_FILES)

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
