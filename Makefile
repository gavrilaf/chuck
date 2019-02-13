SHELL := /bin/bash

test:
	go test ./utils -v
	go test ./storage -v
	go test ./handlers -v

run-rec:
	go build
	./chuck rec -address=127.0.0.1 -port=8123 -folder=log -prevent_304=1 -new_folder=1

run-dbg:
	go build
	./chuck dbg -address=127.0.0.1 -port=8123 -folder=dbg

run-intg:
	go build
	./chuck dbg -address=127.0.0.1 -port=8123 -folder=intg

run-intg-rec:
	go build
	./chuck intg_rec -address=127.0.0.1 -port=8123 -folder=log -new_folder=0

	