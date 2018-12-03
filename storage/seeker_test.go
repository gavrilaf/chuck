package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	//"bufio"
	"bytes"
	//"fmt"
	"github.com/spf13/afero"
	"io/ioutil"
	"net/http"
)

var _ = Describe("Seeker", func() {
	var (
		logger  ReqLogger
		root    *afero.Afero
		path    string
		subject ReqSeeker

		createRequest  func(method string, url string) *http.Request
		createResponse func() *http.Response
	)

	BeforeEach(func() {
		createRequest = func(method string, url string) *http.Request {
			str := "{}"
			req, _ := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewBufferString(str)))
			req.Header.Set("Content-Type", "application/json")
			return req
		}

		createResponse = func() *http.Response {
			str := `{"colors": []}`

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
			resp.Header.Set("Content-Length", "15")

			return resp
		}

		fs := afero.NewMemMapFs()
		root = &afero.Afero{Fs: fs}

		logger, _ = NewLoggerWithFs("test", fs)

		logger.LogRequest(createRequest("POST", "https://secure.api.com/login"), 1)
		logger.LogResponse(createResponse(), 1)

		logger.SetFocusedMode(true)

		logger.LogRequest(createRequest("GET", "https://secure.api.com/users"), 2)
		logger.LogResponse(createResponse(), 2)

		path = "test/" + logger.Name()
	})

	Describe("Open Seeker", func() {
		var (
			err error
		)
		BeforeEach(func() {
			subject, err = NewSeekerWithFs(path, root)
		})

		It("should return nil error", func() {
			Expect(err).To(BeNil())
		})

		It("should return Seeker", func() {
			Expect(subject).ToNot(BeNil())
		})

		Describe("looking for request", func() {
			var (
				resp *http.Response
			)

			Context("when request logged as focused", func() {
				BeforeEach(func() {
					resp = subject.Look("GET", "https://secure.api.com/users")
				})

				XIt("should return request", func() {
					Expect(resp).ToNot(BeNil())
				})
			})

			Context("when request logged as unfocused", func() {
				BeforeEach(func() {
					resp = subject.Look("POST", "https://secure.api.com/login")
				})

				It("should return nil", func() {
					Expect(resp).To(BeNil())
				})
			})
		})
	})
})
