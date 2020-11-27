package authorize_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/clients"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee"
	apigeeFake "gitlab.nordstrom.com/sentry/authorize/clients/apigee/apigeefakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	ioFakes "gitlab.nordstrom.com/sentry/authorize/clients/clientsfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo/dynamofakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter/forterfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	shopperAuthFake "gitlab.nordstrom.com/sentry/authorize/clients/shopperauth/shopperauthfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/verify/verifyfakes"
	"gitlab.nordstrom.com/sentry/authorize/crypto"
	authorizeHandler "gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/authorize"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/login"
	"gitlab.nordstrom.com/sentry/authorize/mfa/mfafakes"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/authorize/shoppertoken/shoppertokenfakes"
	"gitlab.nordstrom.com/sentry/authorize/try"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper/statsd_wrapperfakes"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"
)

const (
	seattleLat  = 47.6062
	seattleLong = -122.3321
	bitsize     = 64
)

var _ = Describe("the authorize post handler", func() {
	var (
		subject http.Handler

		requestedPath  = "/v2/authorize"
		responseWriter *httptest.ResponseRecorder

		encryptor               *crypto.Encryptor
		fakeGologger            *gologgerfakes.FakeLogger
		fakeLogger              *gologgerfakes.FakeLogger
		apigeeFakeClient        *apigeeFake.FakeClient
		saFakeClient            *shopperAuthFake.FakeClient
		verifyFakeClient        *verifyfakes.FakeClient
		dynamoFakeClient        *dynamofakes.FakeClient
		dynamoSvc               *dynamofakes.FakeDynamoService
		forterFakeClient        *forterfakes.FakeClient
		mfaBypassFake           *mfafakes.FakeBypass
		tokenGenerator          *shoppertokenfakes.FakeTokenGenerator
		loginManager            login.Manager
		closer                  *ioFakes.FakeCloser
		statsdClient            statsd_wrapper.Client
		statsdClientFake        *statsd_wrapperfakes.FakeGoStatsClient
		statsdClientFactoryFake *statsd_wrapperfakes.FakeStatsdClientFactory
		body                    string
		authEntity              model.AuthorizationEntity
		fakeSAToken             shopperauth.TokenResponse
		fakeApigeeToken         apigee.Response
		tCtx                    *apmfakes.FakeTransactionContext
	)

	prepareApm := func(r *http.Request) *http.Request {
		ctx := context.WithValue(r.Context(), "tCtx", tCtx)
		return r.WithContext(ctx)
	}

	BeforeEach(func() {
		fakeGologger = &gologgerfakes.FakeLogger{}
		logging.CreateSingleLogger(fakeGologger, "", "")
		apigeeFakeClient = &apigeeFake.FakeClient{}
		saFakeClient = &shopperAuthFake.FakeClient{}
		verifyFakeClient = &verifyfakes.FakeClient{}
		forterFakeClient = &forterfakes.FakeClient{}
		mfaBypassFake = &mfafakes.FakeBypass{}
		tokenGenerator = &shoppertokenfakes.FakeTokenGenerator{}
		statsdClientFake = &statsd_wrapperfakes.FakeGoStatsClient{}
		statsdClientFake.IncrStub = func(_ string, _ []string, _ float64) error {
			Expect(statsdClientFake.CloseCallCount()).To(Equal(0))
			return nil
		}

		statsdClientFactoryFake = &statsd_wrapperfakes.FakeStatsdClientFactory{}
		statsdClientFactoryFake.NewClientReturns(statsdClientFake, nil)

		statsdClient, _ = statsd_wrapper.NewClient(statsdClientFactoryFake)

		authEntity = model.AuthorizationEntity{
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
			AuthGetTimestamp:     time.Now().Format(time.RFC3339Nano),
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
		fakeLogger = &gologgerfakes.FakeLogger{}
		encryptor = crypto.NewEncryptor([]byte(crypto.RandomString(32)), []byte(crypto.RandomString(16)), fakeLogger)
		dynamoFakeClient = &dynamofakes.FakeClient{}
		dynamoSvc = &dynamofakes.FakeDynamoService{}
		dynamoSvc.GetItemReturnsOnCall(0, &dynamodb.GetItemOutput{Item: map[string]*dynamodb.AttributeValue{
			"WebShopperID": {S: aws.String("123")},
			"DeviceID":     {S: aws.String("myappiid")},
			"Appiid":       {S: aws.String("myappiid")},
			"Lat":          {S: aws.String(strconv.FormatFloat(seattleLat, 'f', -1, bitsize))},
			"Long":         {S: aws.String(strconv.FormatFloat(seattleLong, 'f', -1, bitsize))},
			"Time":         {S: aws.String(time.Now().Add(time.Hour * time.Duration(8)).String())},
		}}, nil)

		loginManager = login.NewManager(
			verifyFakeClient,
			apigeeFakeClient,
			forterFakeClient,
			dynamoFakeClient,
			saFakeClient,
			mfaBypassFake,
			encryptor,
			tokenGenerator,
			statsdClient,
		)

		subject = authorizeHandler.NewPostAuthorizeHandler(statsdClient, dynamoFakeClient, loginManager, encryptor)
		responseWriter = httptest.NewRecorder()
		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.SegmentStub = func(s string) apm.Segment {
			return &apmfakes.FakeSegment{}
		}
		tCtx.NewGoRoutineReturns(tCtx)
	})

	Context("when no body is posted", func() {
		JustBeforeEach(func() {
			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser("", closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth zero times", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(0))
		})

		It("calls shopperauth sign in zero times", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(0))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo zero times", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 401", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when a valid request is sent", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			fakeSAToken.AccessToken = "MY_SA_TOKEN"
			fakeSAToken.ShopperID = "MY_SHOPPERID"
			fakeSAToken.WebShopperID = "MY_WEBSHOPPERID"
			saFakeClient.SignInReturns(fakeSAToken, nil)

			fakeApigeeToken.AccessToken = "MY_APIGEE_TOKEN"
			fakeApigeeToken.ShopperID = "MY_SHOPPERID"
			apigeeFakeClient.ExchangeTokenReturns(fakeApigeeToken, nil)

			body = `
			{
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically(">=", 1))

			table, key, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(key["InstallationId"]).To(Equal("fooappiid"))
			Expect(key["IPAddress"]).To(Equal("1.2.3.*"))

			table, key, _, _ = dynamoFakeClient.GetArgsForCall(1)
			Expect(table).To(Equal(dynamo.ShopperIdTable))
			Expect(key["Email"]).To(Equal("testusername@test.com"))
		})

		It("calls shopperauth sign in once", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(1))
			username, password, headers, tc := saFakeClient.SignInArgsForCall(0)
			Expect(tc).NotTo(BeNil())
			Expect(headers).NotTo(BeNil())
			Expect(username).To(Equal("testusername@test.com"))
			Expect(password).To(Equal("testpass"))
		})

		It("calls apigee token exchange once", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(1))
		})

		It("calls put from dynamo", func() {
			Expect(dynamoFakeClient.PutWithTTLCallCount()).To(Equal(1))
			table, e, ttl, _ := dynamoFakeClient.PutWithTTLArgsForCall(0)
			Expect(table).To(Equal(dynamo.ShopperIdTable))
			Expect(ttl).To(Equal(2592000))
			entity := e.(login.ShopperIDCache)
			Expect(entity.Email).To(Equal("testusername@test.com"))
			Expect(entity.ShopperID).To(Equal("MY_WEBSHOPPERID"))
		})

		It("calls delete from dynamo once", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(1))
			deviceId, trueClientIp, _ := dynamoFakeClient.DeleteAuthArgsForCall(0)
			Expect(deviceId).To(Equal("fooappiid"))
			Expect(trueClientIp).To(Equal("1.2.3.*"))
		})

		It("sets the response code to 200", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when code is not present in body", func() {
		JustBeforeEach(func() {
			body = `
			{
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth zero times", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(0))
		})

		It("calls shopperauth sign in zero times", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(0))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo once", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 401", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when code is empty in body", func() {
		JustBeforeEach(func() {
			body = `
			{
				"code":"",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth zero times", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(0))
		})

		It("calls shopperauth sign in zero times", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(0))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo once", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 401", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when auth code is not in dynamo", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetReturns(dynamo.ErrNotFound)
			body = `
			{
				"code":"my_made_up_not_found_key",
				"username__my_made_up_not_found_key":"testusername@test.com",
				"password__my_made_up_not_found_key":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically(">=", 1))
		})

		It("calls shopperauth sign in zero times", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(0))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo zero times", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 401", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when auth code is not in dynamo initially but found later", func() {
		JustBeforeEach(func() {
			// dynamoFakeClient.Get is called for NosaPkce and NosaShopperId tables.
			// In this test we only want to mock the NosaPkce behavior of not found initally and found later
			firstCall := true
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				if t == dynamo.ShopperIdTable {
					data, _ := json.Marshal(authEntity)
					json.Unmarshal(data, o)
					return nil
				}

				// t == dynamo.PkceTable
				if firstCall {
					firstCall = false
					noSubmitCodeEntity := authEntity
					noSubmitCodeEntity.SubmitCode = ""
					data, _ := json.Marshal(noSubmitCodeEntity)
					json.Unmarshal(data, o)
					return nil
				} else {
					data, _ := json.Marshal(authEntity)
					json.Unmarshal(data, o)
					return nil
				}
			}
			body = `
			{
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("returns appropriate code and calls dependencies", func() {
			// One call to ShopperIdTable
			// Two calls to PkceTable (1 retry)
			// One call to ShopperDeviceTable
			try.Until(func() bool { return dynamoFakeClient.GetCallCount() == 4 })
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(4))
			Expect(saFakeClient.SignInCallCount()).To(Equal(1))
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(1))
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(1))
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})
	})

	Context("when code in body doesn't match code in dynamo", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			body = `
			{
				"code":"my_made_up_not_found_key",
				"username__my_made_up_not_found_key":"testusername@test.com",
				"password__my_made_up_not_found_key":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(1))
		})

		It("calls shopperauth sign in zero times", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(0))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo zero times", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 401", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when verifier is invalid", func() {
		JustBeforeEach(func() {
			saFakeClient.SignInReturns(shopperauth.TokenResponse{}, shopperauth.ErrServerErr)

			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			body = `
			{
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"lalala"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(1))
		})

		It("calls shopperauth sign in zero times", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(0))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo zero times", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 401", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})

		It("sets the redirect url to authorize", func() {
			Expect(responseWriter.Header().Get("Location")).To(Equal("/v1/authinit"))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when shopperauth returns Server Error", func() {
		JustBeforeEach(func() {
			saFakeClient.SignInReturns(shopperauth.TokenResponse{}, shopperauth.ErrServerErr)

			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			body = `
			{
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			request = prepareApm(request)

			Expect(err).ToNot(HaveOccurred())

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(2))
			table, _, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			table, _, _, _ = dynamoFakeClient.GetArgsForCall(1)
			Expect(table).To(Equal(dynamo.ShopperIdTable))
		})

		It("calls shopperauth sign in once", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(1))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo once", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 500", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusInternalServerError))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when shopperauth returns Unauthorized", func() {
		JustBeforeEach(func() {
			saFakeClient.SignInReturns(shopperauth.TokenResponse{}, shopperauth.ErrUnauthorized)

			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			body = `
			{
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(2))
		})

		It("calls shopperauth sign in once", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(1))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo once", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 401", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when shopperauth returns an error other than Server Error or Unauthorized", func() {
		JustBeforeEach(func() {
			saFakeClient.SignInReturns(shopperauth.TokenResponse{}, errors.New("testing"))

			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			body = `
			{
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(2))
		})

		It("calls shopperauth sign in once", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(1))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo zero times", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 500", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusInternalServerError))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when apigee returns an error", func() {
		JustBeforeEach(func() {
			fakeSAToken.AccessToken = "MY_SA_TOKEN"
			fakeSAToken.ShopperID = "MY_SHOPPERID"
			fakeSAToken.WebShopperID = "MY_WEBSHOPPERID"
			saFakeClient.SignInReturns(fakeSAToken, nil)

			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			apigeeFakeClient.ExchangeTokenReturns(apigee.Response{}, errors.New("testing"))

			body = `
			{
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically(">=", 1))
		})

		It("calls shopperauth sign in once", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(1))
		})

		It("calls apigee token exchange once", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(1))
		})

		It("calls put from dynamo", func() {
			table, e, ttl, _ := dynamoFakeClient.PutWithTTLArgsForCall(0)
			Expect(table).To(Equal(dynamo.ShopperIdTable))
			Expect(ttl).To(Equal(2592000))
			entity := e.(login.ShopperIDCache)
			Expect(entity.Email).To(Equal("testusername@test.com"))
			Expect(entity.ShopperID).To(Equal("MY_WEBSHOPPERID"))
		})

		It("calls delete from dynamo once", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(1))
		})

		It("sets the response code to 500", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusInternalServerError))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("persistent sign-in", func() {
		It("persists true when passed true", func() {
			// Arrange
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			saFakeClient.SignInReturns(fakeSAToken, nil)
			apigeeFakeClient.ExchangeTokenReturns(fakeApigeeToken, nil)

			body = `
			{
				"persistentOptIn":true,
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			count := 1
			try.Until(func() bool { return dynamoFakeClient.PutCallCount() == count })
			Expect(dynamoFakeClient.PutCallCount()).To(Equal(count))
			var table string
			var e interface{}
			for i := 0; i < count && table != dynamo.ShopperDeviceTable; i++ {
				table, e, _ = dynamoFakeClient.PutArgsForCall(i)
			}
			Expect(table).To(Equal(dynamo.ShopperDeviceTable))
			entity1 := e.(model.DeviceRecord)
			Expect(entity1.PersistentOptIn).To(BeTrue())
		})

		It("persists false when passed false", func() {
			// Arrange
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			saFakeClient.SignInReturns(fakeSAToken, nil)
			apigeeFakeClient.ExchangeTokenReturns(fakeApigeeToken, nil)

			body = `
			{
				"persistentOptIn":false,
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			count := 1
			try.Until(func() bool { return dynamoFakeClient.PutCallCount() == count })
			Expect(dynamoFakeClient.PutCallCount()).To(Equal(count))
			var table string
			var e interface{}
			for i := 0; i < count && table != dynamo.ShopperDeviceTable; i++ {
				table, e, _ = dynamoFakeClient.PutArgsForCall(i)
			}
			Expect(table).To(Equal(dynamo.ShopperDeviceTable))
			entity1 := e.(model.DeviceRecord)
			Expect(entity1.PersistentOptIn).To(BeFalse())
		})

		It("persists false when not passed", func() {
			// Arrange
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			saFakeClient.SignInReturns(fakeSAToken, nil)
			apigeeFakeClient.ExchangeTokenReturns(fakeApigeeToken, nil)

			body = `
			{
				"code":"my_submit_code",
				"username__my_submit_code":"testusername@test.com",
				"password__my_submit_code":"testpass",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			count := 1
			try.Until(func() bool { return dynamoFakeClient.PutCallCount() == count })
			Expect(dynamoFakeClient.PutCallCount()).To(Equal(count))
			var table string
			var e interface{}
			for i := 0; i < count && table != dynamo.ShopperDeviceTable; i++ {
				table, e, _ = dynamoFakeClient.PutArgsForCall(i)
			}
			Expect(table).To(Equal(dynamo.ShopperDeviceTable))
			entity1 := e.(model.DeviceRecord)
			Expect(entity1.PersistentOptIn).To(BeFalse())
		})
	})
})
