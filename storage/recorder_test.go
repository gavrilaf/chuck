package storage_test

import (
	. "chuck/storage"
	. "chuck/utils"

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
		root    *afero.Afero

		createRequest  func(string) *http.Request
		createResponse func() *http.Response
	)

	BeforeEach(func() {
		createRequest = func(method string) *http.Request {
			header := make(http.Header)
			header.Set("Content-Type", "application/json")
			req, _ := MakeRequest2(method, "https://secure.api.com?query=123", header, "{}")
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
				subject, err = NewRecorder(root.Fs, log, "log-1", true, false, true, false)

				path := "log-1/" + subject.Name()
				dirExists, _ = root.DirExists(path)
				indexExists, _ = root.Exists(path + "/" + IndexFileName)
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
				subject, err = NewRecorder(root.Fs, log, "log-2", false, false, true, false)
				indexExists, _ = root.Exists("log-2/" + IndexFileName)
			})

			It("should not error occurred", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should create an index file", func() {
				Expect(indexExists).To(BeTrue())
			})
		})
	})

	Describe("Recording with (onlyNew=false, logRequest=true, appyFilters=false)", func() {
		var (
			basePath string
			dumpPath string
			session  int64
			req      *http.Request
			resp     *http.Response

			err        error
			reqResult  *PendingRequest
			respResult *PendingRequest
		)

		BeforeEach(func() {
			subject, _ = NewRecorder(root.Fs, log, "log-3", true, false, true, false)
			basePath = "log-3/" + subject.Name()
			session = 10
			req = createRequest("POST")
			resp = createResponse()
		})

		It("should contains no pending requests", func() {
			Expect(subject.PendingCount()).To(Equal(0))
		})

		Describe("Record request", func() {
			BeforeEach(func() {
				reqResult, err = subject.RecordRequest(req, session)
				dumpPath = fmt.Sprintf("%s/r_%d/", basePath, reqResult.Id)
			})

			It("should not error occurred", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should return pending request", func() {
				Expect(reqResult).ToNot(BeNil())
			})

			It("should create dump folder", func() {
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
					respResult, err = subject.RecordResponse(resp, session)
				})

				It("should not error occurred", func() {
					Expect(err).ToNot(HaveOccurred())
				})

				It("should index contains log record", func() {
					fi, _ := root.Open(basePath + "/" + IndexFileName)
					defer fi.Close()

					scanner := bufio.NewScanner(fi)
					scanner.Scan()
					line := scanner.Text()

					expected := fmt.Sprintf("N,\t200,\tr_%d,\tPOST,\thttps://secure.api.com?query=123", respResult.Id)
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

				Describe("Create new_only recorder based on the same dir", func() {
					var (
						subjectNew Recorder
					)

					BeforeEach(func() {
						subjectNew, err = NewRecorder(root.Fs, log, basePath, false, true, true, false)
					})

					It("should not error occurred", func() {
						Expect(err).ToNot(HaveOccurred())
					})

					It("should create recorder", func() {
						Expect(subjectNew).ToNot(BeNil())
					})

					Context("when record the request with the same method and url", func() {
						BeforeEach(func() {
							reqResult, err = subjectNew.RecordRequest(req, 101)
						})

						It("should not error occurred", func() {
							Expect(err).ToNot(HaveOccurred())
						})

						It("should not return pending request", func() {
							Expect(reqResult).To(BeNil())
						})
					})

					Context("when record the request with the new method", func() {
						BeforeEach(func() {
							req = createRequest("PUT")
							reqResult, err = subjectNew.RecordRequest(req, 102)
						})

						It("should not error occurred", func() {
							Expect(err).ToNot(HaveOccurred())
						})

						It("should return the pending request", func() {
							Expect(reqResult).ToNot(BeNil())
						})

						Describe("record the response", func() {
							BeforeEach(func() {
								respResult, err = subjectNew.RecordResponse(resp, 102)
							})

							It("should not error occurred", func() {
								Expect(err).ToNot(HaveOccurred())
							})

							It("should return the same result", func() {
								Expect(reqResult).To(Equal(respResult))
							})
						})
					})
				})
			})

			Describe("Record focused", func() {
				BeforeEach(func() {
					subject.SetFocusedMode(true)

					reqResult, _ = subject.RecordRequest(req, session)
					subject.RecordResponse(resp, session)
				})

				It("should record request as focused", func() {
					fi, _ := root.Open(basePath + "/" + IndexFileName)
					defer fi.Close()
					scanner := bufio.NewScanner(fi)
					scanner.Scan()
					line := scanner.Text()

					expected := fmt.Sprintf("F,\t200,\tr_%d,\tPOST,\thttps://secure.api.com?query=123", reqResult.Id)
					Expect(expected).To(Equal(line))
				})
			})

			Describe("Record request/response with empty header & body", func() {
				BeforeEach(func() {
					header := make(http.Header)
					req, _ := MakeRequest("GET", "www.google.com", header, nil)
					resp := MakeResponse(200, header, nil, 0)

					reqResult, _ = subject.RecordRequest(req, session)
					subject.RecordResponse(resp, session)

					dumpPath = fmt.Sprintf("%s/r_%d/", basePath, reqResult.Id)
				})

				It("should record request as usual", func() {
					fi, _ := root.Open(basePath + "/" + IndexFileName)
					defer fi.Close()
					scanner := bufio.NewScanner(fi)
					scanner.Scan()
					line := scanner.Text()

					expected := fmt.Sprintf("N,\t200,\tr_%d,\tGET,\twww.google.com", reqResult.Id)
					Expect(expected).To(Equal(line))
				})

				It("should create dump folder", func() {
					dirExists, _ := root.DirExists(dumpPath)
					Expect(dirExists).To(BeTrue())
				})

				It("should not create request headers dump", func() {
					dumpExists, _ := root.Exists(dumpPath + "req_header.json")
					Expect(dumpExists).ToNot(BeTrue())
				})

				It("should not create request body dump", func() {
					dumpExists, _ := root.Exists(dumpPath + "req_body.json")
					Expect(dumpExists).ToNot(BeTrue())
				})

				It("should not create response headers dump", func() {
					dumpExists, _ := root.Exists(dumpPath + "resp_header.json")
					Expect(dumpExists).ToNot(BeTrue())
				})

				It("should not create response body dump", func() {
					dumpExists, _ := root.Exists(dumpPath + "resp_body.json")
					Expect(dumpExists).ToNot(BeTrue())
				})
			})
		})

		Describe("with (onlyNew=true, logRequest=false, appyFilters=true)", func() {
			var (
				basePath string
			)

			BeforeEach(func() {
				subject, err = NewRecorder(root.Fs, log, "log-4", true, true, false, true)
				basePath = "log-4/" + subject.Name()

				req := createRequest("POST")
				resp := createResponse()

				subject.RecordRequest(req, 10)
				subject.RecordResponse(resp, 10)

				subject.RecordRequest(req, 11)
				subject.RecordResponse(resp, 11)

				req = createRequest("GET")

				subject.RecordRequest(req, 12)
				subject.RecordResponse(resp, 12)

				req = createRequest("PUT")
				resp.StatusCode = 404

				subject.RecordRequest(req, 13)
				subject.RecordResponse(resp, 13)
			})

			It("should should contains only two successful requests", func() {
				fi, _ := root.Open(basePath + "/" + IndexFileName)
				defer fi.Close()

				scanner := bufio.NewScanner(fi)

				Expect(scanner.Scan()).To(BeTrue())
				Expect(scanner.Text()).To(Equal(fmt.Sprintf("N,\t200,\tr_1,\tPOST,\thttps://secure.api.com?query=123")))

				Expect(scanner.Scan()).To(BeTrue())
				Expect(scanner.Text()).To(Equal(fmt.Sprintf("N,\t200,\tr_2,\tGET,\thttps://secure.api.com?query=123")))

				Expect(scanner.Scan()).To(BeFalse(), "should contain only two records")
			})

			It("should create folders only for recorded requests", func() {
				Expect(root.Exists(basePath + "/r_1")).To(BeTrue())
				Expect(root.Exists(basePath + "/r_2")).To(BeTrue())

				Expect(root.Exists(basePath + "/r_3")).ToNot(BeTrue())
			})

			It("should record only response header & body", func() {
				Expect(root.Exists(basePath + "/r_1/req_header.json")).ToNot(BeTrue())
				Expect(root.Exists(basePath + "/r_1/req_body.json")).ToNot(BeTrue())

				Expect(root.Exists(basePath + "/r_1/resp_header.json")).To(BeTrue())
				Expect(root.Exists(basePath + "/r_1/resp_body.json")).To(BeTrue())
			})
		})
	})
})
