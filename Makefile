TEST1?=./builtin/bins/provider-sakuracloud
TEST2?=./builtin/providers/sakuracloud
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
CURRENT_VERSION = $(shell git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $$1;exit}')

BUILD_LDFLAGS = "-s -w \
	  -X github.com/yamamoto-febc/terraform-provider-sakuracloud/version.Revision=`git rev-parse --short HEAD` \
	  -X github.com/yamamoto-febc/terraform-provider-sakuracloud/version.Version=$(CURRENT_VERSION)"

default: test vet

clean:
	rm -Rf $(CURDIR)/bin/*

build: clean vet
	OS="`go env GOOS`" ARCH="`go env GOARCH`" ARCHIVE= BUILD_LDFLAGS=$(BUILD_LDFLAGS) sh -c "'$(CURDIR)/scripts/build.sh'"

build-x: build-darwin build-windows build-linux

build-darwin: bin/terraform-provider-sakuracloud_darwin-386.zip bin/terraform-provider-sakuracloud_darwin-amd64.zip

build-windows: bin/terraform-provider-sakuracloud_windows-386.zip bin/terraform-provider-sakuracloud_windows-amd64.zip

build-linux: bin/terraform-provider-sakuracloud_linux-386.zip bin/terraform-provider-sakuracloud_linux-amd64.zip

bin/terraform-provider-sakuracloud_darwin-386.zip:
	OS="darwin"  ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_darwin-amd64.zip:
	OS="darwin"  ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_windows-386.zip:
	OS="windows" ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_windows-amd64.zip:
	OS="windows" ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_linux-386.zip:
	OS="linux"   ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_linux-amd64.zip:
	OS="linux"   ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) sh -c "'$(CURDIR)/scripts/build.sh'"



test: vet
	TF_ACC= go test $(TEST1) $(TESTARGS) -timeout=30s -parallel=4 ; \
	TF_ACC= go test $(TEST2) $(TESTARGS) -timeout=30s -parallel=4

testacc: vet
	TF_ACC=1 go test $(TEST1) -v $(TESTARGS) -timeout 120m ; \
	TF_ACC=1 go test $(TEST2) -v $(TESTARGS) -timeout 120m

testacc-resource: vet
	TF_ACC=1 go test $(TEST1) -v $(TESTARGS) -run="^TestAccResource" -timeout 120m ; \
	TF_ACC=1 go test $(TEST2) -v $(TESTARGS) -run="^TestAccResource" -timeout 120m

vet: fmt
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

.PHONY: build-docs serve-docs
build-docs:
	sh -c "'$(CURDIR)/scripts/build_docs.sh'"

serve-docs:
	sh -c "'$(CURDIR)/scripts/serve_docs.sh'"

docker-test: 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'test'"

docker-testacc: 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'testacc'"

docker-testacc-resource:
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'testacc-resource'"

docker-build: clean 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'build-x'"


.PHONY: default test vet testacc fmt fmtcheck
