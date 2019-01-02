package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/gavrilaf/chuck/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"net/http"
)

var _ = Describe("ScenarioSeeker", func() {
	var (
		log  Logger
		root *afero.Afero

		createRequest  func(method string, url string) *http.Request
		createResponse func() *http.Response
	)

	BeforeEach(func() {
		log = NewLogger(&cli.MockUi{})

		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		header.Set("Access-Token", "Bearer-12234")

		body := "{}"

		createRequest = func(method string, url string) *http.Request {
			req, _ := MakeRequest2(method, url, header, body)
			return req
		}

		createResponse = func() *http.Response {
			return MakeResponse2(200, header, body)
		}

		fs := afero.NewMemMapFs()
		root = &afero.Afero{Fs: fs}

		recorder1, _ := NewRecorderWithFs(fs, "test/scenario-1", false, false, log)
		recorder1.SetFocusedMode(true)

		recorder1.RecordRequest(createRequest("POST", "https://secure.api.com/login"), 1)
		recorder1.RecordResponse(createResponse(), 1)

		recorder2, _ := NewRecorderWithFs(fs, "test/scenario-2", false, false, log)
		recorder2.SetFocusedMode(true)

		recorder2.RecordRequest(createRequest("GET", "https://secure.api.com/users/113/on"), 1)
		recorder2.RecordResponse(createResponse(), 1)
	})

	Describe("Open Scenario", func() {
		var (
			err     error
			subject ScenarioSeeker
		)

		BeforeEach(func() {
			subject, err = NewScenarioSeekerWithFs(root, "test", log)
		})

		It("should not error occurred", func() {
			Expect(err).ToNot(HaveOccurred())
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
				err  error
			)

			Context("when request from scenario 1", func() {
				BeforeEach(func() {
					resp, _ = subject.Look("scenario-1", "POST", "https://secure.api.com/login")
				})

				It("should find response", func() {
					Expect(resp).ToNot(BeNil())
				})
			})

			Context("when request from scenario 2; looking using prefix", func() {
				BeforeEach(func() {
					resp, _ = subject.Look("scenario-2", "GET", "https://secure.api.com/users")
				})

				It("should find response", func() {
					Expect(resp).ToNot(BeNil())
				})
			})

			Context("when request from unknown scenarion", func() {

				BeforeEach(func() {
					resp, err = subject.Look("scenarion-6666", "GET", "https://secure.api.com/users")
				})

				It("should return nil response", func() {
					Expect(resp).To(BeNil())
				})

				It("should return error", func() {
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})
})
