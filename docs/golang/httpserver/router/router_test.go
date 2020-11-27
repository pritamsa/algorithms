package router_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/authorizeconstants"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/gohealthcheck"
	. "gitlab.nordstrom.com/sentry/authorize/httpserver/router"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/router/routerfakes"
)

var _ = Describe("the authorize service router", func() {
	var subject http.Handler

	BeforeEach(func() {
		basicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("basic-healthcheck"))
		})

		advancedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("advanced-healthcheck"))
		})

		authorizeGetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		})

		authorizePostHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		})

		notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("There is nothing here."))
		})

		gohealthcheckRoutes := gohealthcheck.Routes{
			Basic: gohealthcheck.HealthcheckRoute{
				Path:    "/status/basic",
				Handler: basicHandler,
			},
			Advanced: gohealthcheck.HealthcheckRoute{
				Path:    "/status/advanced",
				Handler: advancedHandler,
			},
		}

		fakeHealthCheckRouteBuilder := new(routerfakes.FakeHealthcheckRouteBuilder)
		fakeHealthCheckRouteBuilder.RoutesForHealthchecksReturns(gohealthcheckRoutes)
		hcRoutes := fakeHealthCheckRouteBuilder.RoutesForHealthchecks()

		subject = NewRouter(notFoundHandler).
			AddRoute("GET", hcRoutes.Basic.Path, hcRoutes.Basic.Handler).
			AddRoute("GET", hcRoutes.Advanced.Path, hcRoutes.Advanced.Handler).
			SetPrefix(authorizeconstants.APIVersion1).
			AddRoute("GET", AuthorizePagePath, authorizeGetHandler).
			AddRoute("POST", AuthorizePostPath, authorizePostHandler).
			GetRoutes()
	})

	Describe("the router", func() {
		var responseWriter *httptest.ResponseRecorder

		BeforeEach(func() {
			responseWriter = httptest.NewRecorder()
		})

		It("routes to the provided basic healthcheck handler", func() {
			performRequestForPath("/status/basic", subject, responseWriter)
			Expect(responseWriter.Body.String()).To(ContainSubstring("basic-healthcheck"))
		})

		It("routes through to the provided advanced healthcheck handler", func() {
			performRequestForPath("/status/advanced", subject, responseWriter)
			Expect(responseWriter.Body.String()).To(ContainSubstring("advanced-healthcheck"))
		})

		It("uses 404 when routing to unmatched route", func() {
			performRequestForPath("/v1/lalalal", subject, responseWriter)
			Expect(responseWriter.Code).To(Equal(http.StatusNotFound))
			Expect(responseWriter.Body.String()).To(Equal("There is nothing here."))
		})

		It("routes to the provided http handler for the authorize get endpoint", func() {
			performRequestForPath("/v1/authorize", subject, responseWriter)
			Expect(responseWriter.Code).To(Equal(http.StatusTeapot))
		})

		It("routes to the provided http handler for the authorize post endpoint", func() {
			performPostForPath("/v1/authorize", subject, responseWriter)
			Expect(responseWriter.Code).To(Equal(http.StatusTeapot))
		})

		It("uses 404 when routing to unmatched route", func() {
			performRequestForPath("/v1/authorize/lala", subject, responseWriter)
			Expect(responseWriter.Code).To(Equal(http.StatusNotFound))
			Expect(responseWriter.Body.String()).To(Equal("There is nothing here."))
		})

		It("uses the 404 handler for un-matched routes", func() {
			performRequestForPath("/v2/whoops", subject, responseWriter)
			Expect(responseWriter.Code).To(Equal(http.StatusNotFound))
			Expect(responseWriter.Body.String()).To(Equal("There is nothing here."))
		})
	})
})

func performRequestForPath(requestedPath string, handler http.Handler, responseWriter http.ResponseWriter) {
	request, err := http.NewRequest("GET", requestedPath, nil)
	Expect(err).ToNot(HaveOccurred())

	handler.ServeHTTP(responseWriter, request)
}

func performPostForPath(requestedPath string, handler http.Handler, responseWriter http.ResponseWriter) {
	request, err := http.NewRequest("POST", requestedPath, nil)
	Expect(err).ToNot(HaveOccurred())

	handler.ServeHTTP(responseWriter, request)
}
