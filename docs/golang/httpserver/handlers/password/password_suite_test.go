package password

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAuthorize(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Forgot confirm password Suite")
}
