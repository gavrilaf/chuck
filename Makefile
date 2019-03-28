SHELL := /bin/bash
export GO111MODULE=on

test:
	go test ./utils -v
	go test ./storage -v
	go test ./handlers -v

test-intg:
	go build
	./chuck intg -address=127.0.0.1 -port=8123 -folder=test-runner/stubs

build:
	go build

run-rec:
	go build
	./chuck rec -address=127.0.0.1 -port=8123 -folder=log -prevent_304=1 -new_folder=1

run-dbg:
	go build
	./chuck dbg -address=127.0.0.1 -port=8123 -folder=dbg

run-intg:
	go build
	./chuck intg -address=127.0.0.1 -port=8123 -folder=intg

run-intg-rec:
	go build
	./chuck intg_rec -address=127.0.0.1 -port=8123 -folder=log-intg -new_folder=1

	