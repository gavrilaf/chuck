package handlers_test

import (
	. "chuck/handlers"
	. "chuck/storage"
	. "chuck/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"gopkg.in/elazarl/goproxy.v1"
	"net/http"
)

var _ = Describe("Seeker handler", func() {
	var (
		log Logger
		fs  afero.Fs

		header  http.Header
		req     *http.Request
		resp    *http.Response
		err     error
		context *goproxy.ProxyCtx

		subject ProxyHandler
	)

	BeforeEach(func() {
		log = NewLogger(cli.NewMockUi())
		fs = afero.NewMemMapFs()

		context = &goproxy.ProxyCtx{Session: 100}
	})

	Describe("open seeker handler on the folder with index", func() {
		BeforeEach(func() {
			header = make(http.Header)
			header.Set("Content-Type", "application/json")

			req, _ = MakeRequest2("POST", "https://secure.api.com/login", header, "")
			resp = MakeResponse2(200, header, "{}")

			recorder, _ := NewRecorder(fs, log, "test", false, false, false, true)
			recorder.SetFocusedMode(true)

			recorder.RecordRequest(req, 1)
			recorder.RecordResponse(resp, 1)

			cfg := &SeekerConfig{
				BaseConfig: BaseConfig{
					Folder: "test",
				},
			}

			subject, err = NewSeekerHandler(cfg, fs, log)
		})

		It("should no error occured", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return seeker proxy handler", func() {
			Expect(subject).ToNot(BeNil())
		})

		Context("when handle focused request", func() {
			BeforeEach(func() {
				resp = subject.Request(req, nil) // pass context as nil because we don't use it if we found response for request
			})

			It("should return valid response", func() {
				Expect(resp).ToNot(BeNil())
				Expect(resp.StatusCode).To(Equal(200))
				Expect(resp.Header).To(Equal(header))
			})
		})

		Context("when handle new request", func() {
			BeforeEach(func() {
				reqNew, _ := MakeRequest2("GET", "www.unknown-host.net", header, "")
				resp = subject.Request(reqNew, context)
			})

			It("should return nil", func() {
				Expect(resp).To(BeNil())
			})

			// TODO: How to test request/response tracking? Add mock for the tracker? Listen stdout?
		})
	})
})
