package handlers

import (
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

	if method == "POST" {
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
