package handlers

import (
	"chuck/utils"
	"io"
	"net/http"

	"github.com/spf13/afero"
)

type scenarioSeekerNoProxyHandler struct {
	proxyHandler ProxyHandler
}

func NewScenarioSeekerNoProxyHandler(config *ScenarioSeekerConfig, fs afero.Fs, log utils.Logger) (http.Handler, error) {
	handler, err := NewScenarioSeekerHandler(config, fs, log)
	if err != nil {
		return nil, err
	}

	return &scenarioSeekerNoProxyHandler{proxyHandler: handler}, nil
}

//

func (self *scenarioSeekerNoProxyHandler) Request(req *http.Request) *http.Response {
	return self.proxyHandler.Request(req, nil)
}

func (self *scenarioSeekerNoProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sc := ParseActivateScenarioRequest(r)
	if sc != nil {
		self.proxyHandler.NonProxyHandler(w, r)
	} else {
		resp := self.proxyHandler.Request(r, nil)

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
		io.Copy(w, resp.Body)
		resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
	}
}
