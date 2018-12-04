package main

import (
	"github.com/gavrilaf/chuck/storage"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
)

type proxyHandler struct {
	recorder storage.Recorder
}

func (ph *proxyHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	ph.recorder.RecordRequest(req, ctx.Session)
	return nil
}

func (ph *proxyHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	ph.recorder.RecordResponse(resp, ctx.Session)
}

func NewHandler(recorder storage.Recorder) *proxyHandler {
	return &proxyHandler{
		recorder: recorder,
	}
}
