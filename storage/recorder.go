package storage

import (
	"fmt"
	"net/http"
	"os"
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
	index     Index
	counter   int64
	pending   map[int64]pendingRequest
	mux       *sync.Mutex
	log       utils.Logger
}

func NewRecorderWithFs(fs afero.Fs, folder string, createNewFolder bool, onlyNew bool, log utils.Logger) (Recorder, error) {
	folder = strings.Trim(folder, " \\/")
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

	name := ""
	path := folder
	if createNewFolder {
		tm := time.Now()
		name = fmt.Sprintf("%d_%d_%d_%d_%d_%d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
		path = folder + "/" + name

		err = fs.Mkdir(path, 0777)
		if err != nil {
			return nil, err
		}
	}

	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, path)}
	indexFp, err := root.OpenFile("index.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	counter := 1
	var index Index
	if onlyNew {
		//content, err := root.ReadFile("index.txt")
		//fmt.Printf("\n***Err: %v, Index: \n%s\n\n", err, string(content))
		index, err = LoadIndex2(root, "index.txt", false)
		if err != nil {
			return nil, err
		}

		counter = index.Size() + 1
		//fmt.Printf("\n*** Loading index with %d records\n\n", index.Size())
	}

	return &recorderImpl{
		name:      name,
		root:      root,
		indexFile: indexFp,
		index:     index,
		counter:   int64(counter),
		pending:   make(map[int64]pendingRequest, 10),
		mux:       &sync.Mutex{},
		log:       log,
	}, nil
}

func (recorder *recorderImpl) Name() string {
	return recorder.name
}

func (recorder *recorderImpl) SetFocusedMode(focused bool) {
	recorder.focused = focused
}

func (recorder *recorderImpl) Close() {
	recorder.indexFile.Close()
}

func (recorder *recorderImpl) PendingCount() int {
	recorder.mux.Lock()
	defer recorder.mux.Unlock()

	return len(recorder.pending)
}

func (recorder *recorderImpl) RecordRequest(req *http.Request, session int64) (int64, error) {
	recorder.mux.Lock()
	defer recorder.mux.Unlock()

	method := req.Method
	url := req.URL.String()

	if recorder.index != nil {
		if recorder.index.Find(method, url, SEARCH_SUBSTR) != nil {
			return -1, nil
		} else {
			recorder.index.Add(IndexItem{
				Focused: false,
				Method:  method,
				Url:     url,
				Code:    0,
				Folder:  "",
			})
		}
	}

	id := recorder.counter
	folder := "r_" + strconv.FormatInt(id, 10)

	err := recorder.root.Mkdir(folder, 0777)
	if err != nil {
		return 0, err
	}

	recorder.pending[session] = pendingRequest{
		id:      id,
		method:  method,
		url:     url,
		started: time.Now(),
	}

	recorder.counter += 1

	err = recorder.writeHeader(folder+"/req_header.json", req.Header)
	if err != nil {
		recorder.log.Error("Couldn't write request header: %v", err)
	}

	recorder.writeRequesteBody(folder+"/req_body.json", req)
	if err != nil {
		recorder.log.Error("Couldn't write request body: %v", err)
	}

	return id, nil
}

func (recorder *recorderImpl) RecordResponse(resp *http.Response, session int64) (int64, error) {
	recorder.mux.Lock()
	defer recorder.mux.Unlock()

	req, ok := recorder.pending[session]
	if !ok {
		if recorder.index != nil {
			return -1, nil
		} else {
			recorder.log.Panic("Could not find request for session: %d", session)
		}
	}

	delete(recorder.pending, session)

	elapsed := time.Since(req.started)
	recorder.log.Request(req.id, req.method, req.url, resp.StatusCode, elapsed)

	folder := "r_" + strconv.FormatInt(req.id, 10)
	line := FormatIndexItem(req.method, req.url, resp.StatusCode, folder, recorder.focused)
	_, err := recorder.indexFile.WriteString(line + "\n")
	if err != nil {
		return 0, err
	}

	err = recorder.writeHeader(folder+"/resp_header.json", resp.Header)
	if err != nil {
		recorder.log.Error("Couldn't write response header: %v", err)
	}

	err = recorder.writeResponseBody(folder+"/resp_body.json", resp)
	if err != nil {
		recorder.log.Error("Couldn't write response body: %v", err)
	}

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
