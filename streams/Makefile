include ../Makefile.variable

all: test build docker

setup: dependencies

#TODO Try getting rid of the vendor flag env var after 1.6 is out
dependencies:export GO15VENDOREXPERIMENT=1
dependencies:
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
	@gofmt -s -w `glide nv | sed 's/\.\.\./*.go/g' | sed 's/.\///'`
	@$(PRINT_OK)

install:export GO15VENDOREXPERIMENT=1
install: test
	@go install -ldflags=$(GOLDFLAGS)
	@$(PRINT_OK)

build:export GO15VENDOREXPERIMENT=1
build:
	@mkdir -p bin/
	@go build -ldflags $(GOLDFLAGS) -o bin/streams
	@$(PRINT_OK)

rebuild: clean build

docker: test
	@docker build -t streams . > /dev/null 2>&1
	@$(PRINT_OK)

test:export GO15VENDOREXPERIMENT=1
test:export LOGXI=dat:sqlx=off
test: fmt vet lint
	@go test `glide novendor` -cover
	@$(PRINT_OK)

test-w:
	@fswatch -o . -r | xargs -n1 -I{} make test

server:export GO15VENDOREXPERIMENT=1
server:
	@go run main.go

server-w:
	@gin -a 8080 -i run

deploy: docker
	#stuff to send docker image up somewhere here
	@$(PRINT_OK)


.PHONY: setup cloc errcheck vet lint fmt install build test deploy docker
