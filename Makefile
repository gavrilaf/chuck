SHELL := /bin/bash

run-rec:
	go build
	./chuck rec

test:
	go test ./... -v


	
	
	