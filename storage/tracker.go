package storage

import (
	"net/http"
	"time"

	"chuck/utils"
)

type trackerImpl struct {
	counter int64
	pending map[int64]*PendingRequest
	log     utils.Logger
}

func NewTracker(counter int64, log utils.Logger) Tracker {
	return &trackerImpl{
		counter: counter,
		pending: make(map[int64]*PendingRequest, 10),
		log:     log,
	}
}

func (self *trackerImpl) PendingCount() int {
	return len(self.pending)
}

func (self *trackerImpl) RecordRequest(req *http.Request, session int64) (*PendingRequest, error) {
	p := &PendingRequest{
		Id:      self.counter,
		Method:  req.Method,
		Url:     req.URL.String(),
		Started: time.Now(),
	}

	self.pending[session] = p
	self.counter += 1

	return p, nil
}

func (self *trackerImpl) RecordResponse(resp *http.Response, session int64) (*PendingRequest, error) {
	req, ok := self.pending[session]
	if !ok {
		return nil, ErrRequestNotFound
	}

	delete(self.pending, session)

	elapsed := time.Since(req.Started)
	self.log.Request(req.Id, req.Method, req.Url, resp.StatusCode, elapsed)

	return req, nil
}
