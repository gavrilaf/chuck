package handlers

import (
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
)

type ProxyHandler interface {
	Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response
	Response(resp *http.Response, ctx *goproxy.ProxyCtx)

	NonProxyHandler(w http.ResponseWriter, req *http.Request)
}
