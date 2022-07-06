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

GO             ?= go
BIN            ?= TODO_PLEASE_SET_BIN_VARIABLE
GO_ENTRY_FILE  ?= main.go
GO_FILES       ?= $(shell find . -name '*.go')
BUILD_LDFLAGS  ?=

.PHONY: install
install:
	@echo "running 'go install'..."
	$(GO) install

.PHONY: build
build: $(BIN)

$(BIN): $(GO_FILES) go.mod go.sum
	@echo "running 'go build'..."
	@GOOS=$${OS:-"`$(GO) env GOOS`"} GOARCH=$${ARCH:-"`$(GO) env GOARCH`"} CGO_ENABLED=0 $(GO) build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) $(GO_ENTRY_FILE)

.PHONY: clean
clean:
	@echo "cleaning..."
	rm -rf $(BIN)

DEFAULT_GOALS += build