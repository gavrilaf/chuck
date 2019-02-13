package handlers

import (
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"

	"chuck/storage"
	"chuck/utils"
)

type scenarioRecordHandler struct {
	recorder       storage.ScenarioRecorder
	log            utils.Logger
	scenarios      map[string]string
	preventCaching bool
}

func (p *scenarioRecordHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	_, err := p.recorder.RecordRequest(req, ctx.Session)
	if err != nil {
		p.log.Error("Record request error: %v", err)
	}

	if p.preventCaching {
		Prevent304HttpAnswer(req)
	}
	return nil
}

func (p *scenarioRecordHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	_, err := p.recorder.RecordResponse(resp, ctx.Session)
	if err != nil {
		p.log.Error("Record response error: %v", err)
	}
}

func (p *scenarioRecordHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	p.tryToActivateScenario(w, req)
}

func (p *scenarioRecordHandler) tryToActivateScenario(w http.ResponseWriter, req *http.Request) {
	sc := ParseActivateScenarioRequest(req)
	if sc == nil {
		return
	}

	err := p.recorder.ActivateScenario(sc.Scenario)
	if err != nil {
		p.log.Error("Couldn't activate scenario %s, %v", sc.Scenario, err)
	} else {
		p.log.Info("Activated scenario %s with id %s", sc.Scenario, sc.Id)
		p.scenarios[sc.Id] = sc.Scenario
		w.WriteHeader(200)
		return
	}

	w.WriteHeader(404)
}
