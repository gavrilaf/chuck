package main

import (
	"net/http"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

type proxyHandler struct {
	recorder storage.Recorder
	log      utils.Logger
}

func (ph *proxyHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	_, err := ph.recorder.RecordRequest(req, ctx.Session)
	if err != nil {
		ph.log.Error("Record request error: %v", err)
	}

	return nil
}

func (ph *proxyHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	_, err := ph.recorder.RecordResponse(resp, ctx.Session)
	if err != nil {
		ph.log.Error("Record response error: %v", err)
	}
}

func NewHandler(log utils.Logger) *proxyHandler {
	recorder, err := storage.NewRecorder("", log)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return &proxyHandler{
		recorder: recorder,
		log:      log,
	}
}
