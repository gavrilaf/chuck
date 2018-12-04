package storage

import (
	"github.com/spf13/afero"
	"net/http"
)

/*
 */
type Recorder interface {
	Name() string
	SetFocusedMode(focused bool)
	RecordRequest(req *http.Request, session int64) (int64, error)
	RecordResponse(resp *http.Response, session int64) (int64, error)
	PendingCount() int
}

func NewRecorder(folder string) (Recorder, error) {
	fs := afero.NewOsFs()
	return NewRecorderWithFs(folder, fs)
}

/*
 */
type Seeker interface {
	Look(method string, url string) *http.Response
}

func NewSeeker(folder string) (Seeker, error) {
	fs := afero.NewOsFs()
	return NewSeekerWithFs(folder, fs)
}
