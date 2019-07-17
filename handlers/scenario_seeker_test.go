package handlers_test

import (
	. "chuck/handlers"
	. "chuck/storage"
	. "chuck/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
)

var _ = Describe("ScenarioSeeker handler", func() {
	var (
		log Logger
		fs  afero.Fs

		respRecorder *httptest.ResponseRecorder
		resp         *http.Response
		emptyHeader  http.Header
		jsonHeader   http.Header
		jsonBody     string

		err     error
		subject ProxyHandler
	)

	BeforeEach(func() {
		log = NewLogger(cli.NewMockUi())
		fs = afero.NewMemMapFs()

		emptyHeader = make(http.Header)

		jsonHeader = make(http.Header)
		jsonHeader.Set("Content-Type", "application/json")

		jsonBody = "{}"
	})

	Describe("open scenario seeker", func() {

		BeforeEach(func() {
			resp = MakeResponse2(200, jsonHeader, jsonBody)

			req1, _ := MakeRequest("POST", "https://secure.api.com/login", emptyHeader, nil)
			req2, _ := MakeRequest("GET", "https://secure.api.com/users", emptyHeader, nil)

			recorder1, _ := NewRecorder(fs, log, "test/scenario-1", false, false, true)
			defer recorder1.Close()

			recorder1.SetFocusedMode(true)

			recorder1.RecordRequest(req1, 1)
			recorder1.RecordResponse(resp, 1)

			recorder2, _ := NewRecorder(fs, log, "test/scenario-2", false, false, true)
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
				req, _ := MakeRequest("PUT", "https://127.0.0.1/scenario/scenario-3/device-1-id/no", emptyHeader, nil)
				respRecorder = httptest.NewRecorder()
				subject.NonProxyHandler(respRecorder, req)
			})

			It("should return 404", func() {
				Expect(respRecorder.Code).To(Equal(404))
			})
		})

		Context("when activate existing scenario", func() {
			BeforeEach(func() {
				req, _ := MakeRequest("PUT", "https://127.0.0.1/scenario/scenario-1/device-1-id/no", emptyHeader, nil)
				respRecorder = httptest.NewRecorder()
				subject.NonProxyHandler(respRecorder, req)
			})

			It("should activate scenario", func() {
				Expect(respRecorder.Code).To(Equal(200))
			})

			Context("when handle request from the scenario", func() {
				BeforeEach(func() {
					header := make(http.Header)
					header.Set(ScenarioIdHeader, "device-1-id")
					req, _ := MakeRequest("POST", "https://secure.api.com/login", header, nil)
					resp = subject.Request(req, nil)
				})

				It("should return valid response", func() {
					Expect(resp).ToNot(BeNil())
					Expect(resp.StatusCode).To(Equal(200))
					Expect(resp.Header).To(Equal(jsonHeader))

					var buf bytes.Buffer
					buf.ReadFrom(resp.Body)
					s := strings.TrimSpace(string(buf.Bytes()))
					Expect(s).To(Equal("{}"))
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

		Describe("service requests", func() {
			Describe("execute script request", func() {
				// TODO: implement it
			})

			Describe("execute reload request", func() {
				BeforeEach(func() {
					// record additional scenario
					resp = MakeResponse2(200, jsonHeader, jsonBody)
					req, _ := MakeRequest("POST", "https://secure.api.com/login", emptyHeader, nil)

					recorder, _ := NewRecorder(fs, log, "test/scenario-3", false, false, true)
					defer recorder.Close()
					recorder.SetFocusedMode(true)
					recorder.RecordRequest(req, 1)
					recorder.RecordResponse(resp, 1)

					// reload scenario seeker
					reloadReq, _ := MakeRequest("PUT", "https://127.0.0.1/scenarios/reload", emptyHeader, nil)
					respRecorder = httptest.NewRecorder()
					subject.NonProxyHandler(respRecorder, reloadReq)
				})

				It("should return 200", func() {
					Expect(respRecorder.Code).To(Equal(200))
				})

				Describe("activate new scenario", func() {
					BeforeEach(func() {
						req, _ := MakeRequest("PUT", "https://127.0.0.1/scenario/scenario-3/device-1-id/no", emptyHeader, nil)
						respRecorder = httptest.NewRecorder()
						subject.NonProxyHandler(respRecorder, req)
					})

					It("should return 200", func() {
						Expect(respRecorder.Code).To(Equal(200))
					})
				})
			})
		})
	})
})
