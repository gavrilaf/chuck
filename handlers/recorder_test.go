package handlers_test

import (
	. "chuck/handlers"
	. "chuck/utils"

	"bufio"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"gopkg.in/elazarl/goproxy.v1"
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
				OnlyNew:         false,
				CreateNewFolder: false,
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

		Describe("Record requests", func() {
			var (
				sendRequest func(url string, ctx int64)
			)

			BeforeEach(func() {
				sendRequest = func(url string, ctx int64) {
					header := make(http.Header)
					context := &goproxy.ProxyCtx{Session: ctx}

					req, _ := MakeRequest("GET", url, header, nil)
					resp := MakeResponse2(200, header, "{\"code\" = 111}")

					subject.Request(req, context)
					subject.Response(resp, context)
				}

				sendRequest("http://test.net/users", 101)
				sendRequest("http://test.net/sessions", 102)
				sendRequest("http://test.net/events", 103)
				sendRequest("http://test.net/info", 104)
				sendRequest("http://test.net/users", 105)
				sendRequest("http://test.net/info", 106)
			})

			Describe("Read index file directly", func() {
				var (
					lines []string
				)

				BeforeEach(func() {
					fi, _ := fs.Open("test/index.txt")
					defer fi.Close()

					scanner := bufio.NewScanner(fi)
					for scanner.Scan() {
						lines = append(lines, scanner.Text())
					}
				})

				It("should contains expected lines", func() {
					expected := []string{
						"N,\t200,\tr_1,\tGET,\thttp://test.net/users",
						"N,\t200,\tr_2,\tGET,\thttp://test.net/sessions",
						"N,\t200,\tr_3,\tGET,\thttp://test.net/events",
						"N,\t200,\tr_4,\tGET,\thttp://test.net/info",
						"N,\t200,\tr_5,\tGET,\thttp://test.net/users",
						"N,\t200,\tr_6,\tGET,\thttp://test.net/info",
					}

					Expect(expected).To(Equal(lines))
				})
			})
		})
	})
})
