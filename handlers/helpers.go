package handlers

import (
	"net/http"
)

var (
	ScenarioIdHeader = http.CanonicalHeaderKey("automation-test-identifier")
)

/*
 * Modify request headers top prevent 304 server answer
 */
func Prevent304HttpAnswer(req *http.Request) {
	req.Header.Set("If-Modified-Since", "off")
	req.Header.Set("Last-Modified", "")
}

/*
 * Return scenario id from header (or empty string)
 */
func GetScenarioId(req *http.Request) string {
	return req.Header.Get(ScenarioIdHeader)
}
