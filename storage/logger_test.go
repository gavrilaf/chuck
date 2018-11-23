package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/spf13/afero"
)

var _ = Describe("Logger", func() {
	var (
		subject ReqLogger
		afr     *afero.Afero
	)

	BeforeEach(func() {
		afr = &afero.Afero{Fs: afero.NewMemMapFs()}
		subject = NewLoggerWithFs(afr)
	})

	Describe("Start logger", func() {
		var (
			err         error
			dirExists   bool
			indexExists bool
		)

		BeforeEach(func() {
			err = subject.Start()

			path := subject.Name()
			dirExists, _ = afr.DirExists(path)
			indexExists, _ = afr.Exists(path + "/" + "index.txt")
		})

		It("should return nil error", func() {
			Expect(err).To(BeNil())
		})

		It("should create a logger folder", func() {
			Expect(dirExists).To(BeTrue())
		})

		It("should create an index file", func() {
			Expect(indexExists).To(BeTrue())
		})
	})
})
