package passwordv2

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAuthorize(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Forgot confirm password v2 Suite")
}
