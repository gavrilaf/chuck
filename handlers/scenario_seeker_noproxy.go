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
	rt := DetectServiceRequest(r)
	if rt != ServiceReq_None {
		self.proxyHandler.NonProxyHandler(w, r)
	} else {
		self.serverStubRequest(w, r)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////

func (self *scenarioSeekerNoProxyHandler) serverStubRequest(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	comps := strings.Split(url.Path, "/")

	// replace Chuck host on 'real' host. 'http://127.0.0.1/my.real.host/v1/profile' -> 'http://my.real.host/v1/profile'
	url.Host = comps[0]
	url.Path = strings.Join(comps[1:], "/")
	url.Scheme = "https"
	r.URL = url

	// handle request with standard proxy handler
	resp := self.proxyHandler.Request(r, nil)

	// copy result
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
	resp.Body.Close()
}
