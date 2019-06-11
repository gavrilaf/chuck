SHELL := /bin/bash
export GO111MODULE=on

define PROJECT_HELP_MSG
Usage:\n
    \t make help:\t\t			show this message\n\n

    \t make build:\t\t			build Chuck\n
    \t make install-tools:\t	install additional tools requirements\n\n

	\t make test:\t\t			run Chuck unit tests\n
    \t make test-intg-c:\t		run Chuck for the integration tests (Chuck integration tests)\n
    \t make test-intg-r:\t		run Chuck integration tests (exec make test-intg-c before)\n\n

    \t make run-rec:\t\t		run Chuck in the recording mode\n
	\t make run-dbg:\t\t		run Chuck in the debug mode\n

	\t make run-intg-rec:\t		run Chuck in the integration tests recording mode\n\n
	\t make run-intg:\t			run Chuck in the integration tests mode as local server\n\n
	\t make run-intg-r:\t		run Chuck in the integration tests mode for real devices(not simulator)\n\n

	\t make copy-intg-auto:\t	run integration tests copying & cleaning utility (auto mode)\n
endef
export PROJECT_HELP_MSG

help:
	echo -e $$PROJECT_HELP_MSG

build:
	go build

install-tools:
	python3 -m venv ext-tools/test-runner/venv
	source ext-tools/test-runner/venv/bin/activate && pip install --upgrade pip && pip install -r ext-tools/test-runner/requirements.txt
	python3 -m venv ext-tools/sc-copy/venv
	source ext-tools/sc-copy/venv/bin/activate && pip install --upgrade pip && pip install -r ext-tools/sc-copy/requirements.txt

test:
	go test ./... -v

test-intg-c:
	go build
	./chuck intg -address=127.0.0.1 -port=8123 -folder=ext-tools/test-runner/stubs -verbose=1

test-intg-r:
	( \
		source ext-tools/test-runner/venv/bin/activate; \
		python ext-tools/test-runner/test_chuck.py; \
	)

run-rec:
	./chuck rec -address=127.0.0.1 -port=8123 -folder=log -prevent_304=1 -new_folder=1 -focused=0 -requests=1

run-dbg:
	./chuck dbg -address=127.0.0.1 -port=8123 -folder=dbg

run-intg:
	./chuck intg-noproxy -address=127.0.0.1 -port=8123 -folder=intg -verbose=1

run-intg-rec:
	./chuck intg_rec -address=127.0.0.1 -port=8123 -folder=log-intg -new_folder=1 -requests=0

run-intg-r:
	./chuck intg-noproxy -address=0.0.0.0 -port=8123 -folder=intg -verbose=1

copy-intg-auto: 
	( \
		source ext-tools/sc-copy/venv/bin/activate; \
		python ext-tools/sc-copy/main.py copy log-intg sc-cleaned auto; \
	)

	