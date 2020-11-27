package challenge

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/clients"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	ioFakes "gitlab.nordstrom.com/sentry/authorize/clients/clientsfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth/shopperauthfakes"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"
)

var _ = Describe("The Challenge respond handler", func() {
	var (
		subject http.Handler

		responseWriter *httptest.ResponseRecorder

		fakeGologger *gologgerfakes.FakeLogger
		saFakeClient *shopperauthfakes.FakeClient
		closer       *ioFakes.FakeCloser
		body         string
		tCtx         *apmfakes.FakeTransactionContext
	)

	prepareApm := func(r *http.Request) *http.Request {
		ctx := context.WithValue(r.Context(), "tCtx", tCtx)
		return r.WithContext(ctx)
	}

	BeforeEach(func() {
		fakeGologger = &gologgerfakes.FakeLogger{}
		logging.CreateSingleLogger(fakeGologger, "", "")
		saFakeClient = &shopperauthfakes.FakeClient{}

		subject = NewChallengeRespondHandler(saFakeClient)
		responseWriter = httptest.NewRecorder()
		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.SegmentStub = func(s string) apm.Segment {
			return &apmfakes.FakeSegment{}
		}
		tCtx.NewGoRoutineReturns(tCtx)
	})

	Context("when a valid request is sent", func() {
		JustBeforeEach(func() {

			body = `
			{
				"sessionId":"ghghgh8988",
                "code":"990090"
			}`

			requestedPath := CHALLENGE_PATH

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set(CLIENT_IP_HEADER, DUMMY_CLIENT_IP)
			request.Header.Set(XNOR_APPIID_HEADER, DUMMY_APPIID)
			request.Header.Set(XNOR_CLIENT_ID, DUMMY_CLIENT_ID)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls customer auth challenge respond", func() {
			Expect(saFakeClient.ChallengeRespondCallCount()).To(BeNumerically("==", 1))
			model, request, _, _, _ := saFakeClient.ChallengeRespondArgsForCall(0)
			headers := request.Header
			Expect(model.SessionId).To(Equal("ghghgh8988"))
			Expect(model.Code).To(Equal("990090"))
			Expect(headers).ToNot(BeNil())
			Expect(headers.Get(XNOR_APPIID_HEADER)).To(Equal(DUMMY_APPIID))
			Expect(headers.Get(CLIENT_IP_HEADER)).To(Equal(DUMMY_CLIENT_IP))
		})

		It("writes the content type header for the response", func() {
			someHeaders := responseWriter.Header()
			Expect(someHeaders).NotTo(BeNil())
			Expect(someHeaders.Get("Content-Type")).NotTo(BeNil())
			Expect(someHeaders.Get("Content-Type")).To(Equal("application/json"))
		})

		It("returns a successful response", func() {
			responseCode := responseWriter.Code
			Expect(responseCode).To(Equal(200))
		})
	})

	Context("when an invalid request with unmarshallable body is sent ", func() {
		JustBeforeEach(func() {
			body = `
			{
				"sessionId":"ghghgh8988",
                "code":"9999%`

			code := DUMMY_AUTH
			verifier := DUMMY_VERIFIER
			requestedPath := CHALLENGE_PATH + code + "&verifier=" + verifier

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set(CLIENT_IP_HEADER, DUMMY_CLIENT_IP)
			request.Header.Set(XNOR_APPIID_HEADER, DUMMY_APPIID)
			request.Header.Set(XNOR_CLIENT_ID, DUMMY_CLIENT_ID)

			subject.ServeHTTP(responseWriter, request)
		})

		It("does not call customer auth challenge respond", func() {
			Expect(saFakeClient.ChallengeRespondCallCount()).To(Equal(0))
		})

		It("returns a bad request response", func() {
			responseCode := responseWriter.Code
			Expect(responseCode).To(Equal(400))
		})
	})

	Context("when a request missing code is sent ", func() {
		JustBeforeEach(func() {

			body = `
			{
				"sessionId":"",
                "code":""
			}`

			code := DUMMY_AUTH
			verifier := DUMMY_VERIFIER
			requestedPath := CHALLENGE_PATH + code + "&verifier=" + verifier

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set(CLIENT_IP_HEADER, DUMMY_CLIENT_IP)
			request.Header.Set(XNOR_APPIID_HEADER, DUMMY_APPIID)
			request.Header.Set(XNOR_CLIENT_ID, DUMMY_CLIENT_ID)

			subject.ServeHTTP(responseWriter, request)
		})

		It("does not call customer auth challenge init", func() {
			Expect(saFakeClient.ChallengeRespondCallCount()).To(Equal(0))
		})

		It("returns a bad request response", func() {
			responseCode := responseWriter.Code
			Expect(responseCode).To(Equal(400))
		})
	})

})
