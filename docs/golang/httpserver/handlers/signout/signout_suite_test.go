package signout_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSignout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Signout Suite")
}
