# Copyright 2016-2022 terraform-provider-sakuracloud authors
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
#====================
AUTHOR          ?= The sacloud/terraform-provider-sakuracloud Authors
COPYRIGHT_YEAR  ?= 2016-2022

BIN            ?= terraform-provider-sakuracloud
BUILD_LDFLAGS   ?= "-s -w -X github.com/sacloud/terraform-provider-sakuracloud/sakuracloud.Revision=`git rev-parse --short HEAD`"

include includes/go/common.mk
include includes/go/single.mk
#====================
export GOPROXY=https://proxy.golang.org

default: fmt set-license go-licenses-check goimports lint docscheck testfake

PKG_NAME     ?= sakuracloud
WEBSITE_REPO  = github.com/hashicorp/terraform-website

UNIT_TEST_UA ?= (Unit Test)
ACC_TEST_UA ?= (Acceptance Test)

export GO111MODULE=on

.PHONY: tools
tools: dev-tools
	go install github.com/bflad/tfproviderlint/cmd/tfproviderlintx@v0.28.1
	go install github.com/bflad/tfproviderdocs@v0.9.1

.PHONY: testfake
testfake:
	FAKE_MODE=1 TF_ACC=1 SAKURACLOUD_APPEND_USER_AGENT="$(ACC_TEST_UA)" go test -v $(TESTARGS) -timeout 240m ./...

.PHONY: tflint
tflint:
	tfproviderlintx \
        -AT001 -AT002 -AT003 -AT004 -AT005 -AT006 -AT007 -AT008 -AT009 \
        -R001 -R002 -R004 -R005 -R006 -R007 -R008 -R009 -R010 -R011 -R012 -R013 -R014 -R015 \
        -R016 -R017       -R019 \
        -S001 -S002 -S003 -S004 -S005 -S006 -S007 -S008 -S009 -S010 -S011 -S012 -S013 -S014 -S015 \
        -S016 -S017 -S018 -S019 -S020 -S021 -S022 -S023 -S024 -S025 -S026 -S027 -S028 -S029 -S030 \
        -S031 -S032 -S033 -S034 -S035 -S036 -S037 \
        -V001 -V002 -V003 -V004 -V005 -V006 -V007 -V008 -V009 -V010 \
        -XR001 -XR004 \
        ./$(PKG_NAME)

.PHONY: docscheck
docscheck:
	tfproviderdocs check \
		-require-resource-subcategory \
		-require-guide-subcategory

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

