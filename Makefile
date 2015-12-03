include Makefile.variable

announce:
	$(call becho,"=== Ello Go Project ===")

setup: announce get-tools
	@$(PRINT_LINE)
	$(call becho,"~~~    github.com/ello/ello-go/streams   ~~~")
	@$(PRINT_LINE)
	@cd streams && $(MAKE) setup

get-tools: get-tools-ci
	@brew rm --force fswatch readline > /dev/null 2>&1
	@brew install fswatch readline > /dev/null 2>&1
	@$(PRINT_OK)

setup-ci: announce get-tools-ci
	@$(PRINT_LINE)
	$(call becho,"~~~    github.com/ello/ello-go/streams   ~~~")
	@$(PRINT_LINE)
	@cd streams && $(MAKE) setup

get-tools-ci:
	@go install -u "github.com/Masterminds/glide"
	@go build "github.com/Masterminds/glide"
	@go get -u "github.com/alecthomas/gometalinter" > /dev/null 2>&1
	@gometalinter --install --update --force  > /dev/null 2>&1
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


install:export GO15VENDOREXPERIMENT=1
install: test
	@echo "=== go install ==="
	@go install -ldflags=$(GOLDFLAGS)

all:
	@$(PRINT_LINE)
	$(call becho,"~~~    github.com/ello/ello-go/streams   ~~~")
	@$(PRINT_LINE)
	@cd streams && $(MAKE) all

test:
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
