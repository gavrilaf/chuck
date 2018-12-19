package handlers_test

import (
	. "github.com/gavrilaf/chuck/handlers"
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/gavrilaf/chuck/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Scenario", func() {
	var (
		log      Logger
		root     *afero.Afero
		scSeeker ScenarioSeeker

		createRequest  func(method string, url string) *http.Request
		createResponse func() *http.Response
	)

	BeforeEach(func() {
		log = NewLogger(cli.NewMockUi())

		createRequest = func(method string, url string) *http.Request {
			req, _ := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewBufferString("")))
			return req
		}

		createResponse = func() *http.Response {
			respBody := "{}"

			resp := &http.Response{
				StatusCode:    200,
				Header:        make(http.Header),
				Body:          ioutil.NopCloser(bytes.NewBufferString(respBody)),
				ContentLength: int64(len(respBody)),
			}
			return resp
		}

		fs := afero.NewMemMapFs()
		root = &afero.Afero{Fs: fs}

		recorder1, _ := NewRecorderWithFs(fs, "test/scenario-1", false, log)
		recorder1.SetFocusedMode(true)

		recorder1.RecordRequest(createRequest("POST", "https://secure.api.com/login"), 1)
		recorder1.RecordResponse(createResponse(), 1)

		recorder2, _ := NewRecorderWithFs(fs, "test/scenario-2", false, log)
		recorder2.SetFocusedMode(true)

		recorder2.RecordRequest(createRequest("GET", "https://secure.api.com/users"), 1)
		recorder2.RecordResponse(createResponse(), 1)

		scSeeker, _ = NewScenarioSeekerWithFs(root, "test", log)
	})

	Describe("Open scenario proxy handler", func() {
		var (
			subject  ProxyHandler
			recorder *httptest.ResponseRecorder
			resp     *http.Response
		)

		BeforeEach(func() {
			subject = NewScenarioHandlerWithSeeker(scSeeker, log)
		})

		It("should not be null", func() {
			Expect(subject).ToNot(BeNil())
		})

		Context("when activate unknown scenario", func() {
			BeforeEach(func() {
				req := createRequest("PUT", "https://127.0.0.1/scenario/scenario-111/scenario-111-id/no")
				recorder = httptest.NewRecorder()
				subject.NonProxyHandler(recorder, req)
			})

			It("should activate scenario", func() {
				Expect(recorder.Code).To(Equal(404))
			})
		})

		Context("when activate existing scenario", func() {
			BeforeEach(func() {
				req := createRequest("PUT", "https://127.0.0.1/scenario/scenario-1/scenario-1-id/no")
				recorder = httptest.NewRecorder()
				subject.NonProxyHandler(recorder, req)
			})

			It("should activate scenario", func() {
				Expect(recorder.Code).To(Equal(200))
			})

			Context("when request from the scenario", func() {
				BeforeEach(func() {
					req := createRequest("POST", "https://secure.api.com/login")
					req.Header = make(http.Header)
					req.Header.Set(AADHIIdentifier, "scenario-1-id")
					resp = subject.Request(req, nil)
				})

				It("should return response", func() {
					Expect(resp).ToNot(BeNil())
				})
			})

			Context("when request not from the scenario", func() {
				BeforeEach(func() {
					req := createRequest("POST", "https://secure.api.com/login")
					resp = subject.Request(req, nil)
				})

				It("should return nil", func() {
					Expect(resp).To(BeNil())
				})
			})
		})
	})
})
