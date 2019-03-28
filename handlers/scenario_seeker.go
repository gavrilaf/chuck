package handlers

import (
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"

	"chuck/storage"
	"chuck/utils"
)

type scenarioSeekerHandler struct {
	seeker    storage.ScenarioSeeker
	log       utils.Logger
	scenarios map[string]string
}

func (p *scenarioSeekerHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
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
			p.log.Error("Scenario isn't found for id %v, %s : %s", id, method, url)
		}
	} else {
		p.log.Error("Integration test header not found for %s : %s", method, url)
	}
	return utils.MakeResponse2(404, make(http.Header), "")
}

func (p *scenarioSeekerHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {

}

func (p *scenarioSeekerHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	p.tryToActivateScenario(w, req)
}

////////////////////////////////////////////////////////////////////////////////////////

func (p *scenarioSeekerHandler) tryToActivateScenario(w http.ResponseWriter, req *http.Request) {
	sc := ParseActivateScenarioRequest(req)
	if sc == nil {
		p.log.Error("Wrong activate scenario request: %v", req.URL.String())
		w.WriteHeader(404)
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
