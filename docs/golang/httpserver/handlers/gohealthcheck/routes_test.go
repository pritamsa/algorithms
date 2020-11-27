package gohealthcheck_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/gohealthcheck"
	fakes "gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/gohealthcheck/gohealthcheckfakes"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"
	"gitlab.nordstrom.com/sentry/gologger/tracecontext"
)

var _ = Describe("Handlers", func() {
	var (
		checks         []Healthcheckable
		routes         Routes
		healthy        *fakes.FakeHealthcheckable
		unhealthy      *fakes.FakeHealthcheckable
		responseWriter *httptest.ResponseRecorder
		fakeLogger     *gologgerfakes.FakeLogger
		request        *http.Request
	)

	BeforeEach(func() {
		var err error
		responseWriter = httptest.NewRecorder()

		healthy = &fakes.FakeHealthcheckable{}
		healthy.NameReturns("i-am-alive")
		healthy.IsHealthyReturns(true, nil)

		unhealthy = &fakes.FakeHealthcheckable{}
		unhealthy.NameReturns("i-am-dead")
		unhealthy.IsHealthyReturns(false, errors.New("this-service-is-dead"))

		checks = []Healthcheckable{healthy, unhealthy}
		fakeLogger = &gologgerfakes.FakeLogger{}

		request, err = http.NewRequest(http.MethodGet, "", nil)
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		routes = RoutesForHealthchecksWithLogger(fakeLogger, checks...)
	})

	It("has two routes", func() {
		Expect(routes.Basic).To(BeAssignableToTypeOf(HealthcheckRoute{}))
		Expect(routes.Advanced).To(BeAssignableToTypeOf(HealthcheckRoute{}))
	})

	Describe("the basic route", func() {
		var basicRoute HealthcheckRoute

		JustBeforeEach(func() {
			basicRoute = routes.Basic
		})

		It("should have a good name", func() {
			Expect(basicRoute.Name).To(Equal("BasicHealthcheck"))
		})

		It("should have a suggested path for its handler", func() {
			Expect(basicRoute.Path).To(Equal("/status/basic"))
		})

		Describe("handles requests", func() {
			JustBeforeEach(func() {
				basicRoute.Handler.ServeHTTP(responseWriter, nil)
			})

			It("returns a 200 by default, with no checkers", func() {
				Expect(responseWriter.Code).To(Equal(http.StatusOK))
				Expect(responseWriter.Body).To(MatchJSON(`{"healthy":[], "unhealthy":[]}`))
			})

		})

		Describe("handles Interrupts", func() {
			JustBeforeEach(func() {
				if handlerShutdown, ok := basicRoute.Handler.(ShutdownHandler); ok {
					handlerShutdown.Shutdown()
				}
				basicRoute.Handler.ServeHTTP(responseWriter, nil)
			})

			It("returns 503", func() {
				Expect(responseWriter.Code).To(Equal(http.StatusServiceUnavailable))
			})
		})
	})

	Describe("the advanced route", func() {
		var advancedRoute HealthcheckRoute

		JustBeforeEach(func() {
			advancedRoute = routes.Advanced
		})

		It("should have a good name", func() {
			Expect(advancedRoute.Name).To(Equal("AdvancedHealthcheck"))
		})

		It("should have a suggested path", func() {
			Expect(advancedRoute.Path).To(Equal("/status/advanced"))
		})

		Describe("handles requests", func() {
			JustBeforeEach(func() {
				advancedRoute.Handler.ServeHTTP(responseWriter, request)
			})

			It("should use content-type application/json", func() {
				Expect(responseWriter.HeaderMap.Get("Content-Type")).To(Equal("application/json"))
			})

			It("filters the checks based on their health", func() {
				Expect(responseWriter.Body).To(MatchJSON(`{"healthy":["i-am-alive"], "unhealthy":["i-am-dead"]}`))
			})

			Context("when all of the checks are healthy", func() {
				BeforeEach(func() {
					checks = []Healthcheckable{healthy}
				})

				It("has a 200 status code", func() {
					Expect(responseWriter.Code).To(Equal(http.StatusOK))
				})
			})

			Context("when one of the checks is unhealthy", func() {
				var expectedTc tracecontext.TraceContext

				BeforeEach(func() {
					checks = []Healthcheckable{healthy, unhealthy}
					expectedTc = tracecontext.NewTraceContext()
					expectedTc.Set(request)
				})

				It("has a 500 status code", func() {
					Expect(responseWriter.Code).To(Equal(http.StatusInternalServerError))
				})

				It("should log the returned error from the unhealthy service", func() {
					Expect(fakeLogger.ErrorCallCount()).To(Equal(1))
					tc, message, err, stackTrace := fakeLogger.ErrorArgsForCall(0)
					Expect(tc).To(Equal(expectedTc))
					Expect(message).To(Equal("i-am-dead.health-check-failed"))
					Expect(err.Error()).To(Equal("this-service-is-dead"))
					Expect(stackTrace).To(BeEmpty())
				})
			})

		})
	})
})
