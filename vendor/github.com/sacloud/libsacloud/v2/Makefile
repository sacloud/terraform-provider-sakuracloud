#
# Copyright 2016-2020 The Libsacloud Authors
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
AUTHOR          ?="The Libsacloud Authors"
COPYRIGHT_YEAR  ?="2016-2020"
COPYRIGHT_FILES ?=$$(find . -name "*.go" -print | grep -v "/vendor/")

default: gen fmt set-license goimports lint test

.PHONY: test
test:
	TESTACC= go test ./... $(TESTARGS) -v -timeout=120m -parallel=8 -race;

.PHONY: testacc
testacc:
	TESTACC=1 go test ./... $(TESTARGS) -v -timeout=120m -parallel=8 ;

.PHONY: tools
tools:
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get golang.org/x/tools/cmd/stringer
	GO111MODULE=off go get github.com/sacloud/addlicense
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/v1.23.8/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.23.8

.PHONY: clean
clean:
	rm -f sacloud/zz_*.go; \
	rm -f sacloud/fake/zz_*.go \
	rm -f sacloud/naked/zz_*.go \
	rm -f sacloud/stub/zz_*.go \
	rm -f sacloud/trace/zz_*.go

.PHONY: gen
gen: _gen fmt goimports set-license

.PHONY: _gen
_gen:
	go generate ./...

.PHONY: gen_fake_data
gen_fake_data: _gen_fake_data fmt goimports set-license

.PHONY: _gen_fake_data
_gen_fake_data:
	go run internal/tools/gen-api-fake-data/main.go

.PHONY: goimports
goimports: fmt
	goimports -l -w .

.PHONY: fmt
fmt:
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

.PHONY: godoc
godoc:
	@echo "URL: http://localhost:6060/pkg/github.com/sacloud/libsacloud/"; \
	docker run -it --rm -v $$PWD:/go/src/github.com/sacloud/libsacloud -p 6060:6060 golang:1.12 godoc -http=:6060

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: set-license
set-license:
	@addlicense -c $(AUTHOR) -y $(COPYRIGHT_YEAR) $(COPYRIGHT_FILES)

