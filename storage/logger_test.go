package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bufio"
	"bytes"
	"github.com/spf13/afero"
	"io/ioutil"
	"net/http"
)

var _ = Describe("Logger", func() {
	var (
		subject ReqLogger
		afr     *afero.Afero
	)

	BeforeEach(func() {
		afr = &afero.Afero{Fs: afero.NewMemMapFs()}
		subject = NewLoggerWithFs(afr)
	})

	Describe("Start logger", func() {
		var (
			err         error
			dirExists   bool
			indexExists bool
		)

		BeforeEach(func() {
			err = subject.Start()

			path := subject.Name()
			dirExists, _ = afr.DirExists(path)
			indexExists, _ = afr.Exists(path + "/" + "index.txt")
		})

		It("should return nil error", func() {
			Expect(err).To(BeNil())
		})

		It("should create a logger folder", func() {
			Expect(dirExists).To(BeTrue())
		})

		It("should create an index file", func() {
			Expect(indexExists).To(BeTrue())
		})
	})

	Describe("Start logger", func() {
		var (
			path string
			err  error
		)
		BeforeEach(func() {
			_ = subject.Start()
			path = subject.Name()

			reqBody := "{request: 1}"
			req, _ := http.NewRequest("POST", "https://secure.api.com?query=123", ioutil.NopCloser(bytes.NewBufferString(reqBody)))

			respBody := "{error: 1}"
			resp := &http.Response{
				Status:        "200 OK",
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Body:          ioutil.NopCloser(bytes.NewBufferString(respBody)),
				ContentLength: int64(len(respBody)),
			}

			err = subject.SaveReqMeta(ReqMeta{Req: req, Resp: resp})
		})

		It("should return nil error", func() {
			Expect(err).To(BeNil())
		})

		It("should index.txt contains log record", func() {
			fi, _ := afr.Open(path + "/" + "index.txt")
			scanner := bufio.NewScanner(fi)
			scanner.Scan()
			line := scanner.Text()

			expected := "N\thttps://secure.api.com?query=123\t200"
			Expect(line).To(Equal(expected))
		})

	})
})
