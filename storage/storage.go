package storage

import (
	"github.com/spf13/afero"
	"net/http"
)

type ReqMeta struct {
	Req  *http.Request
	Resp *http.Response
}

type ReqLogger interface {
	Start() error
	Name() string
	SaveReqMeta(meta ReqMeta)
}

func NewLogger() ReqLogger {
	base := afero.NewOsFs()
	afr := &afero.Afero{Fs: afero.NewBasePathFs(base, "/log")}
	return NewLoggerWithFs(afr)
}

func NewLoggerWithFs(afr *afero.Afero) ReqLogger {
	return &reqLogger{
		base: afr,
	}
}
