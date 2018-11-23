package storage

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"time"
)

type reqLogger struct {
	name string

	base *afero.Afero
	root *afero.Afero

	indexFile afero.File
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

	log.indexFile = file
	return nil
}

func (log *reqLogger) Name() string {
	return log.name
}

func (log *reqLogger) SaveReqMeta(meta ReqMeta) error {
	line := fmt.Sprintf("N\t%s\t%d\n", meta.Req.URL.String(), meta.Resp.StatusCode)
	_, err := log.indexFile.WriteString(line)
	if err != nil {
		return err
	}

	return nil
}
