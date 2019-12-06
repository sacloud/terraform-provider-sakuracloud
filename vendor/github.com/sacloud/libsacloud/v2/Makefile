#
# Copyright 2016-2019 The Libsacloud Authors
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
TEST            ?=$$(go list ./... | grep -v vendor)
VETARGS         ?=-all
GOFMT_FILES     ?=$$(find . -name '*.go' | grep -v vendor)
AUTHOR          ?="The Libsacloud Authors"
COPYRIGHT_YEAR  ?="2016-2019"
COPYRIGHT_FILES ?=$$(find . -name "*.go" -print | grep -v "/vendor/")
export GO111MODULE=on

default: gen fmt set-license goimports golint vet test

test:
	TESTACC= go test ./... $(TESTARGS) -v -timeout=120m -parallel=8 ;

testacc:
	TESTACC=1 go test ./... $(TESTARGS) -v -timeout=120m -parallel=8 ;

.PHONY: clean
clean:
	rm -f sacloud/zz_*.go; \
	rm -f sacloud/fake/zz_*.go \
	rm -f sacloud/naked/zz_*.go \
	rm -f sacloud/stub/zz_*.go \
	rm -f sacloud/trace/zz_*.go

.PHONY: gen
gen: _gen set-license

.PHONY: _gen
_gen:
	go generate ./...; gofmt -s -l -w $(GOFMT_FILES); goimports -l -w $(GOFMT_FILES)

.PHONY: gen_fake_data
gen_fake_data: _gen_fake_data set-license

.PHONY: _gen_fake_data
_gen_fake_data:
	go run -mod=vendor internal/tools/gen-api-fake-data/main.go ; gofmt -s -l -w $(GOFMT_FILES); goimports -l -w $(GOFMT_FILES)

vet: golint
	go vet ./...

golint: 
	test -z "$$(golint ./... | grep -v 'tools/' | grep -v 'vendor/' | grep -v '_string.go' | tee /dev/stderr )"

goimports: fmt
	goimports -l -w $(GOFMT_FILES)

fmt:
	gofmt -s -l -w $(GOFMT_FILES)

godoc:
	@echo "URL: http://localhost:6060/pkg/github.com/sacloud/libsacloud/"; \
	docker run -it --rm -v $$PWD:/go/src/github.com/sacloud/libsacloud -p 6060:6060 golang:1.12 godoc -http=:6060

.PHONY: tools
tools:
	GO111MODULE=off go get golang.org/x/lint/golint
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports
	GO111MODULE=off go get golang.org/x/tools/cmd/stringer
	GO111MODULE=off go get github.com/motemen/gobump
	GO111MODULE=off go get github.com/sacloud/addlicense

.PHONY: bump-patch bump-minor bump-major version
bump-patch:
	@gobump patch -w ; echo "next version is v`gobump show -r`"

bump-minor:
	@gobump minor -w ; echo "next version is v`gobump show -r`"

bump-major:
	@gobump major -w ; echo "next version is v`gobump show -r`"

version:
	@gobump show -r

git-tag:
	git tag v`gobump show -r`

set-license:
	@addlicense -c $(AUTHOR) -y $(COPYRIGHT_YEAR) $(COPYRIGHT_FILES)

