include Makefile.variable

announce:
	$(call becho,"=== Ello Streams Project ===")

get-tools:
	@go get -u "github.com/Masterminds/glide"
	@go build "github.com/Masterminds/glide"
	@go get -u "github.com/alecthomas/gometalinter"
	# This is broken for the moment due to https://github.com/opennota/check/issues/25
	# When that's fixed we can go back to the `gometalinter` command instead of individual installs
	# @gometalinter --install --update --force
	@go get -u -f "golang.org/x/tools/cmd/goimports"
	@go get -u -f "github.com/kisielk/errcheck"
	@go get -u -f "github.com/gordonklaus/ineffassign"
	@go get -u -f "github.com/mibk/dupl"
	@go get -u -f "github.com/alecthomas/gometalinter"
	@go get -u -f "golang.org/x/tools/cmd/gotype"
	@go get -u -f "github.com/tsenart/deadcode"
	@go get -u -f "github.com/alecthomas/gocyclo"
	@go get -u -f "github.com/mvdan/interfacer/cmd/interfacer"
	@go get -u -f "github.com/golang/lint/golint"
	@$(PRINT_OK)

install:export GO15VENDOREXPERIMENT=1
install: test
	@echo "=== go install ==="
	@go install -ldflags=$(GOLDFLAGS)
	@$(PRINT_OK)

# From streams

all: test build docker

setup: announce get-tools dependencies

#TODO Try getting rid of the vendor flag env var after 1.6 is out
dependencies:export GO15VENDOREXPERIMENT=1
dependencies:
	@glide install
	@glide rebuild
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
	@gometalinter --vendor  --deadline=10s --disable=gotype --disable=varcheck --disable=aligncheck --disable=structcheck --disable=errcheck --disable=interfacer --dupl-threshold=70 `glide novendor`
	@$(PRINT_OK)

fmt:export GO15VENDOREXPERIMENT=1
fmt:
	@gofmt -s -w `glide nv | sed 's/\.\.\./*.go/g' | sed 's/.\///'`
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

server:export GO15VENDOREXPERIMENT=1
server:
	@go run main.go

server-w:
	@gin -a 8080 -i run

.PHONY: setup cloc errcheck vet lint fmt install build test deploy docker
