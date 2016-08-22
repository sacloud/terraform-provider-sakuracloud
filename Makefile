TEST1?=./builtin/bins/provider-sakuracloud
TEST2?=./builtin/providers/sakuracloud
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

default: test vet

clean:
	rm -Rf $(CURDIR)/bin/*

build: clean vet
	govendor build -ldflags "-s -w" -o $(CURDIR)/bin/terraform-provider-sakuracloud $(CURDIR)/builtin/bins/provider-sakuracloud/main.go

build-x: clean vet
	sh -c "'$(CURDIR)/scripts/build.sh'"

test: vet
	TF_ACC= govendor test $(TEST1) $(TESTARGS) -timeout=30s -parallel=4 ; \
	TF_ACC= govendor test $(TEST2) $(TESTARGS) -timeout=30s -parallel=4

testacc: vet
	TF_ACC=1 govendor test $(TEST1) -v $(TESTARGS) -timeout 120m ; \
	TF_ACC=1 govendor test $(TEST2) -v $(TESTARGS) -timeout 120m

testacc-resource: vet
	TF_ACC=1 govendor test $(TEST1) -v $(TESTARGS) -run="^TestAccResource" -timeout 120m ; \
	TF_ACC=1 govendor test $(TEST2) -v $(TESTARGS) -run="^TestAccResource" -timeout 120m

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

docker-test: 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'test'"

docker-testacc: 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'testacc'"

docker-testacc-resource:
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'testacc-resource'"

docker-build: clean 
	sh -c "'$(CURDIR)/scripts/build_on_docker.sh' 'build-x'"


.PHONY: default test vet testacc fmt fmtcheck
