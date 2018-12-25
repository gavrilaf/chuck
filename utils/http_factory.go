package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

func MakeRequest(method string, url string, header http.Header, body io.ReadCloser) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = header
	return req, nil
}

func MakeRequest2(method string, url string, header http.Header, body string) (*http.Request, error) {
	return MakeRequest(method, url, header, ioutil.NopCloser(bytes.NewBufferString(body)))
}

func MakeResponse(code int, header http.Header, body io.ReadCloser, contentLen int64) *http.Response {
	return &http.Response{
		StatusCode:    code,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		Body:          body,
		ContentLength: contentLen,
	}
}

func MakeResponse2(code int, header http.Header, body string) *http.Response {
	return MakeResponse(code, header, ioutil.NopCloser(bytes.NewBufferString(body)), int64(len(body)))
}
