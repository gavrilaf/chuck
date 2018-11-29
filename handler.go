package main

import (
	"github.com/gavrilaf/chuck/storage"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
)

type proxyHandler struct {
	logger storage.ReqLogger
}

func (ph *proxyHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	ph.logger.LogRequest(req, ctx.Session)
	return nil
}

func (ph *proxyHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	ph.logger.LogResponse(resp, ctx.Session)
}

func NewHandler(logger storage.ReqLogger) *proxyHandler {
	return &proxyHandler{
		logger: logger,
	}
}
