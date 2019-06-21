package handlers_test

import (
	. "chuck/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io/ioutil"
	"net/http"
)

var _ = Describe("Helpers", func() {
	var (
		createRequest func(method string, url string) *http.Request
	)

	BeforeEach(func() {
		createRequest = func(method string, url string) *http.Request {
			req, _ := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewBufferString("")))
			return req
		}
	})

	Describe("Prevent 304 http answer", func() {
		var (
			req *http.Request
		)

		BeforeEach(func() {
			req = createRequest("GET", "www.google.com")
			Prevent304HttpAnswer(req)
		})

		It("should update headers", func() {
			Expect(req.Header.Get("If-Modified-Since")).To(Equal("off"))
			Expect(req.Header.Get("Last-Modified")).To(Equal(""))
		})
	})

	Describe("Get scenario id", func() {
		var (
			req *http.Request
		)

		BeforeEach(func() {
			req = createRequest("GET", "www.google.com")
		})

		Context("when request contains scenario id", func() {
			BeforeEach(func() {
				req.Header.Set(ScenarioIdHeader, "123456")
			})

			It("should return scenario id", func() {
				id := GetScenarioId(req)
				Expect(id).To(Equal("123456"))
			})
		})

		Context("when request does not contain scenario id", func() {
			It("should return empty string", func() {
				id := GetScenarioId(req)
				Expect(id).To(Equal(""))
			})
		})
	})
})
