package httpserver_test

import (
	"net/http"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitlab.nordstrom.com/sentry/authorize/httpserver"
	"gitlab.nordstrom.com/sentry/gohttp/middleware"
)

var _ = Describe("TimeoutMiddleware", func() {
	var (
		subject            middleware.Middleware
		handlerCalledCount int32
		handler            http.Handler
	)

	BeforeEach(func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&handlerCalledCount, 1)
		})

		subject = NewTimeoutMiddleware(25)
	})

	It("provides a middleware wrapped with a timeout handler", func() {
		request, _ := http.NewRequest("GET", "http://example.org", nil)

		wrappedHandler := subject.Wrap(handler)
		wrappedHandler.ServeHTTP(fakeResponseWriter{}, request)

		Eventually(func() int32 { return atomic.LoadInt32(&handlerCalledCount) }).Should(Equal(int32(1)))
		Expect(wrappedHandler).To(BeAssignableToTypeOf(http.TimeoutHandler(handler, 3, "message")))
	})
})

type fakeResponseWriter struct {
	http.ResponseWriter
}

func (fakeResponseWriter) Header() http.Header {
	return nil
}
func (fakeResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (fakeResponseWriter) WriteHeader(int) {}
