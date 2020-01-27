#
# Copyright 2016-2020 The terraform-provider-sakuracloud Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
TEST1           ?=./
TEST2           ?=./sakuracloud
VETARGS         ?=-all
GOFMT_FILES     ?=$$(find . -name '*.go' | grep -v vendor)
GOGEN_FILES     ?=$$(go list ./... | grep -v vendor)
GOLINT_TARGETS  ?= $$(golint github.com/sacloud/terraform-provider-sakuracloud/sakuracloud | grep -v 'underscores in Go names' | tee /dev/stderr)
CURRENT_VERSION ?= $(shell git log --merges --oneline | perl -ne 'if(m/^.+Merge pull request \#[0-9]+ from .+\/bump-version-([0-9\.]+)/){print $$1;exit}')
AUTHOR          ?="terraform-provider-sakuracloud authors"
COPYRIGHT_YEAR  ?="2016-2020"
COPYRIGHT_FILES ?=$$(find . \( -name "*.dockerfile" -or -name "*.go" -or -name "*.sh" -or -name "*.pl" -or -name "*.bash" \) -print | grep -v "/vendor/")

BUILD_LDFLAGS = "-s -w \
	  -X github.com/sacloud/terraform-provider-sakuracloud/sakuracloud.Revision=`git rev-parse --short HEAD` \
	  -X github.com/sacloud/terraform-provider-sakuracloud/sakuracloud.Version=$(CURRENT_VERSION)"

export GO111MODULE=on

default: vet build

clean:
	rm -Rf $(CURDIR)/bin/*

.PHONY: tools
tools:
	GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=off go get github.com/x-motemen/gobump/cmd/gobump
	GO111MODULE=off go get github.com/sacloud/addlicense
	GO111MODULE=off go get github.com/tcnksm/ghr

build:
	OS="`go env GOOS`" ARCH="`go env GOARCH`" ARCHIVE= BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

build-x: build-darwin build-windows build-linux shasum

build-darwin: bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-amd64.zip

build-windows: bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-amd64.zip

build-linux: bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-386.zip bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-amd64.zip

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-386.zip:
	OS="darwin"  ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_darwin-amd64.zip:
	OS="darwin"  ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-386.zip:
	OS="windows" ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_windows-amd64.zip:
	OS="windows" ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-386.zip:
	OS="linux"   ARCH="386"   ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

bin/terraform-provider-sakuracloud_$(CURRENT_VERSION)_linux-amd64.zip:
	OS="linux"   ARCH="amd64" ARCHIVE=1 BUILD_LDFLAGS=$(BUILD_LDFLAGS) CURRENT_VERSION=$(CURRENT_VERSION) sh -c "'$(CURDIR)/scripts/build.sh'"

shasum:
	(cd bin/; shasum -a 256 * > terraform-provider-sakuracloud_$(CURRENT_VERSION)_SHA256SUMS)

.PHONY: release
release:
	ghr v${CURRENT_VERSION} bin/

test:
	TF_ACC= go test $(TEST1) -v $(TESTARGS) -timeout=30s -parallel=4 ; \
	TF_ACC= go test $(TEST2) -v $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST1) -v $(TESTARGS) -timeout 240m ; \
	TF_ACC=1 go test $(TEST2) -v $(TESTARGS) -timeout 240m

testacc-resource:
	TF_ACC=1 go test $(TEST1) -v $(TESTARGS) -run="^TestAccResource" -timeout 240m ; \
	TF_ACC=1 go test $(TEST2) -v $(TESTARGS) -run="^TestAccResource" -timeout 240m

vet: golint
	@echo "go vet $(VETARGS) ."
	@go vet $(VETARGS) ; if [ $$? -eq 1 ]; then \
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

.PHONY: default test vet testacc fmt fmtcheck

set-license:
	@addlicense -c $(AUTHOR) -y $(COPYRIGHT_YEAR) $(COPYRIGHT_FILES)
