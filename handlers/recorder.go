package handlers

import (
	"net/http"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

type recordHandler struct {
	recorder       storage.Recorder
	log            utils.Logger
	preventCaching bool
}

func (p *recordHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	_, err := p.recorder.RecordRequest(req, ctx.Session)
	if err != nil {
		p.log.Error("Record request error: %v", err)
	}

	if p.preventCaching {
		Prevent304HttpAnswer(req)
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
	p.log.Warn("*** Non-proxy request, %s : %s", req.Method, req.URL.String())
	w.WriteHeader(404)
	w.Write([]byte("Not supported in record mode"))
}
