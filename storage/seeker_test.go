package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/gavrilaf/chuck/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"io/ioutil"
	"net/http"
)

var _ = Describe("Seeker", func() {
	var (
		log      Logger
		recorder Recorder
		root     *afero.Afero
		path     string
		subject  Seeker

		header   http.Header
		respBody string

		createRequest  func(method string, url string) *http.Request
		createResponse func() *http.Response
	)

	BeforeEach(func() {
		log = NewLogger(cli.NewMockUi())

		header = make(http.Header)
		header.Set("Content-Type", "application/json")
		header.Set("Access-Token", "Bearer-12234")

		createRequest = func(method string, url string) *http.Request {
			str := "{}"
			req, _ := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewBufferString(str)))
			req.Header.Set("Content-Type", "application/json")
			return req
		}

		createResponse = func() *http.Response {
			respBody = `{"colors": []}`

			resp := &http.Response{
				Status:        "200 OK",
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        header,
				Body:          ioutil.NopCloser(bytes.NewBufferString(respBody)),
				ContentLength: int64(len(respBody)),
			}
			return resp
		}

		fs := afero.NewMemMapFs()
		root = &afero.Afero{Fs: fs}

		recorder, _ = NewRecorderWithFs(fs, "test", false, log)

		recorder.RecordRequest(createRequest("POST", "https://secure.api.com/login"), 1)
		recorder.RecordResponse(createResponse(), 1)

		recorder.SetFocusedMode(true)

		recorder.RecordRequest(createRequest("GET", "https://secure.api.com/users/678/off"), 2)
		recorder.RecordResponse(createResponse(), 2)

		path = "test"
	})

	Describe("Open Seeker", func() {
		var (
			err error
		)
		BeforeEach(func() {
			subject, err = NewSeekerWithFs(root, path)
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
					resp, _ = subject.Look("GET", "https://secure.api.com/users")
				})

				It("should return request", func() {
					Expect(resp).ToNot(BeNil())
				})

				It("should response has correct headers", func() {
					Expect(resp.Header).To(Equal(header))
				})

				It("should response has correct body", func() {
					var buf []byte
					if resp.Body == nil {
						buf = make([]byte, 0)
					} else {
						buf, _ = ioutil.ReadAll(resp.Body)
					}

					Expect(string(buf)).To(Equal(respBody))
				})
			})

			Context("when request logged as unfocused", func() {
				BeforeEach(func() {
					resp, _ = subject.Look("POST", "https://secure.api.com/login")
				})

				It("should return nil", func() {
					Expect(resp).To(BeNil())
				})
			})
		})
	})
})
