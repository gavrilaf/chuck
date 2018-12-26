package storage_test

import (
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/gavrilaf/chuck/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bufio"
	"fmt"
	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"net/http"
)

var _ = Describe("Recorder", func() {
	var (
		log     Logger
		subject Recorder
		folder  string
		root    *afero.Afero

		createRequest  func() *http.Request
		createResponse func() *http.Response
	)

	BeforeEach(func() {
		createRequest = func() *http.Request {
			header := make(http.Header)
			header.Set("Content-Type", "application/json")
			req, _ := MakeRequest2("POST", "https://secure.api.com?query=123", header, "{}")
			return req
		}

		createResponse = func() *http.Response {
			body := `{"colors": []}`
			header := make(http.Header)
			header.Set("Content-Type", "application/json")
			header.Set("Content-Length", "6573")

			return MakeResponse2(200, header, body)
		}

		log = NewLogger(cli.NewMockUi())

		folder = "log-folder"

		fs := afero.NewMemMapFs()
		root = &afero.Afero{Fs: fs}
	})

	Describe("Create Recorder", func() {
		var (
			err         error
			dirExists   bool
			indexExists bool
		)

		Context("when createNewFolder is true", func() {
			BeforeEach(func() {
				subject, err = NewRecorderWithFs(root.Fs, folder, true, log)

				path := folder + "/" + subject.Name()
				dirExists, _ = root.DirExists(path)
				indexExists, _ = root.Exists(path + "/index.txt")
			})

			It("should not error occurred", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return Recorder object", func() {
				Expect(subject).ToNot(BeNil())
			})

			It("should create a recorder root logger folder", func() {
				Expect(dirExists).To(BeTrue())
			})

			It("should create an index file", func() {
				Expect(indexExists).To(BeTrue())
			})
		})

		Context("when createNewFolder is false", func() {
			BeforeEach(func() {
				subject, err = NewRecorderWithFs(root.Fs, folder, false, log)
				indexExists, _ = root.Exists(folder + "/index.txt")
			})

			It("should not error occurred", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create an index file", func() {
				Expect(indexExists).To(BeTrue())
			})
		})
	})

	Describe("Recording", func() {
		var (
			basePath string
			dumpPath string
			err      error
			reqId    int64
			respId   int64
			session  int64
			req      *http.Request
			resp     *http.Response
		)

		BeforeEach(func() {
			subject, _ = NewRecorderWithFs(root.Fs, "log", true, log)
			basePath = "log/" + subject.Name()
			session = 10
			req = createRequest()
			resp = createResponse()
		})

		It("should contains no pending requests", func() {
			Expect(subject.PendingCount()).To(Equal(0))
		})

		Describe("Record request", func() {
			BeforeEach(func() {
				reqId, err = subject.RecordRequest(req, session)
				dumpPath = fmt.Sprintf("%s/r_%d/", basePath, reqId)
			})

			It("should not error occurred", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should contain one pending request", func() {
				Expect(subject.PendingCount()).To(Equal(1))
			})

			It("should create request dump folder", func() {
				dirExists, _ := root.DirExists(dumpPath)
				Expect(dirExists).To(BeTrue())
			})

			It("should create request headers dump", func() {
				dumpExists, _ := root.Exists(dumpPath + "req_header.json")
				Expect(dumpExists).To(BeTrue())
			})

			It("should create request body dump", func() {
				dumpExists, _ := root.Exists(dumpPath + "req_body.json")
				Expect(dumpExists).To(BeTrue())
			})

			Describe("Record response", func() {
				BeforeEach(func() {
					respId, err = subject.RecordResponse(resp, session)
				})

				It("should not error occurred", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should contains no pending requests", func() {
					Expect(subject.PendingCount()).To(Equal(0))
				})

				It("should request id equal to response id", func() {
					Expect(reqId).To(Equal(respId))
				})

				It("should index.txt contains log record", func() {
					fi, _ := root.Open(basePath + "/" + "index.txt")
					scanner := bufio.NewScanner(fi)
					scanner.Scan()
					line := scanner.Text()

					expected := fmt.Sprintf("N,\t200,\tr_%d,\tPOST,\thttps://secure.api.com?query=123", reqId)
					Expect(expected).To(Equal(line))
				})

				It("should create response headers dump", func() {
					dumpExists, _ := root.Exists(dumpPath + "resp_header.json")
					Expect(dumpExists).To(BeTrue())
				})

				It("should create response body dump", func() {
					dumpExists, _ := root.Exists(dumpPath + "resp_body.json")
					Expect(dumpExists).To(BeTrue())
				})
			})

			Describe("Record focused", func() {
				BeforeEach(func() {
					subject.SetFocusedMode(true)

					reqId, _ = subject.RecordRequest(req, session)
					subject.RecordResponse(resp, session)
				})

				It("should record request as focused", func() {
					fi, _ := root.Open(basePath + "/" + "index.txt")
					scanner := bufio.NewScanner(fi)
					scanner.Scan()
					line := scanner.Text()

					expected := fmt.Sprintf("F,\t200,\tr_%d,\tPOST,\thttps://secure.api.com?query=123", reqId)
					Expect(expected).To(Equal(line))
				})
			})
		})
	})
})
