#
# Copyright 2022 The sacloud/makefile Authors
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

AUTHOR                  ?= The sacloud/makefile Authors
COPYRIGHT_YEAR          ?= 2022
COPYRIGHT_FILES         ?= $$(find . -name "*.go" -print | grep -v "/vendor/")
GO                      ?= go
DEFAULT_GOALS           ?= fmt set-license go-licenses-check goimports lint test
GOLANG_CI_LINT_VERSION  ?= v1.55.2
TEXTLINT_ACTION_VERSION ?= v0.0.3

.DEFAULT_GOAL = default

.PHONY: test
test:
	@echo "running 'go test'..."
	TF_ACC= SAKURACLOUD_APPEND_USER_AGENT="$(UNIT_TEST_UA)" go test -v $(TESTARGS) -timeout=30s ./...

.PHONY: testacc
testacc:
	@echo "running 'go test' with TESTACC=1..."
	TF_ACC=1 SAKURACLOUD_APPEND_USER_AGENT="$(ACC_TEST_UA)" go test -v $(TESTARGS) -timeout 240m ./...

.PHONY: dev-tools
dev-tools:
	$(GO) install github.com/rinchsan/gosimports/cmd/gosimports@latest
	$(GO) install golang.org/x/tools/cmd/stringer@latest
	$(GO) install github.com/sacloud/addlicense@latest
	$(GO) install github.com/client9/misspell/cmd/misspell@latest
	$(GO) install github.com/google/go-licenses@v1.0.0
	$(GO) install github.com/rhysd/actionlint/cmd/actionlint@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANG_CI_LINT_VERSION)

.PHONY: goimports
goimports: fmt
	@echo "running gosimports..."
	@gosimports -l -w .

.PHONY: fmt
fmt:
	@echo "running gofmt..."
	@find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

.PHONY: godoc
godoc:
	godoc -http=localhost:6060

.PHONY: lint
lint: lint-go lint-text lint-action

.PHONY: lint-go
lint-go:
	@echo "running golanci-lint..."
	@golangci-lint run --fix ./...

.PHONY: textlint lint-text
textlint: lint-text
lint-text:
	@echo "running textlint..."
	@docker run -t --rm -v $$PWD:/work -w /work ghcr.io/sacloud/textlint-action:$(TEXTLINT_ACTION_VERSION) .

.PHONY: lint-action
lint-action:
	@echo "running rhysd/actionlint..."
	@actionlint

.PHONY: set-license
set-license:
	@addlicense -c "$(AUTHOR)" -y "$(COPYRIGHT_YEAR)" $(COPYRIGHT_FILES)

.PHONY: go-licenses-check
go-licenses-check:
	@echo "running go-licenses..."
	@go-licenses check .
