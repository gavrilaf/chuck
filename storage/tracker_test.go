package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/gavrilaf/chuck/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mitchellh/cli"
	"net/http"
)

var _ = Describe("Tracker", func() {
	var (
		log        Logger
		req        *http.Request
		resp       *http.Response
		reqResult  *PendingRequest
		respResult *PendingRequest
		err        error

		subject Tracker
	)

	BeforeEach(func() {
		header := make(http.Header)

		req, _ = MakeRequest2("GET", "www.google.com", header, "")
		resp = MakeResponse2(200, header, "")

		log = NewLogger(cli.NewMockUi())

		subject = NewTracker(1, log)
	})

	Describe("track request", func() {
		BeforeEach(func() {
			reqResult, err = subject.RecordRequest(req, 10)
		})

		It("should no error occured", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return pending request", func() {
			Expect(reqResult).ToNot(BeNil())
		})

		It("should contain one pending request", func() {
			Expect(subject.PendingCount()).To(Equal(1))
		})

		It("should pending request has correct properties", func() {
			Expect(reqResult.Id).To(Equal(int64(1)))
			Expect(reqResult.Method).To(Equal("GET"))
			Expect(reqResult.Url).To(Equal("www.google.com"))
			Expect(reqResult.Started).ToNot(BeZero())
		})

		Describe("track response", func() {
			Context("when session id is correct", func() {
				BeforeEach(func() {
					respResult, err = subject.RecordResponse(resp, 10)
				})

				It("should no error occured", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should return pending request", func() {
					Expect(respResult).ToNot(BeNil())
				})

				It("should contain no pending requests", func() {
					Expect(subject.PendingCount()).To(BeZero())
				})

				It("should return the same object with RecordRequest", func() {
					Expect(reqResult).To(Equal(respResult))
				})
			})

			Context("when session id is unknown", func() {
				BeforeEach(func() {
					respResult, err = subject.RecordResponse(resp, 101)
				})

				It("should error occured", func() {
					Expect(err).To(MatchError(ErrRequestNotFound))
				})

				It("should not return pending request", func() {
					Expect(respResult).To(BeNil())
				})

				It("should still contain one pending request", func() {
					Expect(subject.PendingCount()).To(Equal(1))
				})

			})
		})
	})
})
