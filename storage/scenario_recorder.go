package storage

import (
	"github.com/spf13/afero"
	"net/http"

	"chuck/utils"
)

type scRecorderImpl struct {
	root        *afero.Afero
	name        string
	onlyNew     bool
	logRequests bool
	recorder    Recorder
	log         utils.Logger
}

func NewScenarioRecorder(fs afero.Fs, log utils.Logger, folder string, createNewFolder bool, onlyNew bool, logRequests bool) (ScenarioRecorder, error) {
	name, path, err := utils.PrepareStorageFolder(fs, folder, createNewFolder)
	if err != nil {
		return nil, err
	}

	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, path)}
	return &scRecorderImpl{
		root:        root,
		name:        name,
		onlyNew:     onlyNew,
		logRequests: logRequests,
		log:         log,
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

func (self *scRecorderImpl) ActivateScenario(name string) error {
	recorder, err := NewRecorder(self.root, self.log, name, false, self.onlyNew, self.logRequests)
	if err != nil {
		return err
	}

	if self.recorder != nil {
		self.recorder.Close()
	}

	recorder.SetFocusedMode(true)
	self.recorder = recorder

	return nil
}

func (self *scRecorderImpl) RecordRequest(req *http.Request, session int64) (*PendingRequest, error) {
	if self.recorder == nil {
		return nil, ErrScenarioRecorderNotActivated
	}

	return self.recorder.RecordRequest(req, session)
}

func (self *scRecorderImpl) RecordResponse(resp *http.Response, session int64) (*PendingRequest, error) {
	if self.recorder == nil {
		return nil, ErrScenarioRecorderNotActivated
	}

	return self.recorder.RecordResponse(resp, session)
}
