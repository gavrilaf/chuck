package handlers_test

import (
	. "chuck/handlers"
	. "chuck/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
)

var _ = Describe("Recorder handler", func() {
	var (
		log Logger
		fs  afero.Fs

		err     error
		subject ProxyHandler
	)

	BeforeEach(func() {
		log = NewLogger(cli.NewMockUi())
		fs = afero.NewMemMapFs()
	})

	Describe("open recorder", func() {
		BeforeEach(func() {
			cfg := &RecorderConfig{
				BaseConfig: BaseConfig{
					Folder: "test",
				},
				CreateNewFolder: true,
				Prevent304:      true,
				LogAsFocused:    false,
			}

			subject, err = NewRecorderHandler(cfg, fs, log)
		})

		It("should no error occured", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return recorder proxy handler", func() {
			Expect(subject).ToNot(BeNil())
		})
	})
})
