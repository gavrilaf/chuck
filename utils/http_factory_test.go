package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gavrilaf/chuck/utils"

	"net/http"
)

var _ = Describe("HttpFactory", func() {
	Describe("Request", func() {
		var (
			subj *http.Request
			err  error

			method string
			url    string
			header http.Header
			body   string
		)

		BeforeEach(func() {
			method = "POST"
			url = "www.google.com"
			header = make(http.Header)
			header.Set("Content-Type", "application/json")
			body = "{}"

			subj, err = MakeRequest2(method, url, header, body)
		})

		It("should not error occurred", func() {
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should create correct request", func() {
			Expect(subj.Method).To(Equal(method))
			Expect(subj.URL.String()).To(Equal(url))
			Expect(subj.Header).To(Equal(header))

			reqBody, _ := DumpReqBody(subj)
			Expect(reqBody).To(Equal([]byte(body)))
		})
	})

	Describe("Response", func() {
		var (
			subj *http.Response

			code   int
			header http.Header
			body   string
		)

		BeforeEach(func() {
			code = 200
			header = make(http.Header)
			header.Set("Content-Type", "application/json")
			body = "{}"

			subj = MakeResponse2(code, header, body)
		})

		It("should create correct response", func() {
			Expect(subj.StatusCode).To(Equal(code))
			Expect(subj.Header).To(Equal(header))

			reqBody, _ := DumpRespBody(subj)
			Expect(reqBody).To(Equal([]byte(body)))
		})
	})
})
