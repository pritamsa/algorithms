package authorize_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAuthorize(t *testing.T) {
	os.Setenv("CMDPATH", "gitlab.nordstrom.com/sentry/authorize")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Authorize Suite")
}
