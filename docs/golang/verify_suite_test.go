package verify_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVerify(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Verify Suite")
}
