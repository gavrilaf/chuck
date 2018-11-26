package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/afero"
	"io/ioutil"
	"net/http"
)

func createRequest() *http.Request {
	str := "{}"
	req, _ := http.NewRequest("POST", "https://secure.api.com?query=123", ioutil.NopCloser(bytes.NewBufferString(str)))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func createResponse() *http.Response {
	str := `{
		"colors": [
		  {
			"color": "black",
			"category": "hue",
			"type": "primary",
			"code": {
			  "rgba": [255,255,255,1],
			  "hex": "#000"
			}
		  },
		  {
			"color": "white",
			"category": "value",
			"code": {
			  "rgba": [0,0,0,1],
			  "hex": "#FFF"
			}
		  },]}`

	resp := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        make(http.Header),
		Body:          ioutil.NopCloser(bytes.NewBufferString(str)),
		ContentLength: int64(len(str)),
	}

	resp.Header.Set("Content-Type", "application/json")
	resp.Header.Set("Content-Length", "6573")

	return resp
}

var _ = Describe("Logger", func() {
	var (
		subject ReqLogger
		root    *afero.Afero
	)

	BeforeEach(func() {
		fs := afero.NewMemMapFs()
		root = &afero.Afero{Fs: fs}
	})

	Describe("Create logger", func() {
		var (
			err error
		)

		BeforeEach(func() {
			subject, err = NewLoggerWithFs(root.Fs)
		})

		It("should return nil error", func() {
			Expect(err).To(BeNil())
		})

		Context("When logger created", func() {
			var (
				dirExists   bool
				indexExists bool
			)

			BeforeEach(func() {
				path := "log/" + subject.Name()
				dirExists, _ = root.DirExists(path)
				indexExists, _ = root.Exists(path + "/index.txt")
			})

			It("should create a logger folder", func() {
				Expect(dirExists).To(BeTrue())
			})

			It("should create an index file", func() {
				Expect(indexExists).To(BeTrue())
			})
		})
	})

	Describe("Log request", func() {
		var (
			path  string
			err   error
			reqId string

			req  *http.Request
			resp *http.Response
		)

		BeforeEach(func() {
			subject, _ = NewLoggerWithFs(root.Fs)
			path = "log/" + subject.Name()

			req = createRequest()
			resp = createResponse()

			reqId, err = subject.LogRequest(req, resp)
		})

		It("should return nil error", func() {
			Expect(err).To(BeNil())
		})

		It("should index.txt contains log record", func() {
			fi, _ := root.Open(path + "/" + "index.txt")
			scanner := bufio.NewScanner(fi)
			scanner.Scan()
			line := scanner.Text()

			expected := fmt.Sprintf("N\thttps://secure.api.com?query=123\t200\t%s", reqId)
			Expect(line).To(Equal(expected))
		})

		It("should create request dump folder", func() {
			dirExists, _ := root.DirExists(path + "/" + reqId)
			Expect(dirExists).To(BeTrue())
		})

		It("should create request headers dump", func() {
			dumpExists, _ := root.Exists(path + "/" + reqId + "/req_header.json")
			Expect(dumpExists).To(BeTrue())
		})

		XIt("should create request body dump", func() {
			dumpExists, _ := root.Exists(path + "/" + reqId + "/req_body.json")
			Expect(dumpExists).To(BeTrue())
		})

		It("should create response headers dump", func() {
			dumpExists, _ := root.Exists(path + "/" + reqId + "/resp_header.json")
			Expect(dumpExists).To(BeTrue())
		})

		XIt("should create response body dump", func() {
			dumpExists, _ := root.Exists(path + "/" + reqId + "/resp_body.json")
			Expect(dumpExists).To(BeTrue())
		})
	})
})
