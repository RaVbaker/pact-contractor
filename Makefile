PROJECTNAME=$(shell basename "$(PWD)")
MAKEFLAGS += --silent

## build: Builds the project binary `bin/pact-contractor`
build:
	go build -o bin/$(PROJECTNAME) main.go

## run: Run given command. e.g; make run cmd="push -b my-bucket"
run:
	go run main.go $(cmd)

## release: Releases new version of the binary and submits to GitHub. Remember to have the GITHUB_TOKEN env var present. Provide VERSION to set the released version. E.g. make release VERSION=v0.1.1
release:
	git tag -a $(VERSION)
	git push origin $(VERSION)
	goreleaser --rm-dist

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo