package handlers

import (
	"github.com/spf13/afero"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"

	"chuck/storage"
	"chuck/utils"
)

type scenarioRecordHandler struct {
	recorder       storage.ScenarioRecorder
	ch             chan scRecMeta
	log            utils.Logger
	scenarios      map[string]string
	preventCaching bool
}

type scRecMeta struct {
	activate bool
	writer   http.ResponseWriter
	req      *http.Request
	resp     *http.Response
	ctx      *goproxy.ProxyCtx
}

func NewScenarioRecorderHandler(config *ScenarioRecorderConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	recorder, err := storage.NewScenarioRecorder(fs, log, config.Folder, config.CreateNewFolder, config.OnlyNew, config.LogRequests)
	if err != nil {
		return nil, err
	}

	handler := &scenarioRecordHandler{
		recorder:       recorder,
		ch:             make(chan scRecMeta),
		log:            log,
		scenarios:      make(map[string]string),
		preventCaching: true,
	}

	go func() {
		for m := range handler.ch {
			if m.activate {
				handler.tryToActivateScenario(m.writer, m.req)
			} else if m.req != nil {
				_, err := handler.recorder.RecordRequest(m.req, m.ctx.Session)
				if err != nil {
					handler.log.Error("Record request error: %v", err)
				}
			} else if m.resp != nil {
				_, err := handler.recorder.RecordResponse(m.resp, m.ctx.Session)
				if err != nil {
					handler.log.Error("Record response error: %v", err)
				}
			}
		}
	}()

	return handler, nil
}

func (self *scenarioRecordHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	self.ch <- scRecMeta{activate: false, writer: nil, req: req, resp: nil, ctx: ctx}

	if self.preventCaching {
		Prevent304HttpAnswer(req)
	}

	return nil
}

func (self *scenarioRecordHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	self.ch <- scRecMeta{activate: false, writer: nil, req: nil, resp: resp, ctx: ctx}
}

func (self *scenarioRecordHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	self.ch <- scRecMeta{activate: true, writer: w, req: req, resp: nil, ctx: nil}
}

func (self *scenarioRecordHandler) tryToActivateScenario(w http.ResponseWriter, req *http.Request) {
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
