package handlers_test

import (
	. "chuck/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"io/ioutil"
	"net/http"
)

var _ = Describe("Service calls detector", func() {
	var (
		createRequest func(method string, url string) *http.Request
	)

	BeforeEach(func() {
		createRequest = func(method string, url string) *http.Request {
			req, _ := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewBufferString("")))
			return req
		}
	})

	Describe("Detect service request", func() {
		var reqType int

		Context("when url recognized as activate scenario url", func() {
			BeforeEach(func() {
				req := createRequest("POST", "https://127.0.0.1/scenario/scenario-1/scenario-1-id/no")
				reqType = DetectServiceRequest(req)
			})

			It("should return activate scenario type", func() {
				Expect(reqType).To(Equal(ServiceReq_ActivateScenario))
			})
		})

		Context("when url recognized as run script url", func() {
			BeforeEach(func() {
				req := createRequest("POST", "https://127.0.0.1/script/delete-app/run")
				reqType = DetectServiceRequest(req)
			})

			It("should return activate scenario type", func() {
				Expect(reqType).To(Equal(ServiceReq_ExecuteScript))
			})
		})

		Context("when url is not recognized as service url", func() {
			BeforeEach(func() {
				req := createRequest("POST", "https://127.0.0.1/auth/v1/verifier=2738438")
				reqType = DetectServiceRequest(req)
			})

			It("should return activate scenario type", func() {
				Expect(reqType).To(Equal(ServiceReq_None))
			})
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
})
