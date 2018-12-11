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
	activateResp := p.tryToActivateScenario(req)
	if activateResp != nil {
		return activateResp
	}

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

////////////////////////////////////////////////////////////////////////////////////////

var activateRe = regexp.MustCompile("/scenario/(.*)/(.*)/no")

func (p *scenarioHandler) tryToActivateScenario(req *http.Request) *http.Response {
	url := req.URL.String()
	matches := activateRe.FindStringSubmatch(url)
	if len(matches) == 3 {
		scenario := matches[1]
		id := matches[2]

		var statusCode int
		if p.seeker.IsScenarioExists(scenario) {
			p.log.Info("Activated scenario %s with id %s", scenario, id)
			p.scenarios[id] = scenario
			statusCode = 200
		} else {
			p.log.Error("Scenario %s not found", scenario)
			statusCode = 404
		}

		return &http.Response{
			StatusCode: statusCode,
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
	}

	return nil
}
