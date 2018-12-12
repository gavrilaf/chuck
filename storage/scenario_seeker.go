package storage

import (
	"net/http"

	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
)

type scSeekerImpl struct {
	root    *afero.Afero
	seekers map[string]Seeker
	log     utils.Logger
}

func NewScSeekerWithFs(folder string, fs afero.Fs, log utils.Logger) (ScSeeker, error) {
	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, folder)}

	content, err := root.ReadDir("")
	if err != nil {
		return nil, err
	}

	seekers := make(map[string]Seeker)
	for _, f := range content {
		if f.IsDir() {
			name := f.Name()
			log.Info(name)
			seeker, err := NewSeekerWithFs(name, root, log)
			if err != nil {
				log.Error("Couldn't create Seeker on %s: %v", name, err)
			} else {
				seekers[name] = seeker
			}
		}
	}

	log.Info("Scenario seeker created, loaded %d scenarious", len(seekers))

	return &scSeekerImpl{
		root:    root,
		seekers: seekers,
		log:     log,
	}, nil
}

func (sc *scSeekerImpl) IsScenarioExists(name string) bool {
	_, ok := sc.seekers[name]
	return ok
}

func (sc *scSeekerImpl) Look(scenario string, method string, url string) *http.Response {
	seeker, ok := sc.seekers[scenario]
	if ok {
		return seeker.Look(method, url)
	}
	return nil
}
