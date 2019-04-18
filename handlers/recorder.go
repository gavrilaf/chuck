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
	ch             chan recordMeta
	log            utils.Logger
	preventCaching bool
}

type recordMeta struct {
	req  *http.Request
	resp *http.Response
	ctx  *goproxy.ProxyCtx
}

func NewRecorderHandler(config *RecorderConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	recorder, err := storage.NewRecorder(fs, log, config.Folder, config.CreateNewFolder, config.OnlyNew, config.LogRequests)
	if err != nil {
		return nil, err
	}

	handler := &recordHandler{
		recorder: recorder,
		ch:       make(chan recordMeta),
		log:      log,

		preventCaching: config.Prevent304,
	}

	go func() {
		for m := range handler.ch {
			switch {
			case m.req != nil:
				_, err := handler.recorder.RecordRequest(m.req, m.ctx.Session)
				if err != nil {
					handler.log.Error("Record request error: %v", err)
				}
			case m.resp != nil:
			}

		}
	}()

	return handler, nil
}

func (self *recordHandler) Request(req *http.Request, ctx *goproxy.ProxyCtx) *http.Response {
	self.ch <- recordMeta{req: req, resp: nil, ctx: ctx}

	if self.preventCaching {
		Prevent304HttpAnswer(req)
	}

	return nil
}

func (self *recordHandler) Response(resp *http.Response, ctx *goproxy.ProxyCtx) {
	self.ch <- recordMeta{req: nil, resp: resp, ctx: ctx}
}

func (self *recordHandler) NonProxyHandler(w http.ResponseWriter, req *http.Request) {
	self.log.Warn("*** Non-proxy request, %s : %s", req.Method, req.URL.String())
	w.WriteHeader(404)
	w.Write([]byte("Not supported in record mode"))
}
