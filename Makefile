EVENTS_DIR = ./events

validate: lambda internal
	sam validate

build: lambda internal
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 sam build

.PHONY: local-invoke
local-invoke: build
	$(foreach file, $(wildcard $(EVENTS_DIR)/*), sam local invoke --env-vars env/local.json --event $(file);)
