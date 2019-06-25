package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

var (
	ActivateScenarioRegx = regexp.MustCompile("/scenario/(.*)/(.*)/no")
	ExecuteScriptRegx    = regexp.MustCompile("/script/(.*)/run")
)

const (
	ServiceReq_None = iota
	ServiceReq_ActivateScenario
	ServiceReq_ExecuteScript
)

func DetectServiceRequest(req *http.Request) int {
	method := req.Method
	url := req.URL.String()

	if method == "PUT" {
		switch {
		case ActivateScenarioRegx.MatchString(url):
			return ServiceReq_ActivateScenario
		case ExecuteScriptRegx.MatchString(url):
			return ServiceReq_ExecuteScript
		}
	}

	return ServiceReq_None
}

/*
 * Parse request url and return scenario name & id if url is recognized as scenario activation url
 */

type ActivateScenario struct {
	Scenario string
	Id       string
}

func ParseActivateScenarioRequest(req *http.Request) *ActivateScenario {
	url := req.URL.String()
	matches := ActivateScenarioRegx.FindStringSubmatch(url)
	if len(matches) == 3 {
		return &ActivateScenario{Scenario: matches[1], Id: matches[2]}
	}
	return nil
}

/*
 * Parse request url and return script name url is recognized as script execute url
 */

type ExecuteScript struct {
	Name string
	Env  map[string]string
}

func ParseExecuteScriptRequest(req *http.Request) *ExecuteScript {
	url := req.URL.String()
	matches := ExecuteScriptRegx.FindStringSubmatch(url)
	if len(matches) == 2 {
		var env map[string]string

		body, err := ioutil.ReadAll(req.Body)
		if err == nil && len(body) > 0 {
			json.Unmarshal(body, &env)
		}

		return &ExecuteScript{Name: matches[1], Env: env}
	}
	return nil
}
