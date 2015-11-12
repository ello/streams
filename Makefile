include Makefile.variable

all: test build docker

setup: get-tools db-setup dependencies

get-tools:
	@echo "=== get tools ==="
	@brew install fswatch readline glide > /dev/null 2>&1
	@go get -u "github.com/alecthomas/gometalinter"
	@gometalinter --install --update --force
	@$(PRINT_OK)

# db-setup:
# 	@echo "=== setup db ==="
# 	@dropdb --if-exists ello_services_test > /dev/null 2>&1
# 	@/bin/echo -n "."
# 	@createdb ello_services_test
# 	@/bin/echo -n "."
# 	@psql ello_services_test < schema.dump > /dev/null 2>&1
# 	@echo "√"

# dependencies:export GO15VENDOREXPERIMENT=1
# dependencies: clean
# 	@echo "=== deps ==="
# 	@glide install
# 	@/bin/echo -n "√"

# clean:
# 	@echo "=== cleaning ==="
# 	@rm -rf vendor
# 	@/bin/echo -n "."
# 	@rm -rf bin
# 	@echo "√"

# errcheck:
# 	@echo "=== errcheck ==="
# 	@errcheck github.com/ello/services/stream/...

vet:export GO15VENDOREXPERIMENT=1
vet:
	@echo "=== go vet ==="
	@go vet `glide novendor`

lint:export GO15VENDOREXPERIMENT=1
lint:
	@echo "=== go lint ==="
	# TODO Re-enable these linters once vendor support is better (potentially 1.6)
	@gometalinter --vendor  --deadline=10s --disable=gotype --disable=aligncheck --disable=structcheck --dupl-threshold=70 `glide novendor`

fmt:export GO15VENDOREXPERIMENT=1
fmt:
	@echo "=== go fmt ==="
	@go fmt `glide novendor`

install:export GO15VENDOREXPERIMENT=1
install: test
	@echo "=== go install ==="
	@go install -ldflags=$(GOLDFLAGS)

build:export GO15VENDOREXPERIMENT=1
build:
	@echo "=== build ==="
	cd common && $(MAKE) all
	cd streams && $(MAKE) all

test:
	@$(PRINT_LINE)
	$(call becho,"~~~    github.com/ello/ello-go/common    ~~~")
	@$(PRINT_LINE)
	@cd common && $(MAKE) test
	@$(PRINT_LINE)
	$(call becho,"~~~    github.com/ello/ello-go/streams   ~~~")
	@$(PRINT_LINE)
	@cd streams && $(MAKE) test

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
