package storage

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	IndexFileName = "index.txt"
)

var (
	ErrScenarioRecorderNotActivated = fmt.Errorf("Scenario recorder not activated (check intergation tests helper code)")
	ErrRequestNotFound              = fmt.Errorf("Request not found")
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

/*
 *
 */
type ScenarioRecorder interface {
	Tracker
	io.Closer
	Name() string
	ActivateScenario(name string) error
}

/*
 *
 */
type Seeker interface {
	Look(method string, url string) (*http.Response, error)
	Count() int
}

/*
 *
 */
type ScenarioSeeker interface {
	ScenariosCount() int
	IsScenarioExists(name string) bool
	Look(scenario string, method string, url string) (*http.Response, error)
}
