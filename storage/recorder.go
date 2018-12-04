package storage

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

type recorderImpl struct {
	name      string
	focused   bool
	root      *afero.Afero
	indexFile afero.File
	counter   int64
	pending   map[int64]pendingRequest
	mux       *sync.Mutex
}

func NewRecorderWithFs(folder string, fs afero.Fs) (Recorder, error) {
	folder = strings.Trim(folder, " \\/")
	if len(folder) == 0 {
		folder = "log"
	}

	logDirExists, err := afero.DirExists(fs, folder)
	if err != nil {
		return nil, err
	}

	if !logDirExists {
		err := fs.Mkdir(folder, 0777)
		if err != nil {
			return nil, err
		}
	}

	tm := time.Now()
	name := fmt.Sprintf("%d_%d_%d_%d_%d_%d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
	path := folder + "/" + name

	err = fs.Mkdir(path, 0777)
	if err != nil {
		return nil, err
	}

	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, path)}
	file, err := root.Create("index.txt")
	if err != nil {
		return nil, err
	}

	return &recorderImpl{
		name:      name,
		root:      root,
		indexFile: file,
		counter:   1,
		pending:   make(map[int64]pendingRequest, 10),
		mux:       &sync.Mutex{}}, nil
}

func (recorder *recorderImpl) Name() string {
	return recorder.name
}

func (recorder *recorderImpl) SetFocusedMode(focused bool) {
	recorder.focused = focused
}

func (recorder *recorderImpl) PendingCount() int {
	recorder.mux.Lock()
	defer recorder.mux.Unlock()

	return len(recorder.pending)
}

func (recorder *recorderImpl) RecordRequest(req *http.Request, session int64) (int64, error) {
	recorder.mux.Lock()
	defer recorder.mux.Unlock()

	id := recorder.counter
	folder := "r_" + strconv.FormatInt(id, 10)

	err := recorder.root.Mkdir(folder, 0777)
	if err != nil {
		return 0, err
	}

	recorder.pending[session] = pendingRequest{
		id:      id,
		method:  req.Method,
		url:     req.URL.String(),
		started: time.Now(),
	}

	recorder.counter += 1

	recorder.writeHeader(folder+"/req_header.json", req.Header)
	recorder.writeRequesteBody(folder+"/req_body.json", req)

	return id, nil
}

func (recorder *recorderImpl) RecordResponse(resp *http.Response, session int64) (int64, error) {
	recorder.mux.Lock()
	defer recorder.mux.Unlock()

	req, ok := recorder.pending[session]
	if !ok {
		panic(fmt.Errorf("Could not find request for session: %d\n", session))
	}

	delete(recorder.pending, session)

	elapsed := time.Since(req.started)
	fmt.Printf("--> [%d] : [%v] %s %s, %v \n", req.id, elapsed, req.method, req.url, resp.Status)

	mode := "N"
	if recorder.focused {
		mode = "F"
	}

	line := fmt.Sprintf("%s\tr_%d\t%s\t%s\t%d\n", mode, req.id, req.method, req.url, resp.StatusCode)
	_, err := recorder.indexFile.WriteString(line)
	if err != nil {
		return 0, err
	}

	folder := "r_" + strconv.FormatInt(req.id, 10)
	recorder.writeHeader(folder+"/resp_header.json", resp.Header)
	recorder.writeResponseBody(folder+"/resp_body.json", resp)

	return req.id, nil
}

////////////////////////////////////////////////////////////////////////////////////
// Private

func (recorder *recorderImpl) writeHeader(fname string, header http.Header) error {
	if len(header) > 0 {
		fp, err := recorder.root.Create(fname)
		if err != nil {
			return err
		}
		defer fp.Close()

		buf, err := utils.EncodeHeaders(header)
		if err != nil {
			return err
		}
		_, err = fp.Write(buf)
		return err
	}
	return nil
}

func (recorder *recorderImpl) writeResponseBody(fname string, resp *http.Response) error {
	b, err := utils.DumpRespBody(resp)
	if err != nil {
		return nil
	}

	fp, err := recorder.root.Create(fname)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(b)
	return err
}

func (recorder *recorderImpl) writeRequesteBody(fname string, req *http.Request) error {
	b, err := utils.DumpReqBody(req)
	if err != nil {
		return nil
	}

	fp, err := recorder.root.Create(fname)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = fp.Write(b)
	return err
}
