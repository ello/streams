include ../Makefile.variable

all: test

setup: get-tools dependencies

get-tools:
	$(call mecho,"=== get tools ===")
	# install fswatch readline and glide through homebrew
	@brew install fswatch readline glide > /dev/null 2>&1
	@$(DOT)
	# install gometalinter
	@go get -u "github.com/alecthomas/gometalinter" > /dev/null 2>&1
	@$(DOT)
	# install gometalinter updates to the linter packages
	@gometalinter --install --update --force > /dev/null 2>&1
	@$(DOT)
	@$(PRINT_CHECK)

dependencies:export GO15VENDOREXPERIMENT=1
dependencies: clean
	$(call mecho,"=== deps ===")
	@glide up
	@glide rebuild
	@$(PRINT_CHECK)

clean:
	$(call mecho,"=== cleaning ===")
	@rm -rf vendor
	@$(DOT)
	@rm -rf bin
	@$(PRINT_CHECK)

vet:export GO15VENDOREXPERIMENT=1
vet:
	$(call mecho,"=== go vet ===")
	@go vet `glide novendor`

# TODO Re-enable these linters once vendor support is better (potentially 1.6)
lint:export GO15VENDOREXPERIMENT=1
lint:
	$(call mecho,"=== go lint ===")
	@gometalinter --vendor  --deadline=10s --disable=gotype --disable=aligncheck --disable=structcheck --dupl-threshold=70 `glide novendor`

fmt:export GO15VENDOREXPERIMENT=1
fmt:
	$(call mecho,"=== go fmt ===")
	@go fmt `glide novendor`

test:export GO15VENDOREXPERIMENT=1
test:export LOGXI=dat:sqlx=off
test: fmt vet lint
	$(call mecho,"=== go test ===")
	@go test `glide novendor` -v -cover

test-w:
	$(call mecho,"=== testing | watch mode ===")
	@fswatch -o . -r | xargs -n1 -I{} make test

deploy: test

.PHONY: setup cloc errcheck vet lint fmt install build test deploy docker
