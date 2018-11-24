package storage

import (
	"fmt"
	"github.com/spf13/afero"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gavrilaf/chuck/utils"
)

type reqLogger struct {
	name string

	base *afero.Afero
	root *afero.Afero

	indexFile afero.File

	counter int
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
	log.counter = 1

	return nil
}

func (log *reqLogger) Name() string {
	return log.name
}

func (log *reqLogger) SaveReqMeta(meta ReqMeta) (string, error) {
	recordID := "rq_" + strconv.Itoa(log.counter)
	line := fmt.Sprintf("N\t%s\t%d\t%s\n", meta.Req.URL.String(), meta.Resp.StatusCode, recordID)
	_, err := log.indexFile.WriteString(line)
	if err != nil {
		return "", err
	}

	err = log.root.Mkdir(recordID, os.ModeDir)
	if err != nil {
		return "", err
	}

	// TODO: error handling

	log.writeHeader(recordID+"/req_header.json", meta.Req.Header)

	log.writeHeader(recordID+"/resp_header.json", meta.Resp.Header)

	return recordID, nil
}

func (log *reqLogger) writeHeader(fname string, header http.Header) error {
	if len(header) > 0 {
		fp, err := log.root.Create(fname)
		if err != nil {
			return nil
		}

		defer fp.Close()

		buff, err := utils.EncodeHeaders(header)
		if err != nil {
			return err
		}

		_, err = fp.Write(buff)
		return err
	}
	return nil
}
