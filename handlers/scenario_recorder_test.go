package handlers_test

import (
	. "chuck/handlers"
	"chuck/storage"
	. "chuck/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"gopkg.in/elazarl/goproxy.v1"
)

var _ = Describe("ScenarioRecorder handler", func() {
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

	Describe("open scenario recorder", func() {
		BeforeEach(func() {
			cfg := &ScenarioRecorderConfig{
				BaseConfig: BaseConfig{
					Folder: "test",
				},
				OnlyNew:         true,
				LogRequests:     false,
				CreateNewFolder: false,
				Prevent304:      true,
			}

			subject, err = NewScenarioRecorderHandler(cfg, fs, log)
		})

		It("should no error occured", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return recorder proxy handler", func() {
			Expect(subject).ToNot(BeNil())
		})

		Describe("Activate scenario 1", func() {
			var (
				respRecorder *httptest.ResponseRecorder
				sendRequest  func(url string, ctx int64)

				index    storage.Index
				indexErr error
			)

			BeforeEach(func() {
				respRecorder = httptest.NewRecorder()

				req := httptest.NewRequest("PUT", "/scenario/scenario-1/device-1/no", nil)
				subject.NonProxyHandler(respRecorder, req)

				sendRequest = func(url string, ctx int64) {
					header := make(http.Header)
					context := &goproxy.ProxyCtx{Session: ctx}

					req, _ := MakeRequest("GET", url, header, nil)
					resp := MakeResponse2(200, header, "{\"code\" = 111}")

					subject.Request(req, context)
					subject.Response(resp, context)
				}

				sendRequest("http://test.net/users", 101)

				index, indexErr = storage.LoadIndex2(fs, "test/scenario-1/index.txt", true)
			})

			It("should activate scenario", func() {
				Expect(200).To(Equal(respRecorder.Code))
			})

			It("should create scenario folder", func() {
				dirExists, _ := afero.IsDir(fs, "test/scenario-1")
				Expect(dirExists).To(BeTrue())
			})

			It("should index contains expected records", func() {
				Expect(1).To(Equal(index.Size()))
				Expect(storage.IndexItem{Focused: true, Method: "GET", Url: "http://test.net/users", Code: 200, Folder: "r_1"}).To(Equal(index.Get(0)))
			})

			Describe("Activate scenario 2", func() {
				BeforeEach(func() {
					respRecorder = httptest.NewRecorder()

					req := httptest.NewRequest("PUT", "/scenario/scenario-2/device-1/no", nil)
					subject.NonProxyHandler(respRecorder, req)

					sendRequest("http://test.net/sessions", 102)
					sendRequest("http://test.net/events", 103)

					index, indexErr = storage.LoadIndex2(fs, "test/scenario-2/index.txt", true)
				})

				It("should load index", func() {
					Expect(indexErr).ToNot(HaveOccurred())
					Expect(index).ToNot(BeNil())
				})

				It("should index contains expected records", func() {
					Expect(2).To(Equal(index.Size()))
					Expect(storage.IndexItem{Focused: true, Method: "GET", Url: "http://test.net/sessions", Code: 200, Folder: "r_1"}).To(Equal(index.Get(0)))
					Expect(storage.IndexItem{Focused: true, Method: "GET", Url: "http://test.net/events", Code: 200, Folder: "r_2"}).To(Equal(index.Get(1)))
				})

				Describe("Activate scenario 1 again", func() {
					BeforeEach(func() {
						respRecorder = httptest.NewRecorder()

						req := httptest.NewRequest("PUT", "/scenario/scenario-1/device-1/no", nil)
						subject.NonProxyHandler(respRecorder, req)

						sendRequest("http://test.net/sessions", 105)
						sendRequest("http://test.net/events", 106)

						sendRequest("http://test.net/users", 107)
						sendRequest("http://test.net/sessions", 108)
						sendRequest("http://test.net/events", 109)

						sendRequest("http://test.net/info", 110)

						index, indexErr = storage.LoadIndex2(fs, "test/scenario-1/index.txt", true)
					})

					It("should load index", func() {
						Expect(indexErr).ToNot(HaveOccurred())
						Expect(index).ToNot(BeNil())
					})

					It("should index contains expected records", func() {
						Expect(4).To(Equal(index.Size()))
						Expect(storage.IndexItem{Focused: true, Method: "GET", Url: "http://test.net/users", Code: 200, Folder: "r_1"}).To(Equal(index.Get(0)))
						Expect(storage.IndexItem{Focused: true, Method: "GET", Url: "http://test.net/sessions", Code: 200, Folder: "r_2"}).To(Equal(index.Get(1)))
						Expect(storage.IndexItem{Focused: true, Method: "GET", Url: "http://test.net/events", Code: 200, Folder: "r_3"}).To(Equal(index.Get(2)))
						Expect(storage.IndexItem{Focused: true, Method: "GET", Url: "http://test.net/info", Code: 200, Folder: "r_4"}).To(Equal(index.Get(3)))
					})
				})
			})
		})
	})
})
