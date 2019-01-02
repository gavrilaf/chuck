package storage

import (
	"fmt"
	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
	"io"
	"net/http"
	"time"
)

var (
	ErrScenarioNotFound = fmt.Errorf("Scenario not found")
	ErrRequestNotFound  = fmt.Errorf("Request not found")
)

/*
 *
 */
type PendingRequest struct {
	Id      int64
	Method  string
	Url     string
	Started time.Time
}

type Tracker interface {
	RecordRequest(req *http.Request, session int64) (*PendingRequest, error)
	RecordResponse(resp *http.Response, session int64) (*PendingRequest, error)
	PendingCount() int
}

/*
 *
 */
type Recorder interface {
	Tracker
	io.Closer
	Name() string
	SetFocusedMode(focused bool)
}

func NewRecorder(folder string, createNewFolder bool, newOnly bool, log utils.Logger) (Recorder, error) {
	fs := afero.NewOsFs()
	return NewRecorderWithFs(fs, folder, createNewFolder, newOnly, log)
}

/*
 *
 */
type ScenarioRecorder interface {
	Tracker
	io.Closer
	Name() string
	ActivateScenario(name string) error
}

func NewScenarioRecorder(folder string, createNewFolder bool, log utils.Logger) (ScenarioRecorder, error) {
	fs := afero.NewOsFs()
	return NewScenarioRecorderWithFs(fs, folder, createNewFolder, log)
}

/*
 *
 */
type Seeker interface {
	Look(method string, url string) (*http.Response, error)
}

func NewSeeker(folder string) (Seeker, error) {
	fs := afero.NewOsFs()
	return NewSeekerWithFs(fs, folder)
}

/*
 *
 */
type ScenarioSeeker interface {
	IsScenarioExists(name string) bool
	Look(scenario string, method string, url string) (*http.Response, error)
}

func NewScenarioSeeker(folder string, log utils.Logger) (ScenarioSeeker, error) {
	fs := afero.NewOsFs()
	return NewScenarioSeekerWithFs(fs, folder, log)
}
