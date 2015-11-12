include ../Makefile.variable

all: test build docker

setup: get-tools db-setup dependencies

get-tools:
	@echo "=== get tools ==="
	@brew install readline glide > /dev/null 2>&1
	@$(DOT)
	@brew install fswatch > /dev/null 2>&1
	@$(DOT)
	@go get -u "github.com/alecthomas/gometalinter"
	@$(DOT)
	@gometalinter --install --update --force
	@$(DOT)
	@$(PRINT_CHECK)

db-setup:
	@echo "=== setup db ==="
	@dropdb --if-exists ello_services_test > /dev/null 2>&1
	@$(DOT)
	@createdb ello_services_test
	@$(DOT)
	@psql ello_services_test < schema.dump > /dev/null 2>&1
	@$(PRINT_CHECK)

dependencies:export GO15VENDOREXPERIMENT=1
dependencies: clean
	@echo "=== deps ==="
	@glide install
	@$(PRINT_CHECK)

clean:
	@echo "=== cleaning ==="
	@rm -rf vendor
	@$(DOT)
	@rm -rf bin
	@$(PRINT_CHECK)

errcheck:
	@echo "=== errcheck ==="
	@errcheck github.com/ello/services/stream/...

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

install:export GO15VENDOREXPERIMENT=1
install: test
	@echo "=== go install ==="
	@go install -ldflags=$(GOLDFLAGS)

build:export GO15VENDOREXPERIMENT=1
build:
	@echo "=== go build ==="
	@mkdir -p bin/
	@go build -ldflags=$(GOLDFLAGS) -o bin/services

docker: test
	@echo "=== docker build ==="
	@docker build -t services .

test:export GO15VENDOREXPERIMENT=1
test:export LOGXI=dat:sqlx=off
test: fmt vet lint
	$(call mecho,"=== go test ===")
	@go test `glide novendor` -v -cover

test-w:
	@echo "=== testing | watch mode ==="
	@fswatch -o . -r | xargs -n1 -I{} make test

server:export GO15VENDOREXPERIMENT=1
server:
	@echo "=== server ==="
	@go run main.go

server-w:
	@echo "=== server | watch mode ==="
	@gin -a 8080 -i run

deploy: test

.PHONY: setup cloc errcheck vet lint fmt install build test deploy docker
