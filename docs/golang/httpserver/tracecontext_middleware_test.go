package httpserver

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/tracecontext"
	"gitlab.nordstrom.com/sentry/gohttp/middleware"
)

var _ = Describe("TraceContextMiddleware", func() {
	var (
		subject            middleware.Middleware
		handlerCalledCount int32
		handler            http.Handler
	)

	BeforeEach(func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&handlerCalledCount, 1)
		})

		subject = NewTraceContextMiddleware()
	})

	It("provides a middleware wrapped with a tracecontext handler", func() {
		wrappedHandler := subject.Wrap(handler)
		Expect(wrappedHandler).To(BeAssignableToTypeOf(traceContextHandler{}), "message")
	})
})

var _ = Describe("TraceContextHandler", func() {
	var (
		subject      http.Handler
		innerHandler http.Handler
		req          *http.Request
		err          error
	)

	BeforeEach(func() {
		innerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//noop
		})
		subject = NewTraceContextHandler(innerHandler)
		httpTestServer := httptest.NewServer(subject)

		req, err = http.NewRequest("GET", httpTestServer.URL, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		subject.ServeHTTP(nil, req)
	})

	Context("When the request does not include a trace context", func() {
		It("creates a tracecontext and adds it to the request", func() {
			Expect(req.Header.Get(tracecontext.TRACE_CONTEXT_HEADER)).NotTo(BeEmpty())
		})
	})

	Context("When the request does include a trace context", func() {
		BeforeEach(func() {
			req.Header.Add(tracecontext.TRACE_CONTEXT_HEADER, "TraceContextValue")
		})
		It("does not change the existing trace context", func() {
			Expect(req.Header.Get(tracecontext.TRACE_CONTEXT_HEADER)).To(Equal("TraceContextValue"))
		})
	})
})
