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
		/*log = NewLogger(&cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		})*/
		log = NewLogger(&cli.MockUi{})

		createRequest = func() *http.Request {
			req, _ := http.NewRequest("POST", "https://secure.api.com?query=123", ioutil.NopCloser(bytes.NewBufferString("")))
			return req
		}

		createResponse = func() *http.Response {
			str := "{}"

			resp := &http.Response{
				StatusCode:    200,
				Header:        make(http.Header),
				Body:          ioutil.NopCloser(bytes.NewBufferString(str)),
				ContentLength: int64(len(str)),
			}
			resp.Header.Set("Content-Type", "application/json")
			return resp
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
				subject, err = NewScenarioRecorderWithFs(root, folder, false, log)
				dirExists, _ = root.DirExists(folder)

				path := folder + "/" + subject.Name()
				dirExists, _ = root.DirExists(path)
			})

			It("should return nil error", func() {
				Expect(err).To(BeNil())
			})

			It("should return ScRecorder object", func() {
				Expect(subject).ToNot(BeNil())
			})

			It("should create a recording folder", func() {
				Expect(dirExists).To(BeTrue())
			})
		})

		Context("when createNewFolder is false", func() {
			BeforeEach(func() {
				subject, err = NewScenarioRecorderWithFs(root, folder, false, log)
				dirExists, _ = root.DirExists(folder)
			})

			It("should return nil error", func() {
				Expect(err).To(BeNil())
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

					It("should return nil error for request recording", func() {
						Expect(reqErr).To(BeNil())
					})

					It("should return nil error for response recording", func() {
						Expect(respErr).To(BeNil())
					})
				})
			})
		})
	})
})
