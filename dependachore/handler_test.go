package dependachore_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/masters-of-cats/dependachore/dependachore"
	"github.com/masters-of-cats/dependachore/dependachore/dependachorefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const releaseMarkerID = 42

var _ = Describe("Handler", func() {
	var (
		trackerClient *dependachorefakes.FakeTrackerClient
		handler       dependachore.Handler
		body          string
		method        string
		req           *http.Request
		resp          *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		trackerClient = new(dependachorefakes.FakeTrackerClient)
		handler = dependachore.NewHandler(trackerClient, releaseMarkerID)
		resp = httptest.NewRecorder()
		method = http.MethodPost
	})

	JustBeforeEach(func() {
		var err error
		req, err = http.NewRequest(method, "", strings.NewReader(body))
		Expect(err).NotTo(HaveOccurred())

		handler.Handle(resp, req)
	})

	When("the request is not a POST", func() {
		BeforeEach(func() {
			method = http.MethodGet
		})

		It("fails", func() {
			Expect(resp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	When("the request body is invalid json", func() {
		BeforeEach(func() {
			body = `{"foo"`
		})

		It("fails", func() {
			Expect(resp.Code).To(Equal(http.StatusBadRequest))
		})
	})

	When("the story is not from dependabot", func() {
		BeforeEach(func() {
			bodyBytes, err := ioutil.ReadFile("assets/not-dependabot-example.json")
			Expect(err).NotTo(HaveOccurred())
			body = string(bodyBytes)
		})

		It("does not call tracker client", func() {
			Expect(trackerClient.MoveAndChorifyCallCount()).To(Equal(0))
			Expect(resp.Code).To(Equal(http.StatusOK))
		})
	})

	When("the activity is not a creation", func() {
		BeforeEach(func() {
			bodyBytes, err := ioutil.ReadFile("assets/not-create-example.json")
			Expect(err).NotTo(HaveOccurred())
			body = string(bodyBytes)
		})

		It("does not call tracker client", func() {
			Expect(trackerClient.MoveAndChorifyCallCount()).To(Equal(0))
			Expect(resp.Code).To(Equal(http.StatusOK))
		})
	})

	When("the story looks like dependabot creation", func() {
		BeforeEach(func() {
			bodyBytes, err := ioutil.ReadFile("assets/full-example.json")
			Expect(err).NotTo(HaveOccurred())
			body = string(bodyBytes)
		})

		It("moves the story under the release marker", func() {
			Expect(trackerClient.MoveAndChorifyCallCount()).To(Equal(1))
			actualStoryID, actualReleaseMarkerID := trackerClient.MoveAndChorifyArgsForCall(0)
			Expect(actualStoryID).To(Equal(169649127))
			Expect(actualReleaseMarkerID).To(Equal(releaseMarkerID))
		})
	})
})
