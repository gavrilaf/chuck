package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
)

var _ = Describe("Index", func() {
	var (
		subject Index

		item1 IndexItem
		item2 IndexItem
	)

	Describe("Create empty index", func() {
		BeforeEach(func() {
			item1 = IndexItem{Focused: false, Code: 200, Folder: "r1", Method: "GET", Url: "https://secure.api.com?query=123"}
			item2 = IndexItem{Focused: true, Code: 400, Folder: "r2", Method: "POST", Url: "https://profile.node.com/user/123"}

			subject = NewIndex()
		})

		It("should be empty", func() {
			Expect(subject.Size()).To(Equal(0))
		})

		Describe("Search", func() {
			BeforeEach(func() {
				subject.Add(item1)
				subject.Add(item2)
			})

			Context("using Equality option", func() {
				It("should return correct object", func() {
					p := subject.Find("GET", "https://secure.api.com?query=123", SEARCH_EQ)
					Expect(p).To(Equal(&item1))
				})
			})

			Context("using SubStr option", func() {
				It("should return correct object", func() {
					p := subject.Find("POST", "https://profile.node.com/user", SEARCH_SUBSTR)
					Expect(p).To(Equal(&item2))
				})
			})
		})

		Describe("Load index", func() {
			var (
				err  error
				root *afero.Afero
			)

			BeforeEach(func() {
				fs := afero.NewMemMapFs()
				root = &afero.Afero{Fs: fs}

				fp, _ := root.Create(IndexFileName)
				fp.WriteString(item1.Format() + "\n")
				fp.WriteString(item2.Format() + "\n")
				fp.Close()
			})

			Context("when focused is false", func() {
				BeforeEach(func() {
					subject, err = LoadIndex2(root, IndexFileName, false)
				})

				It("should return nil error", func() {
					Expect(err).To(BeNil())
				})

				It("should contain 2 items", func() {
					Expect(subject.Size()).To(Equal(2))
				})

				It("should contain correct items", func() {
					Expect(subject.Get(0)).To(Equal(item1))
					Expect(subject.Get(1)).To(Equal(item2))
				})
			})

			Context("when focused is true", func() {
				BeforeEach(func() {
					subject, err = LoadIndex2(root, IndexFileName, true)
				})

				It("should contain 1 item", func() {
					Expect(subject.Size()).To(Equal(1))
				})

				It("should contains correct items", func() {
					Expect(subject.Get(0)).To(Equal(item2))
				})
			})
		})
	})
})
