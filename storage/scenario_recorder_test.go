package storage_test

import (
	. "chuck/storage"
	. "chuck/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"net/http"
)

var _ = Describe("ScenarioRecorder", func() {
	var (
		log     Logger
		subject ScenarioRecorder
		folder  string
		root    *afero.Afero

		createRequest  func() *http.Request
		createResponse func() *http.Response
	)

	BeforeEach(func() {
		log = NewLogger(&cli.MockUi{})

		header := make(http.Header)

		createRequest = func() *http.Request {
			req, _ := MakeRequest2("POST", "https://secure.api.com?query=123", header, "")
			return req
		}

		createResponse = func() *http.Response {
			return MakeResponse2(200, header, "{}")
		}

		folder = "sc-folder"
		root = &afero.Afero{Fs: afero.NewMemMapFs()}
	})

	Describe("Create Scenario recorder", func() {
		var (
			err       error
			dirExists bool
		)

		Context("when createNewFolder is true", func() {
			BeforeEach(func() {
				subject, err = NewScenarioRecorder(root, log, folder, false)
				dirExists, _ = root.DirExists(folder)

				path := folder + "/" + subject.Name()
				dirExists, _ = root.DirExists(path)
			})

			It("should not error occurred", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return scenario recorder object", func() {
				Expect(subject).ToNot(BeNil())
			})

			It("should create a recording folder", func() {
				Expect(dirExists).To(BeTrue())
			})
		})

		Context("when createNewFolder is false", func() {
			BeforeEach(func() {
				subject, err = NewScenarioRecorder(root, log, folder, false)
				dirExists, _ = root.DirExists(folder)
			})

			It("should not error occurred", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return empty Name", func() {
				Expect(subject.Name()).To(Equal(""))
			})

			It("should create a recording folder", func() {
				Expect(dirExists).To(BeTrue())
			})

			Describe("Activate scenario", func() {
				BeforeEach(func() {
					err = subject.ActivateScenario("scenario-1")
				})

				It("should return nil error", func() {
					Expect(err).To(BeNil())
				})

				Describe("Recording request", func() {
					var (
						reqErr  error
						respErr error
					)
					BeforeEach(func() {
						_, reqErr = subject.RecordRequest(createRequest(), 1)
						_, respErr = subject.RecordResponse(createResponse(), 1)
					})

					It("should not request recording error occurred", func() {
						Expect(reqErr).ToNot(HaveOccurred())
					})

					It("should not response recording error occurred", func() {
						Expect(respErr).ToNot(HaveOccurred())
					})
				})
			})
		})
	})
})
