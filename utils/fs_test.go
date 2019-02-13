package utils_test

import (
	. "chuck/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
	. "path"
)

var _ = Describe("FS", func() {
	var (
		root       *afero.Afero
		name, path string
		err        error
	)

	BeforeEach(func() {
		fs := afero.NewMemMapFs()
		root = &afero.Afero{Fs: fs}
	})

	Describe("PrepareStorageFolder", func() {
		Context("when folder is exists and createNewFolder is false", func() {
			BeforeEach(func() {
				root.Mkdir("folder-1", 0777)

				name, path, err = PrepareStorageFolder(root, "folder-1", false)
			})

			It("should no error occured", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return '' and 'folder' path", func() {
				Expect(name).To(BeEmpty())
				Expect(path).To(Equal("folder-1"))
			})
		})

		Context("when folder is exists and createNewFolder is true", func() {
			BeforeEach(func() {
				root.Mkdir("folder-2", 0777)

				name, path, err = PrepareStorageFolder(root, "folder-2", true)
			})

			It("should no error occured", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create new folder; should return 'unique-name' and 'folder/unique_name' path", func() {
				expected := Join("folder-2", name)
				exists, _ := root.DirExists(expected)

				Expect(exists).To(BeTrue())
				Expect(name).ToNot(BeEmpty())
				Expect(path).To(Equal(expected))
			})
		})

		Context("when folder is not exists and createNewFolder is false", func() {
			BeforeEach(func() {
				name, path, err = PrepareStorageFolder(root, "folder-3", false)
			})

			It("should no error occured", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create folder; should return '' and 'folder' path", func() {
				exists, _ := root.DirExists("folder-3")

				Expect(exists).To(BeTrue())
				Expect(name).To(BeEmpty())
				Expect(path).To(Equal("folder-3"))
			})
		})

		Context("when folder is not exists and createNewFolder is true", func() {
			BeforeEach(func() {
				name, path, err = PrepareStorageFolder(root, "folder-4", true)
			})

			It("should no error occured", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create folder and 'unique-name' folder inside; should return 'unique-name' and 'folder/unique_name' path", func() {
				expected := Join("folder-4", name)
				exists, _ := root.DirExists(expected)

				Expect(exists).To(BeTrue())
				Expect(name).ToNot(BeEmpty())
				Expect(path).To(Equal(expected))
			})
		})
	})
})
