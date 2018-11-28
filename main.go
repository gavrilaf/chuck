package main

import (
	//"flag"
	"fmt"

	"gopkg.in/elazarl/goproxy.v1"
	//"io"
	"log"
	//"net"
	"crypto/tls"
	"net/http"

	"github.com/gavrilaf/chuck/storage"
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

	logger, err := storage.NewLogger()
	if err != nil {
		fmt.Printf("Could not create requests logger: %v\n", err)
		panic(err)
	}

	handler = NewHandler(logger)

	proxy := goproxy.NewProxyHttpServer()

	cert, err := tls.LoadX509KeyPair("ca.pem", "key.pem")
	if err != nil {
		log.Fatalf("Unable to load certificate - %v", err)
	}

	log.Printf("Cerfificates loaded ok\n")

	proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		return &goproxy.ConnectAction{
			Action:    goproxy.ConnectMitm,
			TLSConfig: goproxy.TLSConfigFromCA(&cert),
		}, host
	})

	proxy.OnRequest().DoFunc(handleRequest)
	proxy.OnResponse().DoFunc(handleResponse)

	log.Printf("Starting proxy\n")
	log.Fatal(http.ListenAndServe(addr, proxy))
}
