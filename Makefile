LDFLAG_LOCATION := github.com/metal-toolbox/mctl/internal/version
GIT_COMMIT  := $(shell git rev-parse --short HEAD)
GIT_BRANCH  := $(shell git symbolic-ref -q --short HEAD)
GIT_SUMMARY := $(shell git describe --tags --dirty --always)
VERSION     := $(shell git describe --tags 2> /dev/null)
BUILD_DATE  := $(shell date +%s)

.DEFAULT_GOAL := help

## lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61 run --config .golangci.yml

## lint-fix - auto fix lint errors - for the linters that support auto fix
lint-fix:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61 run --fix --config .golangci.yml

## Go test
test:
	go test -timeout 1m -v -covermode=atomic ./...

## build osx bin
build-osx:
ifeq (${GO_VERSION}, 0)
	$(error build requies go version 1.17.n or higher)
endif
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o mctl \
		-ldflags \
		"-X ${LDFLAG_LOCATION}.GitCommit=${GIT_COMMIT} \
		 -X ${LDFLAG_LOCATION}.GitBranch=${GIT_BRANCH} \
		 -X ${LDFLAG_LOCATION}.GitSummary=${GIT_SUMMARY} \
		 -X ${LDFLAG_LOCATION}.AppVersion=${VERSION} \
		 -X ${LDFLAG_LOCATION}.BuildDate=${BUILD_DATE}"

## Build linux bin
build-linux:
ifeq (${GO_VERSION}, 0)
	$(error build requies go version 1.16.n or higher)
endif
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mctl \
		-ldflags \
		"-X ${LDFLAG_LOCATION}.GitCommit=${GIT_COMMIT} \
		 -X ${LDFLAG_LOCATION}.GitBranch=${GIT_BRANCH} \
		 -X ${LDFLAG_LOCATION}.GitSummary=${GIT_SUMMARY} \
		 -X ${LDFLAG_LOCATION}.AppVersion=${VERSION} \
		 -X ${LDFLAG_LOCATION}.BuildDate=${BUILD_DATE}"

## Generate CLI docs
gen-docs:
	CGO_ENABLED=0 go build -o mctl
	./mctl gendocs

# https://gist.github.com/prwhite/8168133
# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-${TARGET_MAX_CHAR_NUM}s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' ${MAKEFILE_LIST}
