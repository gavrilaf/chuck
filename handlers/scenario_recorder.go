package handlers

import (
	"github.com/spf13/afero"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
	"sync"

	"chuck/storage"
	"chuck/utils"
)

type scenarioRecordHandler struct {
	recorder       storage.ScenarioRecorder
	mux            *sync.Mutex
	log            utils.Logger
	scenarios      map[string]string
	preventCaching bool
}

func NewScenarioRecorderHandler(config *ScenarioRecorderConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	recorder, err := storage.NewScenarioRecorder(fs, log, config.Folder, config.CreateNewFolder, config.OnlyNew, config.LogRequests)
	if err != nil {
		return nil, err
	}

	handler := &scenarioRecordHandler{
		recorder:       recorder,
		mux:            &sync.Mutex{},
		log:            log,
		scenarios:      make(map[string]string),
		preventCaching: true,
	}

	return handler, nil
}

func (self *scenarioRecordHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	self.mux.Lock()
	defer self.mux.Unlock()

	_, err := self.recorder.RecordRequest(req, ctx.Session)
	if err != nil {
		self.log.Error("Record request error: %v", err)
	}

	if self.preventCaching {
		Prevent304HttpAnswer(req)
	}

	return nil
}

func (self *scenarioRecordHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	self.mux.Lock()
	defer self.mux.Unlock()

	_, err := self.recorder.RecordResponse(resp, ctx.Session)
	if err != nil {
		self.log.Error("Record response error: %v", err)
	}
}

func (self *scenarioRecordHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	self.tryToActivateScenario(w, req)
}

func (self *scenarioRecordHandler) tryToActivateScenario(w http.ResponseWriter, req *http.Request) {
	self.mux.Lock()
	defer self.mux.Unlock()

	sc := ParseActivateScenarioRequest(req)
	if sc == nil {
		return
	}

	err := self.recorder.ActivateScenario(sc.Scenario)
	if err != nil {
		self.log.Error("Couldn't activate scenario %s, %v", sc.Scenario, err)
	} else {
		self.log.Info("Activated scenario %s with id %s", sc.Scenario, sc.Id)
		self.scenarios[sc.Id] = sc.Scenario
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(404)
}
