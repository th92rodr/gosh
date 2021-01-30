GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_RUN=$(GO_CMD) run
GO_CLEAN=$(GO_CMD) clean

BINARY_NAME=simplesh

help:
	@sed -ne '/@sed/!s/##//p' $(MAKEFILE_LIST)
build: ##
	$(GO_BUILD) -o $(BINARY_NAME) -v
run: ##
	$(GO_RUN) main.go
clean: ##
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
docker-run: ##
	docker run --name $(BINARY_NAME) -v $(shell pwd):/go/src/$(BINARY_NAME) --rm -it golang:1.15.3-buster /bin/bash
