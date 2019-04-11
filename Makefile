SHELL := /bin/bash
export GO111MODULE=on

test:
	go test ./utils -v
	go test ./storage -v
	go test ./handlers -v

test-intg:
	go build
	./chuck intg -address=127.0.0.1 -port=8123 -folder=test-runner/stubs -verbose=1

build:
	go build

run-rec:
	./chuck rec -address=127.0.0.1 -port=8123 -folder=log -prevent_304=1 -new_folder=1 -focused=0 -requests=1 -filters=0

run-dbg:
	./chuck dbg -address=127.0.0.1 -port=8123 -folder=dbg

run-intg:
	./chuck intg -address=127.0.0.1 -port=8123 -folder=intg -verbose=0

run-intg-v:
	./chuck intg -address=127.0.0.1 -port=8123 -folder=intg -verbose=1

run-intg-rec:
	./chuck intg_rec -address=127.0.0.1 -port=8123 -folder=log-intg -new_folder=1 -requests=0 -filters=0

copy-intg-auto:
	source ext-tools/sc-copy/venv/bin/activate; \
	python ext-tools/sc-copy/main.py copy log-intg sc-cleaned auto; \
	deactivate; \

	