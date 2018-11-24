package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/gavrilaf/chuck/utils"

	"net/http"
)

var _ = Describe("HeaderCoder", func() {
	var (
		empty  http.Header
		filled http.Header

		buf         []byte
		encodingErr error
		decodingErr error

		decoded http.Header
	)

	BeforeEach(func() {
		empty = http.Header{}
		filled = http.Header{
			"Access-Control-Allow-Origin": {"*"},
			"Content-Type":                {"text/html; charset=utf-8"},
			"Content-Length":              {"68137"},
		}
	})

	Describe("encode/decode headers", func() {
		Context("empty", func() {
			BeforeEach(func() {
				buf, encodingErr = EncodeHeaders(empty)
				decoded, decodingErr = DecodeHeaders(buf)
			})

			It("should return nil errors", func() {
				Expect(encodingErr).To(BeNil())
				Expect(decodingErr).To(BeNil())
			})

			It("should return well-formed buffer", func() {
				Expect(string(buf)).To(Equal("{}"))
			})

			It("should return correcr header", func() {
				Expect(empty).To(Equal(decoded))
			})
		})

		Context("filled", func() {
			BeforeEach(func() {
				buf, encodingErr = EncodeHeaders(filled)
				decoded, decodingErr = DecodeHeaders(buf)
			})

			It("should return nil errors", func() {
				Expect(encodingErr).To(BeNil())
				Expect(decodingErr).To(BeNil())
			})

			It("should return correcr header", func() {
				Expect(filled).To(Equal(decoded))
			})
		})
	})
})
