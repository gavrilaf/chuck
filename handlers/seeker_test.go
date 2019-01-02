package handlers_test

/*import (
	. "github.com/gavrilaf/chuck/handlers"
	. "github.com/gavrilaf/chuck/storage"
	. "github.com/gavrilaf/chuck/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Seeker", func() {
	var (
		log      Logger
		fs     *afero.Fs

		req *http.Request
		resp *http.Response
		err error

		subject  Seeker
	)

	BeforeEach(func() {
		log = NewLogger(cli.NewMockUi())
		fs = afero.NewMemMapFs()
	})

	Context("when open Seeker on the folder with index", func() {
		BeforeEach(func() {
			header := make(http.Header)
			req, _ = MakeRequest2("POST", "https://secure.api.com/login", header, "")
			resp = MakeResponse2(200, header, "{}")

			recorder, _ := NewRecorderWithFs(fs, "test", false, false, log)
			recorder.SetFocusedMode(true)

			recorder.RecordRequest(req, 1)
			recorder.RecordResponse(resp, 1)

			subject, err =  NewSeekerHandler
			NewSeekerWithFs(fs, "test")
		})

		It("should no error occured", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return Seeker object", func() {
			Expect(err).ToNot(BeNil())
		})

		// Add Count to the seeker - should contains one record

		Describe("handling focused request", func() {
			BeforeEach(func() {
				subject.
			})
		})
	})

	Context("when open Seeker on the empty folder]", func() {
		BeforeEach(func() {
			subject, err = NewSeekerWithFs(fs, "test-12")
		})

		It("should no error occured", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should return Seeker object", func() {
			Expect(err).ToNot(BeNil())
		})

		// Add Count to the seeker - should be empty
	})
})

*/
