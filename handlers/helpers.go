package handlers

import (
	"net/http"
	"regexp"
)

var (
	ActivateScenarioRegx = regexp.MustCompile("/scenario/(.*)/(.*)/no")
	ScenarioIdHeader     = http.CanonicalHeaderKey("aadhi-identifier")
)

/*
 * Modify request headers top prevent 304 server answer
 */
func Prevent304HttpAnswer(req *http.Request) {
	req.Header.Set("If-Modified-Since", "off")
	req.Header.Set("Last-Modified", "")
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
 * Return scenario id from header (or empty string)
 */

func GetScenarioId(req *http.Request) string {
	return req.Header.Get(ScenarioIdHeader)
}
