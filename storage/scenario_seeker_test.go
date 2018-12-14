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
	//"os"
)

var _ = Describe("Scenario", func() {
	var (
		log  Logger
		root *afero.Afero

		createRequest  func(method string, url string) *http.Request
		createResponse func() *http.Response
	)

	BeforeEach(func() {
		/*log = NewLogger(&cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		})*/
		log = NewLogger(&cli.MockUi{})

		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		header.Set("Access-Token", "Bearer-12234")

		respBody := "{}"

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

		recorder1, _ := NewRecorderWithFs("test/scenario-1", false, fs, log)
		recorder1.SetFocusedMode(true)

		recorder1.RecordRequest(createRequest("POST", "https://secure.api.com/login"), 1)
		recorder1.RecordResponse(createResponse(), 1)

		recorder2, _ := NewRecorderWithFs("test/scenario-2", false, fs, log)
		recorder2.SetFocusedMode(true)

		recorder2.RecordRequest(createRequest("GET", "https://secure.api.com/users"), 1)
		recorder2.RecordResponse(createResponse(), 1)
	})

	Describe("Open Scenario", func() {
		var (
			err     error
			subject ScenarioSeeker
		)

		BeforeEach(func() {
			subject, err = NewScenarioSeekerWithFs("test", root, log)
		})

		It("should return nil error", func() {
			Expect(err).To(BeNil())
		})

		It("should create scenario seeker", func() {
			Expect(subject).ToNot(BeNil())
		})

		Describe("checking if scenario exists", func() {
			var (
				exists    bool
				notExists bool
			)

			BeforeEach(func() {
				exists = subject.IsScenarioExists("scenario-1")
				notExists = subject.IsScenarioExists("scenario-1111")
			})

			It("should return correct values", func() {
				Expect(exists).To(BeTrue())
				Expect(notExists).ToNot(BeTrue())
			})
		})

		Describe("looking for request", func() {
			var (
				resp *http.Response
			)

			Context("when request from scenario 1", func() {
				BeforeEach(func() {
					resp = subject.Look("scenario-1", "POST", "https://secure.api.com/login")
				})

				It("should find response", func() {
					Expect(resp).ToNot(BeNil())
				})
			})

			Context("when request from scenario 2", func() {
				BeforeEach(func() {
					resp = subject.Look("scenario-2", "GET", "https://secure.api.com/users")
				})

				It("should find response", func() {
					Expect(resp).ToNot(BeNil())
				})
			})

			Context("when request from unknown scenarion", func() {
				BeforeEach(func() {
					resp = subject.Look("scenarion-6666", "GET", "https://secure.api.com/users")
				})

				It("should return nil", func() {
					Expect(resp).To(BeNil())
				})
			})
		})
	})
})
