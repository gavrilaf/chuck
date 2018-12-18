package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IndexItem", func() {
	Describe("Format item", func() {
		var (
			itemStr        string
			focusedItemStr string
		)

		BeforeEach(func() {
			itemStr = FormatIndexItem("GET", "https://secure.api.com?query=123", 200, "r1", false)
			focusedItemStr = FormatIndexItem("PUT", "https://my.api.com/login", 500, "r2", true)
		})

		It("should create focused line", func() {
			Expect(focusedItemStr).To(Equal("F,\t500,\tr2,\tPUT,\thttps://my.api.com/login"))
		})

		It("should create unfocused line", func() {
			Expect(itemStr).To(Equal("N,\t200,\tr1,\tGET,\thttps://secure.api.com?query=123"))
		})

		Describe("Parsing", func() {
			var (
				item        *IndexItem
				focusedItem *IndexItem

				expectedItem        IndexItem
				expectedFocusedItem IndexItem
			)

			BeforeEach(func() {
				expectedItem = IndexItem{Focused: false, Code: 200, Folder: "r1", Method: "GET", Url: "https://secure.api.com?query=123"}
				expectedFocusedItem = IndexItem{Focused: true, Code: 500, Folder: "r2", Method: "PUT", Url: "https://my.api.com/login"}

				item = ParseIndexItem(itemStr)
				focusedItem = ParseIndexItem(focusedItemStr)
			})

			It("should parse focused item", func() {
				Expect(*focusedItem).To(Equal(expectedFocusedItem))
			})

			It("should parse unfocused item", func() {
				Expect(*item).To(Equal(expectedItem))
			})
		})
	})
})
