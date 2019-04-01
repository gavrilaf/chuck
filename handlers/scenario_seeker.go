package handlers

import (
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"

	"chuck/storage"
	"chuck/utils"
)

type scenarioSeekerHandler struct {
	seeker    storage.ScenarioSeeker
	verbose   bool
	log       utils.Logger
	scenarios map[string]string
}

func (self *scenarioSeekerHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	method := req.Method
	url := req.URL.String()
	id := GetScenarioId(req)

	if len(id) != 0 {
		scenario, ok := self.scenarios[id]
		if ok {
			resp, err := self.seeker.Look(scenario, method, url)
			if err != nil {
				self.log.Error("Searching response error %v, %s, %s : %s, (%v)", id, scenario, method, url, err)
			} else if resp == nil {
				if self.verbose {
					self.log.Warn("Saved response isn't found for client %v, scenario %s, %s : %s", id, scenario, method, url)
				}
			} else {
				if self.verbose {
					self.log.Info("Stubbed response for client %v, scenario %s, request %s : %s", id, scenario, method, url)
				}
				return resp
			}
		} else {
			self.log.Error("Scenario isn't found for id %v, %s : %s", id, method, url)
		}
	} else {
		self.log.Error("Integration test header not found for %s : %s", method, url)
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
