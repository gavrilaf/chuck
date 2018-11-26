package main

import (
	"fmt"
	"gopkg.in/elazarl/goproxy.v1"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gavrilaf/chuck/storage"
)

type pendingRequest struct {
	method  string
	url     string
	started time.Time
}

type proxyHandler struct {
	pending map[int64]pendingRequest
	mux     *sync.Mutex

	logger storage.ReqLogger
}

func (ph *proxyHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	ph.mux.Lock()
	defer ph.mux.Unlock()

	ph.pending[ctx.Session] = pendingRequest{
		method:  req.Method,
		url:     req.URL.String(),
		started: time.Now()}

	return nil
}

func (ph *proxyHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	ph.mux.Lock()
	defer ph.mux.Unlock()

	pending, ok := ph.pending[ctx.Session]
	if !ok {
		log.Fatalf("Could not find request for session: %d\n", ctx.Session)
	}

	elapsed := time.Since(pending.started)
	fmt.Printf("--> [%d] : [%v] %s %s, %v \n", ctx.Session, elapsed, ctx.Req.Method, ctx.Req.URL.String(), resp.Status)

	_, err := ph.logger.LogRequest(ctx.Req, resp)
	if err != nil {
		fmt.Printf("Could not log request: %v\n", err)
	}

	delete(ph.pending, ctx.Session)
}

func NewHandler(logger storage.ReqLogger) *proxyHandler {
	return &proxyHandler{
		pending: make(map[int64]pendingRequest, 10),
		mux:     &sync.Mutex{},
		logger:  logger,
	}
}
