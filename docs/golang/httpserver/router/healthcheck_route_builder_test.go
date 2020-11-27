package router_test

import (
	. "gitlab.nordstrom.com/sentry/authorize/httpserver/router"

	"net/http"

	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/gohealthcheck"
	"gitlab.nordstrom.com/sentry/gohttp/middleware/middlewarefakes"
)

var _ = Describe("HealthcheckRouteBuilder", func() {
	var (
		subject               HealthcheckRouteBuilder
		fakeMiddlewareWrapper *middlewarefakes.FakeMiddlewareWrapper

		routes gohealthcheck.Routes
	)

	BeforeEach(func() {
		fakeMiddlewareWrapper = new(middlewarefakes.FakeMiddlewareWrapper)
		fakeMiddlewareWrapper.AddMiddlewareToHandlerReturns(new(simpleHTTPHandler))

		subject = NewHealthcheckRouteBuilder(fakeMiddlewareWrapper)
	})

	JustBeforeEach(func() {
		routes = subject.RoutesForHealthchecks()
	})

	It("returns routes for healthchecks", func() {
		Expect(routes).To(BeAssignableToTypeOf(gohealthcheck.Routes{}))
	})

	It("should wrap the basic healthcheck with middleware", func() {
		Expect(routes.Basic.Handler).To(BeAssignableToTypeOf(new(simpleHTTPHandler)))

		responseRecorder := httptest.NewRecorder()
		firstWrappedHandler := fakeMiddlewareWrapper.AddMiddlewareToHandlerArgsForCall(0)
		firstWrappedHandler.ServeHTTP(responseRecorder, nil)

		Expect(responseRecorder.Body.String()).To(MatchJSON(`{"healthy":[],"unhealthy":[]}`))
	})

	It("should wrap the advanced healthcheck with middleware", func() {
		Expect(routes.Advanced.Handler).To(BeAssignableToTypeOf(new(simpleHTTPHandler)))

		responseRecorder := httptest.NewRecorder()
		secondWrappedHandler := fakeMiddlewareWrapper.AddMiddlewareToHandlerArgsForCall(1)
		secondWrappedHandler.ServeHTTP(responseRecorder, nil)

		Expect(responseRecorder.Body.String()).To(MatchJSON(`{"healthy":[],"unhealthy":[]}`))
	})
})

type simpleHTTPHandler struct{}

func (h *simpleHTTPHandler) ServeHTTP(_ http.ResponseWriter, _ *http.Request) {

}
