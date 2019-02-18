package storage

import (
	"github.com/spf13/afero"
	"net/http"

	"chuck/utils"
)

type scRecorderImpl struct {
	root     *afero.Afero
	name     string
	recorder Recorder
	log      utils.Logger
}

func NewScenarioRecorder(fs afero.Fs, log utils.Logger, folder string, createNewFolder bool) (ScenarioRecorder, error) {
	name, path, err := utils.PrepareStorageFolder(fs, folder, createNewFolder)
	if err != nil {
		return nil, err
	}

	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, path)}
	return &scRecorderImpl{
		root: root,
		name: name,
		log:  log,
	}, nil
}

func (p *scRecorderImpl) Close() error {
	return nil
}

func (p *scRecorderImpl) Name() string {
	return p.name
}

func (self *scRecorderImpl) PendingCount() int {
	return 0
}

func (p *scRecorderImpl) ActivateScenario(name string) error {
	recorder, err := NewRecorder(p.root, p.log, name, false, true)
	if err != nil {
		return err
	}

	if p.recorder != nil {
		p.recorder.Close()
	}
	p.recorder = recorder

	return nil
}

func (p *scRecorderImpl) RecordRequest(req *http.Request, session int64) (*PendingRequest, error) {
	if p.recorder == nil {
		return nil, ErrScenarioRecorderNotActivated
	}

	return p.recorder.RecordRequest(req, session)
}

func (p *scRecorderImpl) RecordResponse(resp *http.Response, session int64) (*PendingRequest, error) {
	if p.recorder == nil {
		return nil, ErrScenarioRecorderNotActivated
	}

	return p.recorder.RecordResponse(resp, session)
}
