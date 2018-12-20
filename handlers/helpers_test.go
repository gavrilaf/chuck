package handlers_test

import (
	. "github.com/gavrilaf/chuck/handlers"
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

	Describe("Parse activate scenario url", func() {
		var (
			subj *ActivateScenario
		)

		Context("when url recognized as scenario url", func() {
			BeforeEach(func() {
				req := createRequest("GET", "https://127.0.0.1/scenario/scenario-1/scenario-1-id/no")
				subj = ParseActivateScenarioRequest(req)
			})

			It("should return scenario name & id", func() {
				sc := &ActivateScenario{Scenario: "scenario-1", Id: "scenario-1-id"}
				Expect(subj).To(Equal(sc))
			})
		})

		Context("when url is not recognized as scenario url", func() {
			BeforeEach(func() {
				req := createRequest("GET", "www.google.com")
				subj = ParseActivateScenarioRequest(req)
			})

			It("should return nil", func() {
				Expect(subj).To(BeNil())
			})
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
