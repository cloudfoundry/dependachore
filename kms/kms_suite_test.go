package kms_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKms(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kms Suite")
}
