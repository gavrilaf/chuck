package handlers

import (
	"chuck/utils"
	"io"
	"net/http"
	"strings"

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

///////////////////////////////////////////////////////////////////////////////////////////////////////

func (self *scenarioSeekerNoProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sc := ParseActivateScenarioRequest(r)
	if sc != nil {
		self.proxyHandler.NonProxyHandler(w, r)
	} else {
		url := r.URL

		comps := strings.Split(url.Path, "/")

		url.Host = comps[0]
		url.Path = strings.Join(comps[1:], "/")
		url.Scheme = "https"

		r.URL = url

		resp := self.proxyHandler.Request(r, nil)

		for k, vv := range resp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(resp.StatusCode)

		io.Copy(w, resp.Body)
		resp.Body.Close()
	}
}
