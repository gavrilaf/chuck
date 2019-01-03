package handlers

import (
	"github.com/gavrilaf/chuck/storage"
	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
	"sync"
)

// Recorder

func NewRecorderHandler(config *RecorderConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	recorder, err := storage.NewRecorder(fs, log, config.Folder, config.CreateNewFolder, false)
	if err != nil {
		return nil, err
	}

	return &recordHandler{
		recorder:       recorder,
		log:            log,
		preventCaching: config.Prevent304,
	}, nil
}

// Seeker

func NewSeekerHandler(config *SeekerConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	seeker, err := storage.NewSeeker(fs, config.Folder)
	if err != nil {
		return nil, err
	}

	return &seekerHandler{
		seeker:  seeker,
		tracker: storage.NewTracker(0, log),
		mux:     &sync.Mutex{},
		log:     log,
	}, nil
}

// Scenario Seeker

func NewScenarioSeekerHandler(config *ScenarioSeekerConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	seeker, err := storage.NewScenarioSeeker(fs, log, config.Folder)
	if err != nil {
		return nil, err
	}

	return &scenarioSeekerHandler{
		seeker:    seeker,
		log:       log,
		scenarios: make(map[string]string),
	}, nil
}

// Scenario Recorder

func NewScenarioRecorderHandler(config *ScenarioRecorderConfig, fs afero.Fs, log utils.Logger) (ProxyHandler, error) {
	recorder, err := storage.NewScenarioRecorder(fs, log, config.Folder, true)
	if err != nil {
		return nil, err
	}

	return &scenarioRecordHandler{
		recorder:       recorder,
		log:            log,
		scenarios:      make(map[string]string),
		preventCaching: true,
	}, nil
}