package handlers

import (
	"net/http"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

type recordHandler struct {
	recorder storage.Recorder
	log      utils.Logger
}

func (p *recordHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	req.Header.Set("Cache-Control", "public, max-age=80000, immutable")
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

func (p *recordHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("Not supported in record mode"))
}
