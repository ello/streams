include Makefile.variable

announce:
	$(call becho,"=== Ello Streams Project ===")

# maybe activate them later
get-tools:
	# go install "github.com/Masterminds/glide"
	# @go get -u "github.com/Masterminds/glide"
	# @go build "github.com/Masterminds/glide"
	# @go get -u "github.com/alecthomas/gometalinter"
	# This is broken for the moment due to https://github.com/opennota/check/issues/25
	# When that's fixed we can go back to the `gometalinter` command instead of individual installs
	# @gometalinter --install --update --force
	# go install "golang.org/x/tools/cmd/goimports"
	# go install "github.com/kisielk/errcheck"
	# go install "github.com/gordonklaus/ineffassign"
	# go install "github.com/mibk/dupl"
	# go install "github.com/golangci/golangci-lint"
	# go install "golang.org/x/tools/cmd/gotype"
	# go install "github.com/tsenart/deadcode"
	# go install "github.com/alecthomas/gocyclo"
	# go install "github.com/mvdan/interfacer/cmd/interfacer"
	# go install "github.com/golang/lint/golint"
	@$(PRINT_OK)

install: test
	@echo "=== go install ==="
	@go install -ldflags=$(GOLDFLAGS)
	@$(PRINT_OK)

# From streams

all: test build docker

setup: announce get-tools dependencies

dependencies:
	@go mod vendor
	@$(PRINT_OK)

clean:
	@rm -rf vendor
	@rm -rf bin
	@$(PRINT_OK)


# TODO Re-enable these linters once vendor support is better (potentially 1.6)
# lint:export GO15VENDOREXPERIMENT=1
# lint:
# 	@gometalinter --vendor  --deadline=10s --disable=gotype --disable=varcheck --disable=aligncheck --disable=structcheck --disable=errcheck --disable=interfacer --dupl-threshold=70 `glide novendor`
# 	@$(PRINT_OK)

# vet:
# 	@go vet `glide novendor`
# 	@$(PRINT_OK)

fmt:
	@gofmt -s -w `glide nv | sed 's/\.\.\./*.go/g' | sed 's/.\///'`
	@$(PRINT_OK)

build:
	@mkdir -p bin/
	@go build -ldflags $(GOLDFLAGS) -o bin/streams
	@$(PRINT_OK)

rebuild: clean build

docker: test
	@docker build -t streams . > /dev/null 2>&1
	@$(PRINT_OK)

test:export LOGXI=dat:sqlx=off
test: fmt #vet lint
	@go test -v ./...
	@$(PRINT_OK)

server:
	@go run main.go

server-w:
	@gin -a 8080 -i run

.PHONY: setup cloc errcheck vet lint fmt install build test deploy docker
