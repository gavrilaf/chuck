# chuck

The proxy server for mobile application debugging, mocking nonexisting API, record and run integration tests.

## Installation

1. Install go
MacOS using brew: *brew install go*
Other operation systems: https://golang.org/doc/install

2. Clone chuck repository: *git clone https://github.com/gavrilaf/chuck.git*

Chuck uses golang modules so you don't need to clone Chuck repo to the GOPATH directory. Clone Chuck anywhere. But if you clone Chuck to the GOPATH folder then setup GO111MODULE=on environment variable. Or build and run Chuck using makefile as described below.

3. Load dependencies, build & check chuck: *make test* from the Chuck folder
This command build the project and run all unit tests. All tests should be passed.

## How to use

Before run you have to install certificate from the project folder as root certificate (or generate and install new one).

**Chuck** supports 4 modes:

### Record mode

*chuck rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304] [-focused]*
In this mode, the application is working as a transparent proxy but all traffic is saved to the disk.

### Debug mode

Just copy logs recorded in the rec mode to another folder. Make some requests **focused** (changed 'N' on 'F' in the index.txtx file), edit according responses and run **Chuck** in the debug mode.

*chuck dbg [-address=addr] [-port=port] [-folder=folder]*

**Chuck** runs as a transparent proxy but for the focused requests will be returned stored values. Also, you are able to add new requests to the index.txt. It's a very simple way how to stub unexisting REST API.

### Integration tests recording mode

Almost the same with recording mode but with scenarios support. 

*chuck intg_rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304]*

The tested application should call endpoint *https://chuck-url/scenario/scenario_id/app_id/no* before any scenario execution. All following requests will be stored to the folder 'scenario_id'. app_id isn't used right now.

Example
The application generates requests.
```
https://chuck-url/scenario/sc-12/123/no
https://some-endpoint-1
https://some-endpoint-2
https://chuck-url/scenario/sc-674/123/no
https://some-endpoint-1
https://some-endpoint-4
```

Chuck creates following log.

```
root
  |----sc-12
  |       |-----index.txt (some-endpoint-1, some-endpoint-2)
  |
  |----sc-674
  |       |-----index.txt (some-endpoint-1, some-endpoint-4)
  ```
  All requests are logged as focused.


### Integration tests playing mode

*chuck intg_rec [-address=addr] [-port=port] [-folder=folder]*

After recording scenarios using recording mode **Chuck** is able to play it. The tested application should call endpoint *https://chuck-url/scenario/scenario_id/app_id/no* before any scenario execution. Each request should contains *'int-test-identifier = app_id'* in the http header. **Chuck** uses this header to determine which scenario is active for this application.

### Using make

*make run-rec/run-dgb/run-intg/run-intg-rec* is running **Сhuck** in different modes with default parameters.

### How to generate new root certificates

### Install self signed certificates on iOS simulator
1. Launch Safari in the simulator.
2. Drag’n’Drop ‘ca.pem’ file into the browser address line.
3. Install certificate (press "Allow"/"Done"/"Install" on all questions)
4. Navigate to Settings > General > About > Certificate Trust Settings
5. Enable full trust for the new installed certificate

## Developing

I hope **Chuck** is a good example how to write well-formed Golang code and how to use BDD with Golang.
