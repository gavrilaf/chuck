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

var _ = Describe("Seeker", func() {
	var (
		log Logger
		fs  afero.Fs

		err     error
		subject Seeker
	)

	BeforeEach(func() {
		log = NewLogger(cli.NewMockUi())
		fs = afero.NewMemMapFs()
	})

	Context("Open Seeker on nonexisting folder", func() {
		BeforeEach(func() {
			subject, err = NewSeeker(fs, "test-123")
		})

		It("should error occurred", func() {
			Expect(err).To(MatchError("Folder test-123 doesn't exists"))
		})

		It("should return nil Seeker", func() {
			Expect(subject).To(BeNil())
		})

	})

	Context("Open Seeker on folder with index", func() {
		var (
			header   http.Header
			respBody string
		)

		BeforeEach(func() {
			header = make(http.Header)
			header.Set("Content-Type", "application/json")
			header.Set("Access-Token", "Bearer-12234")

			respBody = `{"colors": []}`

			req1, _ := MakeRequest("POST", "https://secure.api.com/login", header, nil)
			req2, _ := MakeRequest("GET", "https://secure.api.com/users/*", header, nil)

			emptyHeader := make(http.Header)
			reqEmpty, _ := MakeRequest("GET", "www.google.com", emptyHeader, nil)
			respEmpty := MakeResponse(200, emptyHeader, nil, 0)

			resp := MakeResponse2(200, header, respBody)

			recorder, _ := NewRecorder(fs, log, "test", false, false, true, false)

			recorder.RecordRequest(req1, 1)
			recorder.RecordResponse(resp, 1)

			recorder.SetFocusedMode(true)

			recorder.RecordRequest(req2, 2)
			recorder.RecordResponse(resp, 2)

			// empty header/body
			recorder.RecordRequest(reqEmpty, 3)
			recorder.RecordResponse(respEmpty, 3)

			subject, err = NewSeeker(fs, "test")
		})

		It("should not error occurred", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should create Seeker", func() {
			Expect(subject).ToNot(BeNil())
		})

		It("should contains two items", func() {
			Expect(subject.Count()).To(Equal(2))
		})

		Describe("looking for request", func() {
			var (
				resp *http.Response
			)

			Context("when request logged as focused", func() {
				BeforeEach(func() {
					resp, _ = subject.Look("GET", "https://secure.api.com/users/678/off")
				})

				It("should return request", func() {
					Expect(resp).ToNot(BeNil())
				})

				It("should response has correct headers", func() {
					Expect(resp.Header).To(Equal(header))
				})

				It("should response has correct body", func() {
					expected := []byte("{\n\t\"colors\": []\n}\n")
					buf, _ := DumpRespBody(resp)
					Expect(buf).To(Equal(expected))
				})
			})

			Context("when request/response has empty header & body", func() {
				BeforeEach(func() {
					resp, _ = subject.Look("GET", "www.google.com")
				})

				It("should return request", func() {
					Expect(resp).ToNot(BeNil())
				})

				It("should response has empty header", func() {
					Expect(len(resp.Header)).To(BeZero())
				})

				It("should response has empty body", func() {
					buf, _ := DumpRespBody(resp)
					Expect(len(buf)).To(BeZero())
				})
			})

			Context("when request logged as unfocused", func() {
				BeforeEach(func() {
					resp, _ = subject.Look("POST", "https://secure.api.com/login")
				})

				It("should return nil", func() {
					Expect(resp).To(BeNil())
				})
			})
		})
	})
})
