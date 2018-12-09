package handlers

import (
	"net/http"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

type seekerHandler struct {
	seeker storage.Seeker
	log    utils.Logger
}

func (p *seekerHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	method, url := req.Method, req.URL.String()
	resp := p.seeker.Look(method, url)
	if resp != nil {
		p.log.FocusedReq(req.Method, req.URL.String(), resp.StatusCode)
	}
	return resp
}

func (p *seekerHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {

}
