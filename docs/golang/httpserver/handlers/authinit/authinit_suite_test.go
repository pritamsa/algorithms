package authinit_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAuthinit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Authinit Suite")
}
