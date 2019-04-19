package handlers

import (
	"github.com/spf13/afero"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"

	"chuck/storage"
	"chuck/utils"
)

type recordHandler struct {
	recorder       storage.Recorder
	log            utils.Logger
	preventCaching bool
}

func NewRecorderHandler(config *RecorderConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	recorder, err := storage.NewRecorder(fs, log, config.Folder, config.CreateNewFolder, config.OnlyNew, config.LogRequests)
	if err != nil {
		return nil, err
	}

	handler := &recordHandler{
		recorder:       recorder,
		log:            log,
		preventCaching: config.Prevent304,
	}

	return handler, nil
}

func (self *recordHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	_, err := self.recorder.RecordRequest(req, ctx.Session)
	if err != nil {
		self.log.Error("Record request error: %v", err)
	}

	if self.preventCaching {
		Prevent304HttpAnswer(req)
	}

	return nil
}

func (self *recordHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	_, err := self.recorder.RecordResponse(resp, ctx.Session)
	if err != nil {
		self.log.Error("Record response error: %v", err)
	}
}

func (self *recordHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	self.log.Warn("*** Non-proxy request, %s : %s", req.Method, req.URL.String())
	w.WriteHeader(404)
	w.Write([]byte("Not supported in record mode"))
}
