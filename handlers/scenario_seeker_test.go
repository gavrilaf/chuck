package handlers_test

import (
	. "github.com/gavrilaf/chuck/handlers"
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/gavrilaf/chuck/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("ScenarioSeeker handler", func() {
	var (
		log Logger
		fs  afero.Fs

		respRecorder *httptest.ResponseRecorder
		resp         *http.Response
		emptyHeader  http.Header

		err     error
		subject ProxyHandler
	)

	BeforeEach(func() {
		log = NewLogger(cli.NewMockUi())
		fs = afero.NewMemMapFs()
	})

	Describe("open scenario recorder", func() {

		BeforeEach(func() {
			emptyHeader = make(http.Header)

			resp = MakeResponse2(200, emptyHeader, "")

			req1, _ := MakeRequest("POST", "https://secure.api.com/login", emptyHeader, nil)
			req2, _ := MakeRequest("GET", "https://secure.api.com/users", emptyHeader, nil)

			recorder1, _ := NewRecorder(fs, log, "test/scenario-1", false, false)
			defer recorder1.Close()

			recorder1.SetFocusedMode(true)

			recorder1.RecordRequest(req1, 1)
			recorder1.RecordResponse(resp, 1)

			recorder2, _ := NewRecorder(fs, log, "test/scenario-2", false, false)
			defer recorder2.Close()

			recorder2.SetFocusedMode(true)

			recorder2.RecordRequest(req2, 1)
			recorder2.RecordResponse(resp, 1)

			cfg := &ScenarioSeekerConfig{
				BaseConfig: BaseConfig{
					Folder: "test",
				},
			}

			subject, err = NewScenarioSeekerHandler(cfg, fs, log)
		})

		It("should no error occured", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return seeker proxy handler", func() {
			Expect(subject).ToNot(BeNil())
		})

		Context("when activate unknown scenario", func() {
			BeforeEach(func() {
				req, _ := MakeRequest("PUT", "https://127.0.0.1/scenario/scenario-111/scenario-111-id/no", emptyHeader, nil)
				respRecorder = httptest.NewRecorder()
				subject.NonProxyHandler(respRecorder, req)
			})

			It("should return 404", func() {
				Expect(respRecorder.Code).To(Equal(404))
			})
		})

		Context("when activate existing scenario", func() {
			BeforeEach(func() {
				req, _ := MakeRequest("PUT", "https://127.0.0.1/scenario/scenario-1/scenario-1-id/no", emptyHeader, nil)
				respRecorder = httptest.NewRecorder()
				subject.NonProxyHandler(respRecorder, req)
			})

			It("should activate scenario", func() {
				Expect(respRecorder.Code).To(Equal(200))
			})

			Context("when handle request from the scenario", func() {
				BeforeEach(func() {
					header := make(http.Header)
					header.Set(ScenarioIdHeader, "scenario-1-id")
					req, _ := MakeRequest("POST", "https://secure.api.com/login", header, nil)
					resp = subject.Request(req, nil)
				})

				It("should return valid response", func() {
					Expect(resp).ToNot(BeNil())
					Expect(resp.StatusCode).To(Equal(200))
					// TODO: Check resp body
				})
			})

			Context("when request not from the scenario", func() {
				BeforeEach(func() {
					newReq, _ := MakeRequest("POST", "https://secure.api.com/login", emptyHeader, nil)
					resp = subject.Request(newReq, nil)
				})

				It("should return 404", func() {
					Expect(resp).ToNot(BeNil())
					Expect(resp.StatusCode).To(Equal(404))
				})
			})
		})
	})
})
