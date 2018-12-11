package cmds

import (
	"crypto/tls"
	. "github.com/gavrilaf/chuck/handlers"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
)

func CreateProxy() (*goproxy.ProxyHttpServer, error) {
	cert, err := tls.LoadX509KeyPair("ca.pem", "key.pem")
	if err != nil {
		return nil, err
	}

	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return &goproxy.ConnectAction{
			Action:    goproxy.ConnectMitm,
			TLSConfig: goproxy.TLSConfigFromCA(&cert),
		}, host
	})

	return proxy, nil
}

func RunProxy(proxy *goproxy.ProxyHttpServer, handler ProxyHandler, addr string) error {
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler.NonProxyHandler(w, req)
	})

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		resp := handler.Request(req, ctx)
		return req, resp
	})

	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		handler.Response(resp, ctx)
		return resp
	})

	return http.ListenAndServe(addr, proxy)
}
