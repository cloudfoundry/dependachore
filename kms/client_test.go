package kms_test

import (
	"github.com/masters-of-cats/dependachore/kms"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kms.Client", func() {

	const (
		project  = "cf-garden-core"
		location = "global"
		keyRing  = "garden"
	)

	var (
		randomBytes []byte
		kmsClient   *kms.Client
		key         string
	)

	BeforeEach(func() {
		key = "test"
		randomBytes = []byte{}
		for i := 0; i < 128; i++ {
			randomBytes = append(randomBytes, byte(rand.Intn(256)))
		}
	})

	JustBeforeEach(func() {
		kmsClient = kms.NewClient(project, location, keyRing, key)
	})

	It("encrypts and decrypts secrets using the provided crypto key", func() {
		cipherText, err := kmsClient.Encrypt(randomBytes)
		Expect(err).NotTo(HaveOccurred())
		Expect(cipherText).ToNot(Equal(randomBytes))

		plainText, err := kmsClient.Decrypt(cipherText)
		Expect(err).NotTo(HaveOccurred())
		Expect(plainText).To(Equal(randomBytes))
	})

	When("we don't have access to the key", func() {
		BeforeEach(func() {
			key = "dependachore"
		})

		It("errors", func() {
			_, err := kmsClient.Encrypt(randomBytes)
			Expect(err).To(HaveOccurred())
		})
	})
})
