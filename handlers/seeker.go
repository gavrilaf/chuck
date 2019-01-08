package handlers

import (
	"net/http"
	"sync"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

type seekerHandler struct {
	seeker  storage.Seeker
	tracker storage.Tracker
	mux     *sync.Mutex
	log     utils.Logger
}

func (self *seekerHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	method, url := req.Method, req.URL.String()

	resp, err := self.seeker.Look(method, url)
	if err != nil {
		self.log.Error("Searching request error, %s : %s, (%v)", method, url, err)
	} else {
		if resp != nil {
			self.log.FocusedReq(req.Method, req.URL.String(), resp.StatusCode)
		} else {
			self.mux.Lock()
			self.tracker.RecordRequest(req, ctx.Session)
			self.mux.Unlock()
		}
	}
	return resp
}

func (self *seekerHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	self.mux.Lock()
	self.tracker.RecordResponse(resp, ctx.Session)
	self.mux.Unlock()
}

func (p *seekerHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("Not supported in debug mode"))
}
