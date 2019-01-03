# chuck

The proxy server for mobile application debugging, mocking nonexisting API, record and run integration tests.

Support 4 modes:

### Record mode

*chuck rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304] [-focused]*
In this mode application is working as transparent proxy but all traffic is saving on the disk.

### Debug mode

Just copy logs recorded in the rec mode to another folder. Make some requests focused (changed 'N' on 'F' in the index.txtx file), edit according responses and run **Chuck** in the debug mode.

*chuck dbg [-address=addr] [-port=port] [-folder=folder]*

Chuck will run as a transparent proxy but for the focused requests will be returned stored values. Also you are able to add new requests to the index.txt.

### Integration tests recording mode

*chuck intg_rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304]*

### Integration tests playing mode

*chuck intg_rec [-address=addr] [-port=port] [-folder=folder]*
