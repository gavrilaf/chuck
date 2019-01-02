package handlers

import (
	"net/http"

	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
	"gopkg.in/elazarl/goproxy.v1"
)

type ProxyHandler interface {
	Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response
	Response(resp *http.Response, ctx *goproxy.ProxyCtx)

	NonProxyHandler(w http.ResponseWriter, req *http.Request)
}

// Factories

func NewRecorderHandler(config *RecorderConfig, fs afero.Fs, log utils.Logger) ProxyHandler {
	recorder, err := storage.NewRecorder(fs, log, config.Folder, config.CreateNewFolder, false)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return &recordHandler{
		recorder:       recorder,
		log:            log,
		preventCaching: config.Prevent304,
	}
}

func NewSeekerHandler(config *SeekerConfig, fs afero.Fs, log utils.Logger) ProxyHandler {
	seeker, err := storage.NewSeeker(fs, config.Folder)
	if err != nil {
		log.Panic("Could not create requests recorder: %v", err)
	}

	return &seekerHandler{
		seeker: seeker,
		log:    log,
	}
}

func NewScenarioSeekerHandler(config *ScenarioSeekerConfig, fs afero.Fs, log utils.Logger) ProxyHandler {
	seeker, err := storage.NewScenarioSeeker(fs, log, config.Folder)
	if err != nil {
		log.Panic("Could not create requests scenario seeker: %v", err)
	}

	return NewScenarioHandlerWithSeeker(seeker, log)
}

func NewScenarioRecorderHandler(config *ScenarioRecorderConfig, fs afero.Fs, log utils.Logger) ProxyHandler {
	recorder, err := storage.NewScenarioRecorder(fs, log, config.Folder, true)
	if err != nil {
		log.Panic("Could not create requests scenario recorder: %v", err)
	}

	return NewScenarioHandlerWithRecorder(recorder, log)
}
