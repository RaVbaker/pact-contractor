PROJECTNAME=$(shell basename "$(PWD)")
MAKEFLAGS += --silent

## build: Builds the project binary `bin/pact-contractor`
build:
	go build -o bin/$(PROJECTNAME) main.go

## run: Run given command. e.g; make run cmd="push -b my-bucket"
run:
	go run main.go $(cmd)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo