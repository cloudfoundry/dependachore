package tracker_test

import (
	"github.com/masters-of-cats/dependachore/tracker"
	"os"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getEnv(key string) string {
	val := os.Getenv(key)
	Expect(val).ToNot(BeEmpty())
	return val
}

func getEnvInt(key string) int {
	strval := getEnv(key)
	intval, err := strconv.Atoi(strval)
	Expect(err).NotTo(HaveOccurred())
	return intval
}

var _ = Describe("Tracker", func() {
	var (
		client *tracker.Client
	)

	BeforeEach(func() {
		apiKey := getEnv("API_KEY")
		projectID := getEnvInt("PROJECT_ID")
		client = tracker.NewClient(apiKey, projectID)
	})

	Describe("retrieving a story", func() {
		var story tracker.Story

		BeforeEach(func() {
			var err error
			story, err = client.CreateFeature("my name", "my description")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(client.Delete(story.ID)).To(Succeed())
		})

		It("can get a story by ID", func() {
			story, err := client.Get(story.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(story.Description).To(Equal("my description"))
		})
	})

	Describe("creating a story", func() {
		var story tracker.Story

		AfterEach(func() {
			Expect(client.Delete(story.ID)).To(Succeed())
		})

		It("can create a feature", func() {
			var err error

			story, err = client.CreateFeature("my name", "my description")
			Expect(err).NotTo(HaveOccurred())

			retrievedStory, err := client.Get(story.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(retrievedStory.Name).To(Equal(story.Name))
			Expect(retrievedStory.Description).To(Equal(story.Description))
		})
	})

	Describe("Move and Chorify", func() {
		var feature, release tracker.Story

		BeforeEach(func() {
			var err error
			release, err = client.CreateRelease("my release", "my release description")
			Expect(err).NotTo(HaveOccurred())
			feature, err = client.CreateFeature("my feature", "my feature description")
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(client.Delete(release.ID)).To(Succeed())
			Expect(client.Delete(feature.ID)).To(Succeed())
		})

		It("moves it after the release marker", func() {
			var err error

			Expect(client.MoveAndChorify(feature.ID, release.ID)).To(Succeed())

			feature, err = client.Get(feature.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(feature.StoryType).To(Equal("chore"))
			Expect(feature.AfterID).To(Equal(release.ID))

			release, err = client.Get(release.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(release.BeforeID).To(Equal(feature.ID))
		})
	})
})
