package main

import (
	"net/http"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

type ProxyHandler interface {
	Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response
	Response(resp *http.Response, ctx *goproxy.ProxyCtx)
}

func NewRecordHandler(folder string, log utils.Logger) *recordHandler {
	recorder, err := storage.NewRecorder(folder, log)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return &recordHandler{
		recorder: recorder,
		log:      log,
	}
}

func NewSeekerHandler(folder string, log utils.Logger) *seekerHandler {
	seeker, err := storage.NewSeeker(folder, log)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return &seekerHandler{
		seeker: seeker,
		log:    log,
	}
}

/////////////////////////////////////////////////////////////////////////

type recordHandler struct {
	recorder storage.Recorder
	log      utils.Logger
}

func (p *recordHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	_, err := p.recorder.RecordRequest(req, ctx.Session)
	if err != nil {
		p.log.Error("Record request error: %v", err)
	}

	return nil
}

func (p *recordHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	_, err := p.recorder.RecordResponse(resp, ctx.Session)
	if err != nil {
		p.log.Error("Record response error: %v", err)
	}
}

/////////////////////////////////////////////////////////////////////////

type seekerHandler struct {
	seeker storage.Seeker
	log    utils.Logger
}

func (p *seekerHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	method, url := req.Method, req.URL.String()
	resp := p.seeker.Look(method, url)
	if resp != nil {
		p.log.FocusedReq(req.Method, req.URL.String(), resp.StatusCode)
	}
	return resp
}

func (p *seekerHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {

}
