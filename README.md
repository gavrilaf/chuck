# chuck

The proxy server for mobile application debugging, mocking nonexisting API, record and run integration tests.

Support 4 modes:

### Record mode

*chuck rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304] [-focused]*
In this mode, the application is working as a transparent proxy but all traffic is saved to the disk.

### Debug mode

Just copy logs recorded in the rec mode to another folder. Make some requests focused (changed 'N' on 'F' in the index.txtx file), edit according responses and run **Chuck** in the debug mode.

*chuck dbg [-address=addr] [-port=port] [-folder=folder]*

Chuck runs as a transparent proxy but for the focused requests will be returned stored values. Also, you are able to add new requests to the index.txt. It's a very simple way how to stub unexisting REST API.

### Integration tests recording mode

Almost the same with recording mode but with scenarios support. 

*chuck intg_rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304]*

The tested application should call endpoint https://chuck-url/scenario/scenario_id/app_id/no before any scenario execution. All following requests will be stored to the folder 'scenario_id'. app_id isn't used right now.

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
