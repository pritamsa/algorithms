package authorize_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo/dynamofakes"
	authorizeHandler "gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/authorize"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/authorize/try"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper/statsd_wrapperfakes"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"
)

var _ = Describe("the authorize get handler", func() {
	var (
		subject http.Handler

		requestedPath  string
		responseWriter *httptest.ResponseRecorder

		fakeGologger            *gologgerfakes.FakeLogger
		dynamoFakeClient        *dynamofakes.FakeClient
		authEntity              model.AuthorizationEntity
		statsdClient            statsd_wrapper.Client
		statsdClientFake        *statsd_wrapperfakes.FakeGoStatsClient
		statsdClientFactoryFake *statsd_wrapperfakes.FakeStatsdClientFactory
		tCtx                    *apmfakes.FakeTransactionContext
	)

	prepareApm := func(r *http.Request) *http.Request {
		ctx := context.WithValue(r.Context(), "tCtx", tCtx)
		return r.WithContext(ctx)
	}

	BeforeEach(func() {
		statsdClientFake = &statsd_wrapperfakes.FakeGoStatsClient{}
		statsdClientFake.IncrStub = func(_ string, _ []string, _ float64) error {
			Expect(statsdClientFake.CloseCallCount()).To(Equal(0))
			return nil
		}

		statsdClientFactoryFake = &statsd_wrapperfakes.FakeStatsdClientFactory{}
		statsdClientFactoryFake.NewClientReturns(statsdClientFake, nil)
		statsdClient, _ = statsd_wrapper.NewClient(statsdClientFactoryFake)

		fakeGologger = &gologgerfakes.FakeLogger{}
		logging.CreateSingleLoggerForTest(fakeGologger, "", "")
		dynamoFakeClient = &dynamofakes.FakeClient{}
		authEntity = model.AuthorizationEntity{
			IPAddress:            "10.0.0.10",
			InstallationId:       "1234567890abcdefg-?",
			ClientId:             "NINTERNALIOS",
			PKCE:                 "-S1l05-YI9a3yfaw5CcbxKedtiyPXkSwBBgCMzw14VQ*",
			PKCEMethod:           "S256",
			AuthCode:             "my_auth_code",
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

		subject = authorizeHandler.NewGetAuthorizePageHandler(dynamoFakeClient, statsdClient, "int")
		responseWriter = httptest.NewRecorder()
		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.NewGoRoutineReturns(tCtx)
		tCtx.SegmentReturns(&apmfakes.FakeSegment{})
	})

	Context("when a valid request is sent", func() {

		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			requestedPath = "/v2/authorize?code=my_auth_code&verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"

			request, err := http.NewRequest("POST", requestedPath, nil)
			Expect(err).ToNot(HaveOccurred())

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(1))
			table, key, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(key["InstallationId"]).To(Equal("fooappiid"))
			Expect(key["IPAddress"]).To(Equal("1.2.3.*"))
		})

		It("calls setsubmitcode once", func() {
			Expect(dynamoFakeClient.UpdateAuthCallCount()).To(Equal(1))
			deviceId, trueClientIp, m, _ := dynamoFakeClient.UpdateAuthArgsForCall(0)
			Expect(deviceId).To(Equal("fooappiid"))
			Expect(trueClientIp).To(Equal("1.2.3.*"))
			Expect(m["SubmitCode"]).NotTo(BeEmpty())
			Expect(strings.HasSuffix(m["SubmitCode"], "==")).To(BeTrue())
			Expect(m["AuthGetTimestamp"]).NotTo(BeEmpty())
			t, err := time.Parse(time.RFC3339Nano, m["AuthGetTimestamp"])
			Expect(err).To(BeNil())
			Expect(t.Unix()).To(BeNumerically(">", 0))
		})

		It("sets the response code to 200", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})
	})

	Context("when setting submit code fails", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			dynamoFakeClient.UpdateAuthReturns(errors.New("hii"))

			requestedPath = "/v2/authorize?code=my_auth_code&verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"

			request, err := http.NewRequest("POST", requestedPath, nil)
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(1))
			table, key, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(key["InstallationId"]).To(Equal("fooappiid"))
			Expect(key["IPAddress"]).To(Equal("1.2.3.*"))
		})

		It("calls setsubmitcode once", func() {
			Expect(dynamoFakeClient.UpdateAuthCallCount()).To(Equal(1))
			deviceId, trueClientIp, m, _ := dynamoFakeClient.UpdateAuthArgsForCall(0)
			Expect(deviceId).To(Equal("fooappiid"))
			Expect(trueClientIp).To(Equal("1.2.3.*"))
			Expect(m["SubmitCode"]).NotTo(BeEmpty())
			Expect(strings.HasSuffix(m["SubmitCode"], "==")).To(BeTrue())
			Expect(m["AuthGetTimestamp"]).NotTo(BeEmpty())
			t, err := time.Parse(time.RFC3339Nano, m["AuthGetTimestamp"])
			Expect(err).To(BeNil())
			Expect(t.Unix()).To(BeNumerically(">", 0))
		})

		It("sets the response code to 500", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("when code is not present in query param", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			requestedPath = "/v2/authorize?verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"

			request, err := http.NewRequest("POST", requestedPath, nil)
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("sets the response code to 418", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusForbidden))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})
	})

	Context("when code is empty in query param", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			requestedPath = "/v2/authorize?code=&verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"

			request, err := http.NewRequest("POST", requestedPath, nil)
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("sets the response code to 418", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusForbidden))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})
	})

	Context("when auth code is not in dynamo", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetReturns(dynamo.ErrNotFound)
			requestedPath = "/v2/authorize?code=my_made_up_not_found_key&verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"

			request, err := http.NewRequest("GET", requestedPath, nil)
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically(">=", 1))
		})

		It("sets the response code to 418", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusForbidden))
		})

		It("calls logger", func() {
			try.Until(func() bool { return fakeGologger.InfoCallCount() > 0 })
			Expect(fakeGologger.InfoCallCount()).To(BeNumerically(">=", 1))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})
	})

	Context("when code in query does not match code in dynamo", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			requestedPath = "/v2/authorize?code=my_made_up_not_match_key&verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"

			request, err := http.NewRequest("GET", requestedPath, nil)
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(1))
		})

		It("sets the response code to 418", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusForbidden))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})
	})

	Context("when verifier is invalid", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			requestedPath = "/v2/authorize?code=my_auth_code&verifier=lalala"

			request, err := http.NewRequest("POST", requestedPath, nil)
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(1))
		})

		It("calls delete from dynamo zero times", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 418", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusForbidden))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})
	})
})
