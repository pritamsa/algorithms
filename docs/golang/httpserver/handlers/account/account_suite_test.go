package account_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAccount(t *testing.T) {
	os.Setenv("CMDPATH", "gitlab.nordstrom.com/sentry/account")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Account Suite")
}
