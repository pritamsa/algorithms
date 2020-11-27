package authtoken_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAuthToken(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthToken Handler Suite")
}
