include ../Makefile.variable

all: test

setup: dependencies

dependencies:export GO15VENDOREXPERIMENT=1
dependencies: clean
	@glide install > /dev/null 2>&1
	@glide rebuild > /dev/null 2>&1
	@$(PRINT_OK)

clean:
	@rm -rf vendor
	@rm -rf bin
	@$(PRINT_OK)

vet:export GO15VENDOREXPERIMENT=1
vet:
	@go vet `glide novendor`
	@$(PRINT_OK)

# TODO Re-enable these linters once vendor support is better (potentially 1.6)
lint:export GO15VENDOREXPERIMENT=1
lint:
	@gometalinter --vendor  --deadline=10s --disable=gotype --disable=aligncheck --disable=structcheck --dupl-threshold=70 `glide novendor`
	@$(PRINT_OK)

fmt:export GO15VENDOREXPERIMENT=1
fmt:
	@go fmt `glide novendor`
	@$(PRINT_OK)

test:export GO15VENDOREXPERIMENT=1
test:export LOGXI=dat:sqlx=off
test: fmt vet lint
	@go test `glide novendor` -v -cover 
	@$(PRINT_OK)

test-w:
	$(call mecho,"=== testing | watch mode ===")
	@fswatch -o . -r | xargs -n1 -I{} make test

deploy: test

.PHONY: setup cloc errcheck vet lint fmt install build test deploy docker
