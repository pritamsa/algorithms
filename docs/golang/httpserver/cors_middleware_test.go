package httpserver

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/authorizeconstants"
	"gitlab.nordstrom.com/sentry/gohttp/middleware"
)

var _ = Describe("CorsMiddleware", func() {
	var (
		subject            middleware.Middleware
		handlerCalledCount int32
		handler            http.Handler
	)

	BeforeEach(func() {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&handlerCalledCount, 1)
		})

		subject = NewCorsMiddleware("")
	})

	It("provides a middleware wrapped with a cors handler", func() {
		wrappedHandler := subject.Wrap(handler)
		Expect(wrappedHandler).To(BeAssignableToTypeOf(corsHandler{}))
	})
})

var _ = Describe("CorsHandler", func() {
	var (
		subject      http.Handler
		innerHandler http.Handler
		req          *http.Request
		response     *httptest.ResponseRecorder
		err          error
		regex        string
	)

	BeforeEach(func() {
		regex = authorizeconstants.CorsRegex
		response = httptest.NewRecorder()
		innerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		httpTestServer := httptest.NewServer(subject)

		req, err = http.NewRequest("GET", httpTestServer.URL, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {

		subject = newCorsHandler(innerHandler, regex)
		subject.ServeHTTP(response, req)
	})

	Context("When there is no CORS regex provided to the service", func() {

		BeforeEach(func() {
			req.Header.Add("Origin", "https://awesome.nordstromrules.com")
			regex = ``
		})
		It("does not return Access-Control-Allow-Origin header", func() {
			Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
		})
	})

	Context("When there is a CORS regex provided to the service", func() {
		Context("When the request passes an Origin header doesn't end in nordstrom.com", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "https://awesome.nordstromrules.com")
			})

			It("does not return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			})
		})

		Context("When the request passes an Origin header that contains nordstrom.com, but doesn't end with it", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "https://awesome.nordstrom.com.badactor.com")
			})

			It("does not return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			})
		})

		Context("When the request passes an Origin header that matches the RegEx", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "https://reviews.nordstrom.com")
			})

			It("does return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			})
		})

		Context("When the request passes an Origin header with https://shop.nordstrom.com", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "https://shop.nordstrom.com")
			})

			It("does return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(Equal("https://shop.nordstrom.com"))
			})
		})

		Context("When the request passes an Origin header with http://shop.nordstrom.com", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "http://shop.nordstrom.com")
			})

			It("does return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			})
		})

		Context("When the request passes an Origin header with https://dev.shop.nordstrom.com", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "https://dev.shop.nordstrom.com")
			})

			It("does return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			})
		})

		Context("When the request passes an Origin header with http://dev.shop.nordstrom.com", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "http://dev.shop.nordstrom.com")
			})

			It("does return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			})
		})

		Context("When the request passes an Origin header with http://product-page.mwp.nordstrom.com", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "http://product-page.mwp.nordstrom.com")
			})

			It("does return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			})
		})

		Context("When the request passes an Origin header with http://studio2.products.dev.nordstrom.com:12345", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "http://studio2.products.dev.nordstrom.com:12345")
			})

			It("does return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			})
		})
	})
})

var _ = Describe("CorsOptionsHandler", func() {
	var (
		subject  http.Handler
		req      *http.Request
		response *httptest.ResponseRecorder
		err      error
		regex    string
	)

	BeforeEach(func() {
		regex = authorizeconstants.CorsRegex
		response = httptest.NewRecorder()
		httpTestServer := httptest.NewServer(subject)

		req, err = http.NewRequest("OPTIONS", httpTestServer.URL, nil)
		Expect(err).NotTo(HaveOccurred())

		req.Header.Add("Access-Control-Request-Method", "POST")
		req.Header.Add("Access-Control-Request-Headers", "Content-Type,Hello-World")
	})

	JustBeforeEach(func() {
		subject = NewCorsOptionsHandler(regex)
		subject.ServeHTTP(response, req)
	})

	Context("When there is no CORS regex provided to the service", func() {
		BeforeEach(func() {
			req.Header.Add("Origin", "awesome.nordstromrules.com")
		})
		It("does not return CORS header", func() {
			Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
			Expect(response.Header().Get("Access-Control-Allow-Methods")).To(BeEmpty())
			Expect(response.Header().Get("Access-Control-Allow-Headers")).To(BeEmpty())
		})
	})

	Context("When there is a CORS regex provided to the service", func() {
		Context("When the request passes an Origin header that does not match the RegEx", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "awesome.nordstromrules.com")
			})

			It("does not return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(BeEmpty())
				Expect(response.Header().Get("Access-Control-Allow-Methods")).To(BeEmpty())
				Expect(response.Header().Get("Access-Control-Allow-Headers")).To(BeEmpty())
			})
		})

		Context("When the request passes an Origin header that matches the RegEx", func() {
			BeforeEach(func() {
				req.Header.Add("Origin", "https://m.secure.nordstrom.com")
			})

			It("does return Access-Control-Allow-Origin header", func() {
				Expect(response.Header().Get("Access-Control-Allow-Origin")).To(Equal("https://m.secure.nordstrom.com"))
				Expect(response.Header().Get("Access-Control-Allow-Methods")).To(Equal("POST"))
				Expect(response.Header().Get("Access-Control-Allow-Headers")).To(Equal("Content-Type,Hello-World"))
			})
		})
	})
})
