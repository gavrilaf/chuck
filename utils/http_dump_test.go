package utils_test

import (
	. "chuck/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io/ioutil"
	"net/http"
)

var _ = Describe("HttpDump", func() {
	var (
		err error
		buf []byte
	)

	Describe("Dump http response", func() {
		Context("Empty response", func() {
			BeforeEach(func() {
				resp := &http.Response{
					Status:        "200 OK",
					StatusCode:    200,
					Proto:         "HTTP/1.1",
					ProtoMajor:    1,
					ProtoMinor:    1,
					Body:          nil,
					ContentLength: int64(0),
				}

				buf, err = DumpRespBody(resp)
			})

			It("should return nil errors", func() {
				Expect(err).To(BeNil())
			})

			It("should return empty buffer", func() {
				Expect(buf).To(Equal([]byte("")))
			})
		})

		Context("Filled response", func() {
			var str string

			BeforeEach(func() {
				str = `{"allowed": 1}`
				resp := &http.Response{
					Status:        "200 OK",
					StatusCode:    200,
					Proto:         "HTTP/1.1",
					ProtoMajor:    1,
					ProtoMinor:    1,
					Body:          ioutil.NopCloser(bytes.NewBufferString(str)),
					ContentLength: int64(len(str)),
				}

				buf, err = DumpRespBody(resp)
			})

			It("should return nil errors", func() {
				Expect(err).To(BeNil())
			})

			It("should return filled buffer", func() {
				Expect(buf).To(Equal([]byte(str)))
			})
		})

		Describe("Dump http request", func() {
			Context("Empty resporequestnse", func() {
				BeforeEach(func() {
					req, _ := http.NewRequest("GET", "https://secure.api.com?query=123", nil)
					buf, err = DumpReqBody(req)
				})

				It("should return nil errors", func() {
					Expect(err).To(BeNil())
				})

				It("should return empty buffer", func() {
					Expect(buf).To(Equal([]byte("")))
				})
			})

			Context("Filled response", func() {
				var str string

				BeforeEach(func() {
					str = `{"allowed": 1}`
					req, _ := http.NewRequest("GET", "https://secure.api.com?query=123", ioutil.NopCloser(bytes.NewBufferString(str)))
					buf, err = DumpReqBody(req)
				})

				It("should return nil errors", func() {
					Expect(err).To(BeNil())
				})

				It("should return filled buffer", func() {
					Expect(buf).To(Equal([]byte(str)))
				})
			})
		})
	})
})
