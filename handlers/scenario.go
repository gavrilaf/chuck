package handlers

import (
	"net/http"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

const (
	AADHIIdentifier = "aadhi-identifier"
)

type scenarioHandler struct {
	seeker    storage.ScenarioSeeker
	log       utils.Logger
	scenarios map[string]string
}

func NewScenarioHandlerWithSeeker(seeker storage.ScenarioSeeker, log utils.Logger) ProxyHandler {
	return &scenarioHandler{
		seeker:    seeker,
		log:       log,
		scenarios: make(map[string]string),
	}
}

func (p *scenarioHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	method := req.Method
	url := req.URL.String()
	id := GetScenarioId(req)

	if len(id) != 0 {
		scenario, ok := p.scenarios[id]
		if ok {
			resp, err := p.seeker.Look(scenario, method, url)
			if err != nil {
				p.log.Warn("Searching response error %s, %s : %s, (%v)", scenario, method, url, err)
			} else if resp == nil {
				p.log.Warn("Saved response isn't found for scenarion %s, %s : %s", scenario, method, url)
			} else {
				p.log.Info("Stubbed response for scenarion %s, request %s : %s", scenario, method, url)
				return resp
			}
		} else {
			p.log.Error("Scenario isn't found for id %s, %s : %s", id[0], method, url)
		}
	} else {
		p.log.Error("Integration test header not found for %s : %s", method, url)
	}
	return nil
}

func (p *scenarioHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {

}

func (p *scenarioHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	p.tryToActivateScenario(w, req)
}

////////////////////////////////////////////////////////////////////////////////////////

func (p *scenarioHandler) tryToActivateScenario(w http.ResponseWriter, req *http.Request) {
	sc := ParseActivateScenarioRequest(req)
	if sc == nil {
		return
	}

	if p.seeker.IsScenarioExists(sc.Scenario) {
		p.log.Info("Activated scenario %s with id %s", sc.Scenario, sc.Id)
		p.scenarios[sc.Id] = sc.Scenario
		w.WriteHeader(200)
		return
	} else {
		p.log.Error("Scenario %s not found", sc.Scenario)
	}

	w.WriteHeader(404)
}
