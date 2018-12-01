package storage

import (
	"net/http"
)

type ReqLogger interface {
	Name() string
	LogRequest(req *http.Request, session int64) (int64, error)
	LogResponse(resp *http.Response, session int64) (int64, error)
	PendingCount() int
}

type ReqSeeker interface {
	Look(method string, url string) *http.Response
}
