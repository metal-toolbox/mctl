.DEFAULT_GOAL := help

## lint
lint:
	golangci-lint run --config .golangci.yml

## Go test
test:
	CGO_ENABLED=0 go test -timeout 1m -v -covermode=atomic ./...

## build osx bin
build-osx:
ifeq ($(GO_VERSION), 0)
	$(error build requies go version 1.17.n or higher)
endif
	  GOOS=darwin GOARCH=amd64 go build -o mctl


## Build linux bin
build-linux:
ifeq ($(GO_VERSION), 0)
	$(error build requies go version 1.16.n or higher)
endif
	GOOS=linux GOARCH=amd64 go build -o mctl

## Generate CLI docs
gen-docs:
	go build -o mctl
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
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
