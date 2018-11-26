package storage

import (
	"fmt"
	"github.com/spf13/afero"
	"net/http"
	"strconv"
	"time"

	"github.com/gavrilaf/chuck/utils"
)

type reqLogger struct {
	name      string
	root      *afero.Afero
	indexFile afero.File
	counter   int
}

func NewLogger() (ReqLogger, error) {
	fs := afero.NewOsFs()
	return NewLoggerWithFs(fs)
}

func NewLoggerWithFs(fs afero.Fs) (ReqLogger, error) {
	logDirExists, err := afero.DirExists(fs, "log")
	if err != nil {
		return nil, err
	}

	if !logDirExists {
		err := fs.Mkdir("log", 0777)
		if err != nil {
			return nil, err
		}
	}
	fmt.Printf("Created /log folder\n")

	tm := time.Now()
	name := fmt.Sprintf("%d_%d_%d_%d_%d_%d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	path := "log/" + name

	err = fs.Mkdir(path, 0777)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created %s folder\n", path)

	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, path)}
	file, err := root.Create("index.txt")
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created %s file\n", file.Name())

	return &reqLogger{
		name:      name,
		root:      root,
		indexFile: file,
		counter:   1}, nil
}

func (log *reqLogger) Name() string {
	return log.name
}

func (log *reqLogger) LogRequest(req *http.Request, resp *http.Response) (string, error) {
	recordID := "rq_" + strconv.Itoa(log.counter)
	line := fmt.Sprintf("N\t%s\t%d\t%s\n", req.URL.String(), resp.StatusCode, recordID)
	_, err := log.indexFile.WriteString(line)
	if err != nil {
		return "", err
	}

	err = log.root.Mkdir(recordID, 0777)
	if err != nil {
		return "", err
	}

	// TODO: error handling

	log.writeHeader(recordID+"/req_header.json", req.Header)

	log.writeHeader(recordID+"/resp_header.json", resp.Header)

	log.counter += 1

	return recordID, nil
}

////////////////////////////////////////////////////////////////////////////////////
// private

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
