package storage_test

import (
	. "chuck/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
)

var _ = Describe("Index", func() {
	var (
		subject Index

		item1 IndexItem
		item2 IndexItem
		item3 IndexItem
		item4 IndexItem
		item5 IndexItem
	)

	Describe("Create empty index", func() {
		BeforeEach(func() {
			item1 = IndexItem{Focused: false, Code: 200, Folder: "r1", Method: "GET", Url: "https://secure.api.com?query=123"}
			item2 = IndexItem{Focused: true, Code: 400, Folder: "r2", Method: "POST", Url: "https://profile.node.com/user/:user_id"}
			item3 = IndexItem{Focused: true, Code: 200, Folder: "r3", Method: "GET", Url: "https://test.aaa.r53.amazon.net:443/secure.amazon.com/v1/authinit?format=json&apikey=*&code=*"}
			item4 = IndexItem{Focused: true, Code: 200, Folder: "r4", Method: "GET", Url: "https://test.net/breadcrumb/offers?orderby=Boosted&top=20&skip=0&breadcrumb=Home/Men/All%20Men&category=mens-view-all&filterby=store%20eq%201&apikey=1111"}
			item5 = IndexItem{Focused: true, Code: 200, Folder: "r5", Method: "GET", Url: "https://test.net/breadcrumb/offers?orderby=Boosted&top=20&skip=0&breadcrumb=Home/Men/All%20Men&category=mens-view-all&filterby=searchcolorfacet%20eq%20'Black'&apikey=1111"}

			subject = NewIndex()
		})

		It("should be empty", func() {
			Expect(subject.Size()).To(Equal(0))
		})

		Describe("Search", func() {
			BeforeEach(func() {
				subject.Add(item1)
				subject.Add(item2)
				subject.Add(item3)
				subject.Add(item4)
			})

			It("should return the correct item (1)", func() {
				p := subject.Find("GET", "https://secure.api.com?query=123")
				Expect(p).To(Equal(&item1))
			})

			It("should return the correct item (2)", func() {
				p := subject.Find("POST", "https://profile.node.com/user/123")
				Expect(p).To(Equal(&item2))
			})

			It("should return the correct item (3)", func() {
				p := subject.Find("GET", "https://test.aaa.r53.amazon.net:443/secure.amazon.com/v1/authinit?format=json&apikey=1111&code=RMJKAZiOIfW8qWbRJlqYL-QoWcF4L8SMFXtlFtkavZU*")
				Expect(p).To(Equal(&item3))
			})

			It("should return the correct item (4)", func() {
				p := subject.Find("GET", "https://test.net/breadcrumb/offers?orderby=Boosted&top=20&skip=0&breadcrumb=Home/Men/All%20Men&category=mens-view-all&filterby=store%20eq%201&apikey=1111")
				Expect(p).To(Equal(&item4))
			})

			It("should return nil for request with additional query param", func() {
				p := subject.Find("GET", "https://test.net/breadcrumb/offers?orderby=Boosted&top=20&skip=0&breadcrumb=Home/Men/All%20Men&category=mens-view-all&filterby=searchcolorfacet%20eq%20'Black'&apikey=1111")
				Expect(p).To(BeNil())
			})

			Describe("Search by additional search param", func() {
				BeforeEach(func() {
					subject.Add(item5)
				})

				It("should return the correct item (5)", func() {
					p := subject.Find("GET", "https://test.net/breadcrumb/offers?orderby=Boosted&top=20&skip=0&breadcrumb=Home/Men/All%20Men&category=mens-view-all&filterby=searchcolorfacet%20eq%20'Black'&apikey=1111")
					Expect(p).To(Equal(&item5))
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
					subject, err = LoadIndex2(root, IndexFileName, false, nil)
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
					subject, err = LoadIndex2(root, IndexFileName, true, nil)
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
