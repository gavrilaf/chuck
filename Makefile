SHELL := /bin/bash

run:
	go build
	./chuck

test:
	go test ./... -v


	
	
	