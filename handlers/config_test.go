package handlers_test

import (
	. "chuck/handlers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"flag"
	"io/ioutil"
)

func defaultBase(folder string) BaseConfig {
	return BaseConfig{
		Address: "127.0.0.1",
		Port:    8123,
		Folder:  folder,
	}
}

var _ = Describe("Config", func() {
	var (
		flags *flag.FlagSet
	)

	BeforeEach(func() {
		flags = flag.NewFlagSet("test", flag.ContinueOnError)
		flags.SetOutput(ioutil.Discard)
	})

	Describe("Recorder", func() {
		var (
			subject  *RecorderConfig
			expected *RecorderConfig
		)

		Context("when args list is empty", func() {
			BeforeEach(func() {
				subject = NewRecorderConfig(flags, []string{}, "rec")

				expected = &RecorderConfig{
					BaseConfig:      defaultBase("rec"),
					CreateNewFolder: true,
					Prevent304:      true,
					OnlyNew:         false,
					LogAsFocused:    false,
					LogRequests:     true,
				}
			})

			It("should return default parameters set", func() {
				Expect(subject).To(Equal(expected))
			})
		})

		Context("when args list contains all args", func() {
			BeforeEach(func() {
				args := []string{"-address=www.google.com", "-port=9999", "-folder=log-99", "-focused",
					"-prevent_304=0", "-new_folder=0", "-requests=0"}
				subject = NewRecorderConfig(flags, args, "rec")

				expected = &RecorderConfig{
					BaseConfig: BaseConfig{
						Address: "www.google.com",
						Port:    9999,
						Folder:  "log-99",
					},
					CreateNewFolder: false,
					Prevent304:      false,
					OnlyNew:         false,
					LogAsFocused:    true,
					LogRequests:     false,
				}
			})

			It("should return parsed parameters set", func() {
				Expect(subject).To(Equal(expected))
			})
		})

		Context("when args list is invalid", func() {
			BeforeEach(func() {
				args := []string{"-address=www.google.com", "-port=9999", "-folder22=log-99"}
				subject = NewRecorderConfig(flags, args, "rec")
			})

			It("should return parsed parameters set", func() {
				Expect(subject).To(BeNil())
			})
		})
	}) // Recorder

	Describe("ScenarioRecorder", func() {
		var (
			subject  *ScenarioRecorderConfig
			expected *ScenarioRecorderConfig
		)

		Context("when args list is empty", func() {
			BeforeEach(func() {
				subject = NewScenarioRecorderConfig(flags, []string{}, "rec")

				expected = &ScenarioRecorderConfig{
					BaseConfig:      defaultBase("rec"),
					CreateNewFolder: false,
					Prevent304:      true,
					OnlyNew:         true,
					LogRequests:     false,
				}
			})

			It("should return default parameters set", func() {
				Expect(subject).To(Equal(expected))
			})
		})

		Context("when args list contains all args", func() {
			BeforeEach(func() {
				args := []string{"-address=www.google.com", "-port=9999", "-folder=log-99", "-prevent_304=0",
					"-new_folder=0", "-requests=1"}
				subject = NewScenarioRecorderConfig(flags, args, "rec")

				expected = &ScenarioRecorderConfig{
					BaseConfig: BaseConfig{
						Address: "www.google.com",
						Port:    9999,
						Folder:  "log-99",
					},
					CreateNewFolder: false,
					Prevent304:      false,
					OnlyNew:         true,
					LogRequests:     true,
				}
			})

			It("should return parsed parameters set", func() {
				Expect(subject).To(Equal(expected))
			})
		})

		Context("when args list is invalid", func() {
			BeforeEach(func() {
				args := []string{"-address=www.google.com", "-port=9999", "-folder22=log-99"}
				subject = NewScenarioRecorderConfig(flags, args, "rec")
			})

			It("should return parsed parameters set", func() {
				Expect(subject).To(BeNil())
			})
		})
	}) // ScenarioRecorder

	Describe("Seeker", func() {
		var (
			subject  *SeekerConfig
			expected *SeekerConfig
		)

		Context("when args list is empty", func() {
			BeforeEach(func() {
				subject = NewSeekerConfig(flags, []string{}, "seek")

				expected = &SeekerConfig{
					BaseConfig: defaultBase("seek"),
				}
			})

			It("should return default parameters set", func() {
				Expect(subject).To(Equal(expected))
			})
		})

		Context("when args list contains all args", func() {
			BeforeEach(func() {
				args := []string{"-address=www.google.com", "-port=9999", "-folder=log-99"}
				subject = NewSeekerConfig(flags, args, "rec")

				expected = &SeekerConfig{
					BaseConfig: BaseConfig{
						Address: "www.google.com",
						Port:    9999,
						Folder:  "log-99",
					},
				}
			})

			It("should return parsed parameters set", func() {
				Expect(subject).To(Equal(expected))
			})
		})

		Context("when args list is invalid", func() {
			BeforeEach(func() {
				args := []string{"-address=www.google.com", "-port=9999", "-folder22=log-99"}
				subject = NewSeekerConfig(flags, args, "rec")
			})

			It("should return parsed parameters set", func() {
				Expect(subject).To(BeNil())
			})
		})
	}) // Seeker

	Describe("ScenarioSeeker", func() {
		var (
			subject  *ScenarioSeekerConfig
			expected *ScenarioSeekerConfig
		)

		Context("when args list is empty", func() {
			BeforeEach(func() {
				subject = NewScenarioSeekerConfig(flags, []string{}, "seek")

				expected = &ScenarioSeekerConfig{
					BaseConfig: defaultBase("seek"),
					Verbose:    true,
				}
			})

			It("should return default parameters set", func() {
				Expect(subject).To(Equal(expected))
			})
		})

		Context("when args list contains all args", func() {
			BeforeEach(func() {
				args := []string{"-address=www.google.com", "-port=9999", "-folder=log-99", "-verbose=0"}
				subject = NewScenarioSeekerConfig(flags, args, "rec")

				expected = &ScenarioSeekerConfig{
					BaseConfig: BaseConfig{
						Address: "www.google.com",
						Port:    9999,
						Folder:  "log-99",
					},
					Verbose: false,
				}
			})

			It("should return parsed parameters set", func() {
				Expect(subject).To(Equal(expected))
			})
		})

		Context("when args list is invalid", func() {
			BeforeEach(func() {
				args := []string{"-address=www.google.com", "-port=9999", "-folder22=log-99"}
				subject = NewScenarioSeekerConfig(flags, args, "rec")
			})

			It("should return parsed parameters set", func() {
				Expect(subject).To(BeNil())
			})
		})
	}) // ScenarioSeeker
})
