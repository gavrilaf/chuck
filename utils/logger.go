package utils

import (
	"fmt"
	"time"
)

type Logger interface {
	Request(id int64, method string, url string, statusCode int, elapsed time.Duration)

	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Panic(format string, args ...interface{})
}

func NewLogger() Logger {
	return &loggerImpl{}
}

/*
 * Implemenation
 */

type loggerImpl struct {
}

func (log *loggerImpl) Request(id int64, method string, url string, statusCode int, elapsed time.Duration) {
	fmt.Printf("--> [%d] : [%v] %s %s, %d \n", id, elapsed, method, url, statusCode)
}

func (log *loggerImpl) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (log *loggerImpl) Error(format string, args ...interface{}) {
	fmt.Printf("[ERR] "+format+"\n", args...)
}

func (log *loggerImpl) Panic(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
