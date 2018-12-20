package storage

import (
	"fmt"
	"net/http"

	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
)

type scSeekerImpl struct {
	root    *afero.Afero
	seekers map[string]Seeker
}

func NewScenarioSeekerWithFs(fs afero.Fs, folder string, log utils.Logger) (ScenarioSeeker, error) {
	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, folder)}

	content, err := root.ReadDir("")
	if err != nil {
		return nil, err
	}

	seekers := make(map[string]Seeker)
	for _, f := range content {
		if f.IsDir() {
			name := f.Name()
			seeker, err := NewSeekerWithFs(root, name)
			if err != nil {
				log.Error("Couldn't create Seeker on %s: %v", name, err)
			} else {
				seekers[name] = seeker
			}
		}
	}
	return &scSeekerImpl{
		root:    root,
		seekers: seekers,
	}, nil
}

func (p *scSeekerImpl) IsScenarioExists(name string) bool {
	_, ok := p.seekers[name]
	return ok
}

func (p *scSeekerImpl) Look(scenario string, method string, url string) (*http.Response, error) {
	seeker, ok := p.seekers[scenario]
	if ok {
		return seeker.Look(method, url)
	}
	return nil, fmt.Errorf("Unknow scenario: %s", scenario)
}
