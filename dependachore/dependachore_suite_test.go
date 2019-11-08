package dependachore_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDependachore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dependachore Suite")
}
