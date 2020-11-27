package signout_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper/statsd_wrapperfakes"

	"gitlab.nordstrom.com/sentry/authorize/clients"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee/apigeefakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/clientsfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter/forterfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth/shopperauthfakes"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/signout"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/try"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the sign out handler", func() {
	var (
		subject http.Handler

		requestedPath  string
		responseWriter *httptest.ResponseRecorder

		shopperAuthFakeClient *shopperauthfakes.FakeClient
		apigeeFakeClient      *apigeefakes.FakeClient
		forterFakeClient      *forterfakes.FakeClient
		statsdFakeClient      *statsd_wrapperfakes.FakeClient

		body   string
		closer *clientsfakes.FakeCloser
		tCtx   *apmfakes.FakeTransactionContext
	)

	prepareApm := func(r *http.Request) *http.Request {
		ctx := context.WithValue(r.Context(), "tCtx", tCtx)
		return r.WithContext(ctx)
	}

	BeforeEach(func() {
		logging.CreateSingleLoggerForTest(&gologgerfakes.FakeLogger{}, "", "")
		apigeeFakeClient = &apigeefakes.FakeClient{}
		forterFakeClient = &forterfakes.FakeClient{}
		statsdFakeClient = &statsd_wrapperfakes.FakeClient{}

		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.NewGoRoutineReturns(tCtx)
		tCtx.SegmentReturns(&apmfakes.FakeSegment{})
		subject = signout.NewPostSignOutHandler(shopperAuthFakeClient, apigeeFakeClient, forterFakeClient, statsdFakeClient)
		responseWriter = httptest.NewRecorder()
	})

	Context("when a valid request is sent with mobile scope", func() {
		JustBeforeEach(func() {
			requestedPath = "/v1/signout"
			body = `{"access_token": "test"}`

			closer = new(clientsfakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("User", "_8nwnEwkRCZU_6MtpR8WQPTUZaLxZdaHhQ3YBsmG9m55--9s58BPxQ2")
			request.Header.Set("X-Nor-Scope", "MobileRegistered")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dependencies", func() {
			try.Until(func() bool { return forterFakeClient.SendLogoutCallCount() > 0 })
			Expect(apigeeFakeClient.DeleteTokenCallCount()).To(Equal(1))
			Expect(forterFakeClient.SendLogoutCallCount()).To(Equal(1))
		})

		It("sets the response code to 200", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})
	})

	Context("when a valid request is sent with web scope", func() {
		JustBeforeEach(func() {
			requestedPath = "/v1/signout"
			body = `{"access_token": "test"}`

			closer = new(clientsfakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("User", "_8nwnEwkRCZU_6MtpR8WQPTUZaLxZdaHhQ3YBsmG9m55--9s58BPxQ2")
			request.Header.Set("X-Nor-Scope", "WebRegistered")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dependencies", func() {
			try.Until(func() bool { return forterFakeClient.SendLogoutCallCount() > 0 })
			Expect(apigeeFakeClient.DeleteTokenCallCount()).To(Equal(1))
			Expect(forterFakeClient.SendLogoutCallCount()).To(Equal(1))
		})

		It("sets the response code to 200", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})
	})

	Context("when missing User header", func() {
		JustBeforeEach(func() {
			requestedPath = "/v1/signout"
			body = `{"access_token": "test"}`

			closer = new(clientsfakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Scope", "MobileRegistered")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls apigee", func() {
			try.Until(func() bool { return apigeeFakeClient.DeleteTokenCallCount() > 0 })
			Expect(apigeeFakeClient.DeleteTokenCallCount()).To(Equal(1))
			Expect(forterFakeClient.SendLogoutCallCount()).To(Equal(0))
		})

		It("sets the response code to 200", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})
	})

	Context("when missing X-Nor-Scope header", func() {
		JustBeforeEach(func() {
			requestedPath = "/v1/signout"
			body = `{"access_token": "test"}`

			closer = new(clientsfakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("User", "_8nwnEwkRCZU_6MtpR8WQPTUZaLxZdaHhQ3YBsmG9m55--9s58BPxQ2")

			subject.ServeHTTP(responseWriter, request)
		})

		It("doesn't call dependencies", func() {
			try.Until(func() bool { return apigeeFakeClient.DeleteTokenCallCount() > 0 })
			Expect(apigeeFakeClient.DeleteTokenCallCount()).To(Equal(1))
			try.Until(func() bool { return forterFakeClient.SendLogoutCallCount() > 0 })
			Expect(forterFakeClient.SendLogoutCallCount()).To(Equal(1))
		})

		It("sets the response code to 200", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})
	})
})
