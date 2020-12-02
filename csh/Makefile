BINARY_NAME=csh

help:
	@sed -ne '/@sed/!s/##//p' $(MAKEFILE_LIST)
build: ##
	gcc -g -lreadline -o $(BINARY_NAME) main.c
run: ##
	gcc -g -lreadline -o $(BINARY_NAME) main.c
	./$(BINARY_NAME)
clean: ##
	rm -f $(BINARY_NAME)
deps: ##
	apt-get update && apt-get install libreadline-dev -y
