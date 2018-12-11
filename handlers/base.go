package handlers

import (
	"net/http"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"gopkg.in/elazarl/goproxy.v1"
)

type ProxyHandler interface {
	Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response
	Response(resp *http.Response, ctx *goproxy.ProxyCtx)

	NonProxyHandler(w http.ResponseWriter, req *http.Request)
}

func NewRecordHandler(folder string, log utils.Logger) ProxyHandler {
	recorder, err := storage.NewRecorder(folder, true, log)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return &recordHandler{
		recorder: recorder,
		log:      log,
	}
}

func NewSeekerHandler(folder string, log utils.Logger) ProxyHandler {
	seeker, err := storage.NewSeeker(folder, log)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return &seekerHandler{
		seeker: seeker,
		log:    log,
	}
}

func NewScenarioHandler(folder string, log utils.Logger) ProxyHandler {
	seeker, err := storage.NewScSeeker(folder, log)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return NewScenarioHandlerWithSeeker(seeker, log)
}
