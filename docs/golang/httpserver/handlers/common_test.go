package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"

	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo/dynamofakes"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper/statsd_wrapperfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/clients"
	"gitlab.nordstrom.com/sentry/authorize/clients/clientsfakes"
	. "gitlab.nordstrom.com/sentry/authorize/httpserver/handlers"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"
)

var _ = Describe("Common", func() {
	var (
		statsdClient            statsd_wrapper.Client
		statsdClientFake        *statsd_wrapperfakes.FakeGoStatsClient
		statsdClientFactoryFake *statsd_wrapperfakes.FakeStatsdClientFactory
		dynamoClient            *dynamofakes.FakeClient
		tCtx                    *apmfakes.FakeTransactionContext
	)

	prepareApm := func(r *http.Request) *http.Request {
		ctx := context.WithValue(r.Context(), "tCtx", tCtx)
		return r.WithContext(ctx)
	}

	BeforeEach(func() {
		statsdClientFactoryFake = &statsd_wrapperfakes.FakeStatsdClientFactory{}
		statsdClientFactoryFake.NewClientReturns(statsdClientFake, nil)
		statsdClient, _ = statsd_wrapper.NewClient(statsdClientFactoryFake)
		dynamoClient = &dynamofakes.FakeClient{}
		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.SegmentStub = func(s string) apm.Segment {
			return &apmfakes.FakeSegment{}
		}
		tCtx.NewGoRoutineReturns(tCtx)
	})
	Context("VerifyAuth", func() {
		It("returns empty error when request is valid", func() {
			authEntity := getAuthEntity()

			fakeGologger := &gologgerfakes.FakeLogger{}
			logging.CreateSingleLoggerForTest(fakeGologger, "", "")
			responseWriter := httptest.NewRecorder()
			dynamoClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			body := `
			{
				"code":"my_submit_code",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
			}`

			closer := new(clientsfakes.FakeCloser)
			request, _ := http.NewRequest("POST", "doesnot/matter", clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-ClientId", "NINTERNALIOS")
			request = prepareApm(request)
			_, err := VerifyAuth(responseWriter, request, statsdClient, dynamoClient)
			Expect(err).To(BeNil())
		})

		It("returns an error when code doesn't match", func() {
			authEntity := getAuthEntity()

			fakeGologger := &gologgerfakes.FakeLogger{}
			logging.CreateSingleLoggerForTest(fakeGologger, "", "")
			responseWriter := httptest.NewRecorder()
			dynamoClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			body := `
			{
				"code":"wrong_submit_code",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"irst_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
			}`

			closer := new(clientsfakes.FakeCloser)
			request, _ := http.NewRequest("POST", "doesnot/matter", clients.FakeReadCloser(body, closer))
			request = prepareApm(request)
			request.Header.Set("X-Nor-Appiid", "myappiid")
			request.Header.Set("True-Client-Ip", "1.1.1.1")
			request.Header.Set("X-Akamai-Edgescape", "mygeoinfo")
			request.Header.Set("X-Nor-ClientId", "NINTERNALIOS")
			_, err := VerifyAuth(responseWriter, request, statsdClient, dynamoClient)
			Expect(err).To(Equal(errors.New("Unauthorized")))
		})

		It("returns an error when code is empty", func() {

			fakeGologger := &gologgerfakes.FakeLogger{}
			logging.CreateSingleLoggerForTest(fakeGologger, "", "")
			responseWriter := httptest.NewRecorder()
			body := `
			{
				"code":"",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"irst_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
			}`

			closer := new(clientsfakes.FakeCloser)
			request, _ := http.NewRequest("POST", "doesnot/matter", clients.FakeReadCloser(body, closer))
			request = prepareApm(request)
			request.Header.Set("X-Nor-Appiid", "myappiid")
			request.Header.Set("True-Client-Ip", "1.1.1.1")
			request.Header.Set("X-Akamai-Edgescape", "mygeoinfo")
			request.Header.Set("X-Nor-ClientId", "NINTERNALIOS")
			_, err := VerifyAuth(responseWriter, request, statsdClient, dynamoClient)
			Expect(err).To(Equal(errors.New("Unauthorized")))
		})

		It("returns an error when verifier is empty", func() {

			fakeGologger := &gologgerfakes.FakeLogger{}
			logging.CreateSingleLoggerForTest(fakeGologger, "", "")
			responseWriter := httptest.NewRecorder()
			body := `
			{
				"code":"my_submit_code",
				"verifier":"",
				"irst_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
			}`

			closer := new(clientsfakes.FakeCloser)
			request, _ := http.NewRequest("POST", "doesnot/matter", clients.FakeReadCloser(body, closer))
			request = prepareApm(request)
			request.Header.Set("X-Nor-Appiid", "myappiid")
			request.Header.Set("True-Client-Ip", "1.1.1.1")
			request.Header.Set("X-Akamai-Edgescape", "mygeoinfo")
			request.Header.Set("X-Nor-ClientId", "NINTERNALIOS")
			_, err := VerifyAuth(responseWriter, request, statsdClient, dynamoClient)
			Expect(err).To(Equal(errors.New("Unauthorized")))
		})

		It("returns an error when the client IDs don't match", func() {

			fakeGologger := &gologgerfakes.FakeLogger{}
			logging.CreateSingleLoggerForTest(fakeGologger, "", "")
			responseWriter := httptest.NewRecorder()
			body := `
			{
				"code":"my_submit_code",
				"verifier":"",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
			}`

			closer := new(clientsfakes.FakeCloser)
			request, _ := http.NewRequest("POST", "doesnot/matter", clients.FakeReadCloser(body, closer))
			request = prepareApm(request)
			request.Header.Set("X-Nor-Appiid", "myappiid")
			request.Header.Set("True-Client-Ip", "1.1.1.1")
			request.Header.Set("X-Akamai-Edgescape", "mygeoinfo")
			request.Header.Set("X-Nor-ClientId", "incorrectcliendID")
			_, err := VerifyAuth(responseWriter, request, statsdClient, dynamoClient)
			Expect(err).To(Equal(errors.New("Unauthorized")))
		})
	})
})

func getAuthEntity() model.AuthorizationEntity {
	return model.AuthorizationEntity{
		IPAddress:            "10.0.0.10",
		InstallationId:       "1234567890abcdefg-?",
		ClientId:             "NINTERNALIOS",
		PKCE:                 "-S1l05-YI9a3yfaw5CcbxKedtiyPXkSwBBgCMzw14VQ*",
		PKCEMethod:           "S256",
		AuthCode:             "my_auth_code",
		SubmitCode:           "my_submit_code",
		PubKey:               "351A83D6-FFB0-48D9-A857-E1FE7DFC730C",
		RedirectURI:          "http://this.is.redire.ct/test?all=the%2e%things",
		AuthorizationTimeout: "1996-12-19T16:39:57-08:00",
		Scope:                "AUTHORIZED",
		AuthorizationAttempts: []model.AuthorizationAttempt{
			{
				IPAddress:       "10.0.0.9",
				InstallationId:  "1234567890abcdefg-?",
				AttemptDateTime: "1996-12-19T16:39:57-08:00",
				FailureReason:   "BADCLIENTID",
			},
			{
				IPAddress:       "10.0.0.9",
				InstallationId:  "1234567890abcdefg-?",
				AttemptDateTime: "1996-12-19T16:39:57-08:00",
				FailureReason:   "FAILEDPKCE",
			},
		},
	}
}
