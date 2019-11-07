package tracker_test

import (
	"dependachore/tracker"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tracker", func() {
	var (
		storyID int
		client  *tracker.Client
	)

	BeforeEach(func() {
		storyID = 169623658
		apiKey := os.Getenv("API_KEY")
		projectID := 1158420
		client = tracker.NewClient(apiKey, projectID)
	})

	Describe("retrieving a story", func() {
		It("can get a story by ID", func() {
			story, err := client.Get(storyID)
			Expect(err).NotTo(HaveOccurred())
			Expect(story.Description).To(ContainSubstring("garden-shed"))
		})
	})

	Describe("Chorify", func() {
		It("can transform story into chore", func() {
			err := client.Chorify(storyID)
			Expect(err).NotTo(HaveOccurred())

			story, err := client.Get(storyID)
			Expect(err).NotTo(HaveOccurred())
			Expect(story.StoryType).To(Equal("chore"))
		})
	})

	Describe("Move it after", func() {
		releaseID := 169626389

		It("moves it after the release marker", func() {
			Expect(client.MoveAfter(storyID, releaseID)).To(Succeed())
			release, err := client.Get(releaseID)
			Expect(err).NotTo(HaveOccurred())

			Expect(release.BeforeID).To(Equal(storyID))

			story, err := client.Get(storyID)
			Expect(err).NotTo(HaveOccurred())
			Expect(story.AfterID).To(Equal(releaseID))
		})

	})

})
