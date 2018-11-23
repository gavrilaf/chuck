package storage

import (
	"github.com/spf13/afero"
	"os"
	"time"
)

type reqLogger struct {
	name string

	base *afero.Afero
	root *afero.Afero

	index afero.File
}

func (log *reqLogger) Start() error {
	log.name = time.Now().Format("2006_01_02_04_05_01")

	err := log.base.Mkdir(log.name, os.ModeDir)
	if err != nil {
		return err
	}

	log.root = &afero.Afero{Fs: afero.NewBasePathFs(log.base, log.name)}
	file, err := log.root.Create("index.txt")
	if err != nil {
		return err
	}

	log.index = file
	return nil
}

func (log *reqLogger) Name() string {
	return log.name
}

func (log *reqLogger) SaveReqMeta(meta ReqMeta) {
	// Calc request hash
	// Add line to the index
	// Create 'hash' folder
	// Save req_headers.txt
	// Save req_body.*
	// Save resp_headers.txt
	// Save resp_body.*
}
