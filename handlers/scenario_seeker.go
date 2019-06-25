package handlers

import (
	"fmt"
	"net/http"

	"github.com/spf13/afero"
	"gopkg.in/elazarl/goproxy.v1"

	"chuck/storage"
	"chuck/utils"
)

type scenarioSeekerHandler struct {
	seeker    storage.ScenarioSeeker
	verbose   bool
	log       utils.Logger
	scenarios map[string]string
}

func NewScenarioSeekerHandler(config *ScenarioSeekerConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	seeker, err := storage.NewScenarioSeeker(fs, log, config.Folder)
	if err != nil {
		return nil, err
	}

	return &scenarioSeekerHandler{
		seeker:    seeker,
		verbose:   config.Verbose,
		log:       log,
		scenarios: make(map[string]string),
	}, nil
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

func (self *scenarioSeekerHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {}

func (self *scenarioSeekerHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	rt := DetectServiceRequest(req)
	switch rt {
	case ServiceReq_ActivateScenario:
		self.activateScenario(w, req)
	case ServiceReq_ExecuteScript:
		self.executeScript(w, req)
	default:
		fmt.Printf("Unsupported non proxy request: %v", req.URL.String())
		self.log.Error("Unsupported non proxy request: %v", req.URL.String())
		w.WriteHeader(404)
	}
}

////////////////////////////////////////////////////////////////////////////////////////

func (self *scenarioSeekerHandler) activateScenario(w http.ResponseWriter, req *http.Request) {
	sc := ParseActivateScenarioRequest(req)
	if sc == nil {
		self.log.Error("Wrong activate scenario request: %v", req.URL.String())
		w.WriteHeader(404)
		return
	}

	if self.seeker.IsScenarioExists(sc.Scenario) {
		self.log.Info("Activated scenario %s with id %s", sc.Scenario, sc.Id)
		self.scenarios[sc.Id] = sc.Scenario // TODO: fatal error: concurrent map writes
		w.WriteHeader(200)
	} else {
		self.log.Error("Scenario %s not found", sc.Scenario)
		w.WriteHeader(404)
	}
}

func (self *scenarioSeekerHandler) executeScript(w http.ResponseWriter, req *http.Request) {
	sc := ParseExecuteScriptRequest(req)
	if sc == nil {
		self.log.Error("Wrong execute script request: %v", req.URL.String())
		w.WriteHeader(404)
		return
	}

	utils.ExecuteCmd(sc.Name, sc.Env, self.log)
}
