package httpserver_test

import (
	"net/http"
	"net/http/httptest"

	"gitlab.nordstrom.com/sentry/authorize/httpserver"
	"gitlab.nordstrom.com/sentry/authorize/logging/loggingfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("a 404 handler", func() {
	var (
		subject http.Handler
		logger  *loggingfakes.FakeLogger

		requestedPath  string
		responseWriter *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		logger = new(loggingfakes.FakeLogger)
		subject = httpserver.NewNotFoundHandler(logger)

		responseWriter = httptest.NewRecorder()
	})

	JustBeforeEach(func() {
		requestedPath = "/v2/whoops/this/is/wrong"
		request, err := http.NewRequest("GET", requestedPath, nil)
		Expect(err).ToNot(HaveOccurred())

		subject.ServeHTTP(responseWriter, request)
	})

	It("logs the requested path", func() {
		Expect(logger.InfoCallCount()).To(Equal(1))

		_, message, data := logger.InfoArgsForCall(0)
		Expect(message).To(Equal("404_not_found"))
		Expect(data["requested_path"]).To(Equal(requestedPath))
	})

	It("sets the response code to 404", func() {
		Expect(responseWriter.Code).To(Equal(http.StatusNotFound))
	})
})
