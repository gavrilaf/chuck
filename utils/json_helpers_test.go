package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "chuck/utils"

	"net/http"
)

func createResp(content string) *http.Response {
	header := make(http.Header)
	header.Set("Content-Type", content)
	resp := &http.Response{
		Header: header,
	}
	return resp
}

var _ = Describe("JsonHelpers", func() {
	Describe("IsRespHasJsonContent", func() {
		It("Should return true", func() {
			Expect(IsRespHasJsonContent(createResp("application/json"))).To(BeTrue())
			Expect(IsRespHasJsonContent(createResp("application/json, encoding=utf-8"))).To(BeTrue())
			Expect(IsRespHasJsonContent(createResp("encoding=utf-8,application/json"))).To(BeTrue())
			Expect(IsRespHasJsonContent(createResp("my.dummy.service/json"))).To(BeTrue())
		})
	})

	Describe("FormatJson", func() {
		It("Should return pretty printed json", func() {
			s := []byte("{\"var1\": 1, \"var2\": \"2\"}")
			expected := []byte("{\n\t\"var1\": 1,\n\t\"var2\": \"2\"\n}\n")

			Expect(FormatJson(s)).To(Equal(expected))
		})

		It("Should return original sequence", func() {
			s := []byte("just text")
			expected := []byte("just text")

			Expect(FormatJson(s)).To(Equal(expected))
		})
	})
})
