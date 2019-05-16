package storage

import (
	"fmt"
	"github.com/spf13/afero"
	"net/http"
	"os"
	. "path"

	"chuck/utils"
)

type scSeekerImpl struct {
	root    *afero.Afero
	seekers map[string]Seeker
}

func NewScenarioSeeker(fs afero.Fs, log utils.Logger, folder string) (ScenarioSeeker, error) {
	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, folder)}
	seekers := make(map[string]Seeker)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error("Couldn't access a path: %s, %v", path, err)
			return err
		}

		if !info.IsDir() && info.Name() == IndexFileName {
			folder, _ := Split(path)
			scenarioName := Base(folder)
			
			if _, ok := seekers[scenarioName], ok {
				return fmt.Errorf("Scenario %s (%s) already opened", scenarioName, folder)
			}


			seeker, err := NewSeeker(root, folder)
			if err != nil {
				log.Error("Couldn't load index by path %s, %v", path, err)
				return err
			} else {
				log.Info("Loaded scenario %s", folder)
			}
			
			seekers[scenarioName] = seeker			
		}

		return nil
	}

	err := root.Walk("", walkFn)
	if err != nil {
		return nil, err
	}

	return &scSeekerImpl{
		root:    root,
		seekers: seekers,
	}, nil
}

func (self *scSeekerImpl) ScenariosCount() int {
	return len(self.seekers)
}

func (self *scSeekerImpl) IsScenarioExists(name string) bool {
	_, ok := self.seekers[name]
	return ok
}

func (self *scSeekerImpl) Look(scenario string, method string, url string) (*http.Response, error) {
	seeker, ok := self.seekers[scenario]
	if ok {
		return seeker.Look(method, url)
	}
	return nil, fmt.Errorf("Unknow scenario: %s", scenario)
}
