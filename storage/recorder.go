package storage

import (
	"github.com/spf13/afero"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"

	"chuck/utils"
)

type recorderImpl struct {
	name      string
	focused   bool
	root      *afero.Afero
	indexFile afero.File
	index     Index
	tracker   Tracker
	mux       *sync.Mutex
	log       utils.Logger
}

func NewRecorder(fs afero.Fs, log utils.Logger, folder string, createNewFolder bool, onlyNew bool) (Recorder, error) {
	name, path, err := utils.PrepareStorageFolder(fs, folder, createNewFolder)
	if err != nil {
		return nil, err
	}

	root := &afero.Afero{Fs: afero.NewBasePathFs(fs, path)}
	indexFp, err := root.OpenFile(IndexFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	counter := 1
	var index Index
	if onlyNew {
		index, err = LoadIndex2(root, IndexFileName, false)
		if err != nil {
			return nil, err
		}

		counter = index.Size() + 1
	}

	return &recorderImpl{
		name:      name,
		root:      root,
		indexFile: indexFp,
		index:     index,
		tracker:   NewTracker(int64(counter), log),
		mux:       &sync.Mutex{},
		log:       log,
	}, nil
}

func (self *recorderImpl) Name() string {
	return self.name
}

func (self *recorderImpl) SetFocusedMode(focused bool) {
	self.focused = focused
}

func (self *recorderImpl) Close() error {
	return self.indexFile.Close()
}

func (self *recorderImpl) PendingCount() int {
	self.mux.Lock()
	defer self.mux.Unlock()

	return self.tracker.PendingCount()
}

func (self *recorderImpl) RecordRequest(req *http.Request, session int64) (*PendingRequest, error) {
	pendingReq, err := self.checkAndRecordRequest(req, session)
	if pendingReq == nil {
		return nil, err
	}

	folder := "r_" + strconv.FormatInt(pendingReq.Id, 10)

	err = self.root.Mkdir(folder, 0777)
	if err != nil {
		return nil, err
	}

	err = self.writeHeader(path.Join(folder, "req_header.json"), req.Header)
	if err != nil {
		self.log.Error("Couldn't write request header: %v", err)
	}

	self.writeRequesteBody(path.Join(folder, "req_body.json"), req)
	if err != nil {
		self.log.Error("Couldn't write request body: %v", err)
	}

	return pendingReq, nil
}

func (self *recorderImpl) RecordResponse(resp *http.Response, session int64) (*PendingRequest, error) {
	self.mux.Lock()
	pendingReq, err := self.tracker.RecordResponse(resp, session)
	self.mux.Unlock()

	if err != nil {
		if self.index != nil {
			return nil, nil
		} else {
			self.log.Panic("Could not find request for session: %d, %v", session, err)
		}
	}

	folder := "r_" + strconv.FormatInt(pendingReq.Id, 10)
	line := FormatIndexItem(pendingReq.Method, pendingReq.Url, resp.StatusCode, folder, self.focused)
	_, err = self.indexFile.WriteString(line + "\n")
	if err != nil {
		return nil, err
	}

	err = self.writeHeader(path.Join(folder, "resp_header.json"), resp.Header)
	if err != nil {
		self.log.Error("Couldn't write response header: %v", err)
	}

	err = self.writeResponseBody(path.Join(folder, "resp_body.json"), resp)
	if err != nil {
		self.log.Error("Couldn't write response body: %v", err)
	}

	return pendingReq, nil
}

////////////////////////////////////////////////////////////////////////////////////
// Private

func (self *recorderImpl) checkAndRecordRequest(req *http.Request, session int64) (*PendingRequest, error) {
	method := req.Method
	url := req.URL.String()

	self.mux.Lock()
	defer self.mux.Unlock()

	if self.index != nil {
		if self.index.Find(method, url, SEARCH_SUBSTR) != nil {
			return nil, nil
		} else {
			self.index.Add(IndexItem{
				Focused: false,
				Method:  method,
				Url:     url,
				Code:    0,
				Folder:  "",
			})
		}
	}

	return self.tracker.RecordRequest(req, session)
}

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
		return err
	}

	if len(b) == 0 {
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
		return err
	}

	if len(b) == 0 {
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
