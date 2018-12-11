package storage

import (
	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
	"net/http"
)

/*
 *
 */
type Recorder interface {
	Name() string
	SetFocusedMode(focused bool)
	RecordRequest(req *http.Request, session int64) (int64, error)
	RecordResponse(resp *http.Response, session int64) (int64, error)
	PendingCount() int
}

func NewRecorder(folder string, createNewFolder bool, log utils.Logger) (Recorder, error) {
	fs := afero.NewOsFs()
	return NewRecorderWithFs(folder, createNewFolder, fs, log)
}

/*
 *
 */
type Seeker interface {
	Look(method string, url string) *http.Response
}

func NewSeeker(folder string, log utils.Logger) (Seeker, error) {
	fs := afero.NewOsFs()
	return NewSeekerWithFs(folder, fs, log)
}

/*
 *
 */
type ScSeeker interface {
	Look(scenario string, method string, url string) *http.Response
	IsScenarioExists(name string) bool
}

func NewScSeeker(folder string, log utils.Logger) (ScSeeker, error) {
	fs := afero.NewOsFs()
	return NewScSeekerWithFs(folder, fs, log)
}
