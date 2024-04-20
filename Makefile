# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd/goxkcd/
BINARY_NAME := goxkcd

ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif


.PHONY: help
## help: print this help message
help:
	@echo "$$(tput setaf 2)Available rules:$$(tput sgr0)";sed -ne"/^## /{h;s/.*//;:d" -e"H;n;s/^## /---/;td" -e"s/:.*//;G;s/\\n## /===/;s/\\n//g;p;}" ${MAKEFILE_LIST}|awk -F === -v n=$$(tput cols) -v i=4 -v a="$$(tput setaf 6)" -v z="$$(tput sgr0)" '{printf"- %s%s%s\n",a,$$1,z;m=split($$2,w,"---");l=n-i;for(j=1;j<=m;j++){l-=length(w[j])+1;if(l<= 0){l=n-i-length(w[j])-1;}printf"%*s%s\n",-i," ",w[j];}}'

.PHONY: confirm
# confirm: ask for confirmation
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

.PHONY: no-dirty
## no diff
no-dirty:
	git diff --exit-code

.PHONY: tidy
## tidy: format code and tidy modfile
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: audit
## audit: run quality control checks
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

.PHONY: test
## test: run all tests
test:
	go test -v -race -buildvcs ./...

.PHONY: test/cover
## test/cover: run all tests and display coverage
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

.PHONY: build
## build: build the application
build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
	echo "Built /tmp/bin/${BINARY_NAME}"

.PHONY: run
## run: run the  application
run: build
	/tmp/bin/${BINARY_NAME} $(RUN_ARGS)

.PHONY: bench
## bench: run benchmarks
bench:
	go test -bench=. ./... -v

.DEFAULT_GOAL := build