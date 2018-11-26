package storage

import (
	"net/http"
)

type ReqLogger interface {
	Name() string
	LogRequest(req *http.Request, resp *http.Response) (string, error)
}
