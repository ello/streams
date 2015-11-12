include ../Makefile.variable

all: test build docker

setup: get-tools db-setup dependencies

get-tools:
	@echo "=== get tools ==="
	@brew install fswatch readline glide > /dev/null 2>&1
	@go get -u "github.com/alecthomas/gometalinter"
	@gometalinter --install --update --force
	@$(PRINT_OK)

db-setup:
	@echo "=== setup db ==="
	@dropdb --if-exists ello_services_test > /dev/null 2>&1
	@createdb ello_services_test
	@psql ello_services_test < schema.dump > /dev/null 2>&1
	@$(PRINT_OK)

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

install:export GO15VENDOREXPERIMENT=1
install: test
	@echo "=== go install ==="
	@go install -ldflags=$(GOLDFLAGS)

build:export GO15VENDOREXPERIMENT=1
build:
	@mkdir -p bin/
	@go build -ldflags=$(GOLDFLAGS) -o bin/streams
	@$(PRINT_OK)

rebuild: clean build

docker: test
	@echo "=== docker build ==="
	@docker build -t services .

test:export GO15VENDOREXPERIMENT=1
test:export LOGXI=dat:sqlx=off
test: fmt vet lint
	@go test `glide novendor` -v -cover
	@$(PRINT_OK)

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
