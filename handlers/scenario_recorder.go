package handlers

import (
	"net/http"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

type scenarioRecordHandler struct {
	recorder storage.ScenarioRecorder
	log      utils.Logger

	scenarios      map[string]string
	preventCaching bool
}

func NewScenarioHandlerWithRecorder(recorder storage.ScenarioRecorder, log utils.Logger) ProxyHandler {
	return &scenarioRecordHandler{
		recorder:       recorder,
		log:            log,
		scenarios:      make(map[string]string),
		preventCaching: true,
	}
}

func (p *scenarioRecordHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	// TODO: Should check id in header to define scenario, will add later

	_, err := p.recorder.RecordRequest(req, ctx.Session)
	if err != nil {
		p.log.Error("Record request error: %v", err)
	}

	if p.preventCaching {
		p.prevent304(req)
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

func (p *scenarioRecordHandler) prevent304(req *http.Request) {
	req.Header.Set("If-Modified-Since", "off")
	req.Header.Set("Last-Modified", "")
}

func (p *scenarioRecordHandler) tryToActivateScenario(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	matches := activateScRegx.FindStringSubmatch(url)
	if len(matches) == 3 {
		scenario := matches[1]
		id := matches[2]

		err := p.recorder.ActivateScenario(scenario)
		if err != nil {
			p.log.Error("Couldn't activate scenario %s, %v", scenario, err)
		} else {
			p.log.Info("Activated scenario %s with id %s", scenario, id)
			p.scenarios[id] = scenario
			w.WriteHeader(200)
			return
		}
	}

	w.WriteHeader(404)
}
