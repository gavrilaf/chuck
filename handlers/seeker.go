package handlers

import (
	"github.com/spf13/afero"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
	"sync"

	"chuck/storage"
	"chuck/utils"
)

type seekerHandler struct {
	seeker  storage.Seeker
	tracker storage.Tracker
	mux     *sync.Mutex
	log     utils.Logger
}

func NewSeekerHandler(config *SeekerConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	seeker, err := storage.NewSeeker(fs, config.Folder, log)
	if err != nil {
		return nil, err
	}

	return &seekerHandler{
		seeker:  seeker,
		tracker: storage.NewTracker(0, log),
		mux:     &sync.Mutex{},
		log:     log,
	}, nil
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
