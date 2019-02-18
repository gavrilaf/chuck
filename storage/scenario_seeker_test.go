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

var _ = Describe("ScenarioSeeker", func() {
	var (
		log  Logger
		root *afero.Afero
	)

	BeforeEach(func() {
		log = NewLogger(&cli.MockUi{})

		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		header.Set("Access-Token", "Bearer-12234")

		body := "{}"

		req1, _ := MakeRequest("POST", "https://secure.api.com/login", header, nil)
		req2, _ := MakeRequest("GET", "https://secure.api.com/users/113/on", header, nil)

		resp := MakeResponse2(200, header, body)

		fs := afero.NewMemMapFs()
		root = &afero.Afero{Fs: fs}

		recorder1, _ := NewRecorder(fs, log, "test/folder1/scenario-1", false, false)
		recorder1.SetFocusedMode(true)

		recorder1.RecordRequest(req1, 1)
		recorder1.RecordResponse(resp, 1)

		recorder2, _ := NewRecorder(fs, log, "test/folder2/scenario-2", false, false)
		recorder2.SetFocusedMode(true)

		recorder2.RecordRequest(req2, 1)
		recorder2.RecordResponse(resp, 1)

		recorder3, _ := NewRecorder(fs, log, "test/folder2/scenario-3", false, false)
		recorder3.SetFocusedMode(true)

		recorder3.RecordRequest(req1, 1)
		recorder3.RecordResponse(resp, 1)
	})

	Describe("Open Scenario", func() {
		var (
			err     error
			subject ScenarioSeeker
		)

		BeforeEach(func() {
			subject, err = NewScenarioSeeker(root, log, "test")
		})

		It("should not error occurred", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should create scenario seeker", func() {
			Expect(subject).ToNot(BeNil())
		})

		It("should contain three scenarios", func() {
			Expect(subject.ScenariosCount()).To(Equal(3))
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
					resp, _ = subject.Look("scenario-2", "GET", "https://secure.api.com/users/113/on/update")
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
