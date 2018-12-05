package main

import (
	"crypto/tls"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
)

var handler *proxyHandler

func handleRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	resp := handler.Request(req, ctx)
	return req, resp
}

func handleResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	handler.Response(resp, ctx)
	return resp
}

func main() {
	addr := ":8123"

	log := utils.NewLogger()

	handler = NewHandler(log)

	proxy := goproxy.NewProxyHttpServer()

	cert, err := tls.LoadX509KeyPair("ca.pem", "key.pem")
	if err != nil {
		log.Panic("Unable to load certificate: %v", err)
	}

	log.Info("Cerfificates loaded ok")

	proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return &goproxy.ConnectAction{
			Action:    goproxy.ConnectMitm,
			TLSConfig: goproxy.TLSConfigFromCA(&cert),
		}, host
	})

	proxy.OnRequest().DoFunc(handleRequest)
	proxy.OnResponse().DoFunc(handleResponse)

	log.Info("Starting proxy")
	err = http.ListenAndServe(addr, proxy)
	if err != nil {
		log.Panic("Couldn't start proxy: %v", err)
	}
}
