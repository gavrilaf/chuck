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

// Factories

func NewRecorderHandler(config *RecorderConfig, log utils.Logger) ProxyHandler {
	recorder, err := storage.NewRecorder(config.Folder, config.CreateNewFolder, false, log)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return &recordHandler{
		recorder:       recorder,
		log:            log,
		preventCaching: config.Prevent304,
	}
}

func NewSeekerHandler(config *SeekerConfig, log utils.Logger) ProxyHandler {
	seeker, err := storage.NewSeeker(config.Folder)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return &seekerHandler{
		seeker: seeker,
		log:    log,
	}
}

func NewScenarioSeekerHandler(config *ScenarioSeekerConfig, log utils.Logger) ProxyHandler {
	seeker, err := storage.NewScenarioSeeker(config.Folder, log)
	if err != nil {
		log.Panic("Could not create requests scenario seeker: %v", err)
	}

	return NewScenarioHandlerWithSeeker(seeker, log)
}

func NewScenarioRecorderHandler(config *ScenarioRecorderConfig, log utils.Logger) ProxyHandler {
	recorder, err := storage.NewScenarioRecorder(config.Folder, true, log)
	if err != nil {
		log.Panic("Could not create requests scenario recorder: %v", err)
	}

	return NewScenarioHandlerWithRecorder(recorder, log)
}
