package storage

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
)

type pendingRequest struct {
	id      int64
	method  string
	url     string
	started time.Time
}

type reqLogger struct {
	name      string
	root      *afero.Afero
	indexFile afero.File
	counter   int64
	pending   map[int64]pendingRequest
	mux       *sync.Mutex
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

	tm := time.Now()
	name := fmt.Sprintf("%d_%d_%d_%d_%d_%d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	path := "log/" + name

	err = fs.Mkdir(path, 0777)
	if err != nil {
		return nil, err
	}

	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, path)}
	file, err := root.Create("index.txt")
	if err != nil {
		return nil, err
	}

	return &reqLogger{
		name:      name,
		root:      root,
		indexFile: file,
		counter:   1,
		pending:   make(map[int64]pendingRequest, 10),
		mux:       &sync.Mutex{}}, nil
}

func (log *reqLogger) Name() string {
	return log.name
}

func (log *reqLogger) PendingCount() int {
	log.mux.Lock()
	defer log.mux.Unlock()

	return len(log.pending)
}

func (log *reqLogger) LogRequest(req *http.Request, session int64) (int64, error) {
	log.mux.Lock()
	defer log.mux.Unlock()

	id := log.counter
	folder := "r_" + strconv.FormatInt(id, 10)

	err := log.root.Mkdir(folder, 0777)
	if err != nil {
		return 0, err
	}

	log.pending[session] = pendingRequest{
		id:      id,
		method:  req.Method,
		url:     req.URL.String(),
		started: time.Now(),
	}

	log.counter += 1

	log.writeHeader(folder+"/req_header.json", req.Header)

	return id, nil
}

func (log *reqLogger) LogResponse(resp *http.Response, session int64) (int64, error) {
	log.mux.Lock()
	defer log.mux.Unlock()

	req, ok := log.pending[session]
	if !ok {
		panic(fmt.Errorf("Could not find request for session: %d\n", session))
	}

	delete(log.pending, session)

	elapsed := time.Since(req.started)
	fmt.Printf("--> [%d] : [%v] %s %s, %v \n", req.id, elapsed, req.method, req.url, resp.Status)

	line := fmt.Sprintf("N\tr_%d\t%s\t%s\t%d\n", req.id, req.method, req.url, resp.StatusCode)
	_, err := log.indexFile.WriteString(line)
	if err != nil {
		return 0, err
	}

	folder := "r_" + strconv.FormatInt(req.id, 10)
	log.writeHeader(folder+"/resp_header.json", resp.Header)
	log.writeResponseBody(folder+"/resp_body.json", resp)

	return req.id, nil
}

////////////////////////////////////////////////////////////////////////////////////
// Private

func (log *reqLogger) writeHeader(fname string, header http.Header) error {
	if len(header) > 0 {
		fp, err := log.root.Create(fname)
		if err != nil {
			return err
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

func (log *reqLogger) writeResponseBody(fname string, resp *http.Response) error {
	b, err := utils.DumpRespBody(resp)
	if err != nil {
		return nil
	}

	fp, err := log.root.Create(fname)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(b)
	return err
}
