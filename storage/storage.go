package storage

import (
	"fmt"
	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
	"net/http"
)

var (
	ErrScenarioNotFound = fmt.Errorf("Scenario not found")
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
	Close()
}

func NewRecorder(folder string, createNewFolder bool, newOnly bool, log utils.Logger) (Recorder, error) {
	fs := afero.NewOsFs()
	return NewRecorderWithFs(fs, folder, createNewFolder, newOnly, log)
}

/*
 *
 */
type ScenarioRecorder interface {
	Name() string
	ActivateScenario(name string) error
	RecordRequest(req *http.Request, session int64) (int64, error)
	RecordResponse(resp *http.Response, session int64) (int64, error)
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
