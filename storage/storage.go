package storage

import (
	"github.com/spf13/afero"
	"net/http"
)

/*
 * Logger
 */
type ReqLogger interface {
	Name() string
	SetFocusedMode(focused bool)
	LogRequest(req *http.Request, session int64) (int64, error)
	LogResponse(resp *http.Response, session int64) (int64, error)
	PendingCount() int
}

func NewLogger(folder string) (ReqLogger, error) {
	fs := afero.NewOsFs()
	return NewLoggerWithFs(folder, fs)
}

/*
 * Seeker
 */
type ReqSeeker interface {
	Look(method string, url string) *http.Response
}
