package handlers

import (
	"net/http"
	"regexp"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

const (
	AADHIIdentifier = "aadhi-identifier"
)

type scenarioHandler struct {
	seeker    storage.ScSeeker
	log       utils.Logger
	scenarios map[string]string
}

func NewScenarioHandlerWithSeeker(seeker storage.ScSeeker, log utils.Logger) ProxyHandler {
	return &scenarioHandler{
		seeker:    seeker,
		log:       log,
		scenarios: make(map[string]string),
	}
}

func (p *scenarioHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	method := req.Method
	url := req.URL.String()
	id, ok := req.Header[http.CanonicalHeaderKey(AADHIIdentifier)]

	if ok {
		scenario, ok := p.scenarios[id[0]]
		if ok {
			resp := p.seeker.Look(scenario, method, url)
			if resp == nil {
				p.log.Error("Saved response isn't found for scenarion %s, %s : %s", scenario, method, url)
			} else {
				p.log.Info("Stubbed response for scenarion %s, request %s : %s", scenario, method, url)
				return resp
			}
		} else {
			p.log.Error("Scenario isn't found for id %s, %s : %s", id[0], method, url)
		}
	} else {
		p.log.Error("AADHI header not found for %s : %s", method, url)
	}
	return nil
}

func (p *scenarioHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {

}

func (p *scenarioHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	p.log.Info("NonProxyHandler, %s, %s", req.Method, req.URL.String())
	p.tryToActivateScenario(w, req)
}

////////////////////////////////////////////////////////////////////////////////////////

var activateRe = regexp.MustCompile("/scenario/(.*)/(.*)/no")

func (p *scenarioHandler) tryToActivateScenario(w http.ResponseWriter, req *http.Request) {
	url := req.URL.String()
	matches := activateRe.FindStringSubmatch(url)
	if len(matches) == 3 {
		scenario := matches[1]
		id := matches[2]

		if p.seeker.IsScenarioExists(scenario) {
			p.log.Info("Activated scenario %s with id %s", scenario, id)
			p.scenarios[id] = scenario
			w.WriteHeader(200)
			return
		} else {
			p.log.Error("Scenario %s not found", scenario)
		}
	}

	w.WriteHeader(404)
}
