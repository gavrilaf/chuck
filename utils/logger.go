package utils

import (
	"fmt"
	"github.com/mitchellh/cli"
	"time"
)

type Logger interface {
	Request(id int64, method string, url string, statusCode int, elapsed time.Duration)
	FocusedReq(method string, url string, statusCode int)

	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Panic(format string, args ...interface{})
}

func NewLogger(ui cli.Ui) Logger {
	return &loggerImpl{
		ui: ui,
	}
}

/*
 * Implemenation
 */

type loggerImpl struct {
	ui cli.Ui
}

func (log *loggerImpl) Request(id int64, method string, url string, statusCode int, elapsed time.Duration) {
	s := fmt.Sprintf("--> [%d] : [%v] %s : %s, %d", id, elapsed, method, url, statusCode)
	log.printForStatusCode(s, statusCode)
}

func (log *loggerImpl) FocusedReq(method string, url string, statusCode int) {
	s := fmt.Sprintf("<-- %s : %s, %d", method, url, statusCode)
	log.printForStatusCode(s, statusCode)
}

func (log *loggerImpl) Info(format string, args ...interface{}) {
	log.ui.Info(fmt.Sprintf("[INFO] "+format, args...))
}

func (log *loggerImpl) Warn(format string, args ...interface{}) {
	log.ui.Warn(fmt.Sprintf("[WARN] "+format, args...))
}

func (log *loggerImpl) Error(format string, args ...interface{}) {
	log.ui.Error(fmt.Sprintf("[ERR] "+format, args...))
}

func (log *loggerImpl) Panic(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}

func (log *loggerImpl) printForStatusCode(s string, code int) {
	if code < 400 {
		log.ui.Info(s)
	} else if code >= 400 && code < 500 {
		log.ui.Warn(s)
	} else {
		log.ui.Error(s)
	}
}
