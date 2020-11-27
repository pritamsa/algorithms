package account

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/authorizeconstants"
	"gitlab.nordstrom.com/sentry/authorize/clients"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee"
	apigeeFake "gitlab.nordstrom.com/sentry/authorize/clients/apigee/apigeefakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	ioFakes "gitlab.nordstrom.com/sentry/authorize/clients/clientsfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/doomhammer"
	"gitlab.nordstrom.com/sentry/authorize/clients/doomhammer/doomhammerfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo/dynamofakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter/forterfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	shopperAuthFake "gitlab.nordstrom.com/sentry/authorize/clients/shopperauth/shopperauthfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/verify/verifyfakes"
	"gitlab.nordstrom.com/sentry/authorize/crypto"
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

var _ = Describe("the account post handler", func() {
	var (
		subject        http.Handler
		loginManager   login.Manager
		requestedPath  = "/v2/account"
		responseWriter *httptest.ResponseRecorder

		fakeGologger            *gologgerfakes.FakeLogger
		apigeeFakeClient        *apigeeFake.FakeClient
		saFakeClient            *shopperAuthFake.FakeClient
		doomHammerFakeClient    *doomhammerfakes.FakeClient
		verifyFakeClient        *verifyfakes.FakeClient
		forterFakeClient        *forterfakes.FakeClient
		mfaBypassFake           *mfafakes.FakeBypass
		tokenGenerator          *shoppertokenfakes.FakeTokenGenerator
		closer                  *ioFakes.FakeCloser
		statsdClient            statsd_wrapper.Client
		statsdClientFake        *statsd_wrapperfakes.FakeGoStatsClient
		statsdClientFactoryFake *statsd_wrapperfakes.FakeStatsdClientFactory
		dynamoFakeClient        *dynamofakes.FakeClient
		dynamoSvc               *dynamofakes.FakeDynamoService

		body            string
		authEntity      model.AuthorizationEntity
		fakeSAToken     shopperauth.TokenResponse
		fakeApigeeToken apigee.Response
		tCtx            *apmfakes.FakeTransactionContext
	)

	prepareApm := func(r *http.Request) *http.Request {
		ctx := context.WithValue(r.Context(), "tCtx", tCtx)
		return r.WithContext(ctx)
	}

	BeforeEach(func() {
		fakeGologger = &gologgerfakes.FakeLogger{}
		logging.CreateSingleLoggerForTest(fakeGologger, "", "")
		apigeeFakeClient = &apigeeFake.FakeClient{}
		saFakeClient = &shopperAuthFake.FakeClient{}
		doomHammerFakeClient = &doomhammerfakes.FakeClient{}
		verifyFakeClient = &verifyfakes.FakeClient{}
		forterFakeClient = &forterfakes.FakeClient{}
		statsdClientFake = &statsd_wrapperfakes.FakeGoStatsClient{}
		statsdClientFake.IncrStub = func(_ string, _ []string, _ float64) error {
			Expect(statsdClientFake.CloseCallCount()).To(Equal(0))
			return nil
		}
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

		statsdClientFactoryFake = &statsd_wrapperfakes.FakeStatsdClientFactory{}
		statsdClientFactoryFake.NewClientReturns(statsdClientFake, nil)
		statsdClient, _ = statsd_wrapper.NewClient(statsdClientFactoryFake)

		mfaBypassFake = &mfafakes.FakeBypass{}
		tokenGenerator = &shoppertokenfakes.FakeTokenGenerator{}

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
		fakeLogger := &gologgerfakes.FakeLogger{}
		encryptor := crypto.NewEncryptor([]byte(crypto.RandomString(32)), []byte(crypto.RandomString(16)), fakeLogger)

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
		subject = NewPostAccountHandler(
			doomHammerFakeClient,
			saFakeClient,
			forterFakeClient,
			statsdClient,
			loginManager,
			dynamoFakeClient,
			encryptor,
		)
		responseWriter = httptest.NewRecorder()
		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.SegmentReturns(&apmfakes.FakeSegment{})
		tCtx.NewGoRoutineReturns(tCtx)
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
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

		It("calls doomhammer for shopper creation", func() {
			Expect(doomHammerFakeClient.CreateShopperCallCount()).To(Equal(1))
		})

		It("calls forter to get account signup results", func() {
			Expect(forterFakeClient.GetForterAccountSignupResultCallCount()).To(Equal(1))
		})

		It("calls dynamo for Put", func() {
			try.Until(func() bool { return dynamoFakeClient.PutCallCount() == 2 })
			Expect(dynamoFakeClient.PutCallCount()).To(Equal(2))
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically(">=", 1))
			table, key, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(key["InstallationId"]).To(Equal("fooappiid"))
			Expect(key["IPAddress"]).To(Equal("1.2.3.*"))
		})

		It("calls shopperauth sign in once", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(1))
			username, password, headers, tc := saFakeClient.SignInArgsForCall(0)
			Expect(tc).NotTo(BeNil())
			Expect(headers).NotTo(BeNil())
			Expect(username).To(Equal("someemail@email.com"))
			Expect(password).To(Equal("ilikecookies"))
		})

		It("calls apigee token exchange once", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(1))
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

		Context("with no opt-in flag", func() {
			JustBeforeEach(func() {
				body = `
			{
				"code":"my_submit_code",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies"
			}`
				closer = new(ioFakes.FakeCloser)
				request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
				request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
				Expect(err).ToNot(HaveOccurred())
				request = prepareApm(request)

				request.Header.Set("True-Client-Ip", "1.2.3.4")
				request.Header.Set("X-Nor-Appiid", "fooappiid")

				subject.ServeHTTP(responseWriter, request)
			})
			It("sets the response code to 200", func() {
				Expect(responseWriter.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("when no body is posted", func() {
		JustBeforeEach(func() {
			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser("", closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})
		It("calls statsd zero times", func() {
			Expect(statsdClientFake.TimingCallCount()).To(Equal(0))
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

		It("sets the response code to 401", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("On request validation", func() {
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"irst_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
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

		It("with invalid fields returns 400", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Context("when code is not present in body", func() {
		JustBeforeEach(func() {
			body = `
			{
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo get once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically(">=", 1))
		})

		It("calls shopperauth sign in zero times", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(0))
		})

		It("calls apigee token exchange zero times", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
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
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true,
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

	Context("when shopperauth returns server error", func() {
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
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

		It("sets the response code to 500", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusInternalServerError))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("when shopperauth returns an error other than Server Error", func() {
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
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

		It("sets the response code to 500", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusInternalServerError))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})
	})

	Context("create account using DoomHammer", func() {
		TestErrCase := func(dhError error) {
			Context("when DoomHammer returns error "+dhError.Error(), func() {
				JustBeforeEach(func() {
					fakeLogger := &gologgerfakes.FakeLogger{}
					encryptor := crypto.NewEncryptor([]byte(crypto.RandomString(32)), []byte(crypto.RandomString(16)), fakeLogger)

					subject = NewPostAccountHandler(
						doomHammerFakeClient,
						saFakeClient,
						forterFakeClient,
						statsdClient,
						loginManager,
						dynamoFakeClient,
						encryptor,
					)
					doomHammerFakeClient.CreateShopperReturns(doomhammer.CreateShopperResponse{}, dhError)
					dynamoFakeClient.GetFeatureFlagReturns(true)

					dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
						data, _ := json.Marshal(authEntity)
						json.Unmarshal(data, o)
						return nil
					}

					requestedPath = "/v1/account"
					body = `
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

					closer = new(ioFakes.FakeCloser)
					request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
					request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

					Expect(err).ToNot(HaveOccurred())
					request = prepareApm(request)

					subject.ServeHTTP(responseWriter, request)
				})

				It("returns appropriate response", func() {
					body, err := ioutil.ReadAll(responseWriter.Body)
					Expect(err).To(BeNil())
					Expect(string(body)).To(Equal(`{"message":"` + dhError.Error() + `"}`))
				})

				It("calls dynamo get auth once", func() {
					Expect(dynamoFakeClient.GetCallCount()).To(Equal(1))
				})

				It("calls shopperauth sign in zero times", func() {
					Expect(saFakeClient.SignInCallCount()).To(Equal(0))
				})

				It("calls doomhammer once", func() {
					Expect(doomHammerFakeClient.CreateShopperCallCount()).To(Equal(1))
					_, _, _, s := doomHammerFakeClient.CreateShopperArgsForCall(0)
					matched, err := regexp.MatchString("^[A-Z0-9]{32}$", s)
					Expect(err).To(BeNil())
					Expect(matched).To(BeTrue())
				})

				It("calls apigee token exchange zero times", func() {
					Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
				})

				It("sets the response code to 400", func() {
					Expect(responseWriter.Code).To(Equal(http.StatusBadRequest))
				})

				It("closes the request body", func() {
					Expect(closer.CloseCallCount()).To(Equal(1))
				})
			})
		}
		TestErrCase(authorizeconstants.ErrEmailAlreadyExists)
		TestErrCase(authorizeconstants.ErrInvalidEmailFormat)
		TestErrCase(authorizeconstants.ErrMobileNumberAlreadyExists)
		TestErrCase(authorizeconstants.ErrInvalidMobileNumberFormat)
		TestErrCase(authorizeconstants.ErrInvalidFirstNameOrLastName)
		TestErrCase(authorizeconstants.ErrGenericBadRequest)

		Context("when DoomHammer returns an unknown error", func() {
			fakeLogger := &gologgerfakes.FakeLogger{}
			encryptor := crypto.NewEncryptor([]byte(crypto.RandomString(32)), []byte(crypto.RandomString(16)), fakeLogger)

			JustBeforeEach(func() {
				subject = NewPostAccountHandler(
					doomHammerFakeClient,
					saFakeClient,
					forterFakeClient,
					statsdClient,
					loginManager,
					dynamoFakeClient,
					encryptor,
				)
				doomHammerFakeClient.CreateShopperReturns(doomhammer.CreateShopperResponse{}, errors.New("Unknown"))
				dynamoFakeClient.GetFeatureFlagReturns(true)

				dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
					data, _ := json.Marshal(authEntity)
					json.Unmarshal(data, o)
					return nil
				}

				body = `
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

				closer = new(ioFakes.FakeCloser)
				request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
				request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

				Expect(err).ToNot(HaveOccurred())
				request = prepareApm(request)

				subject.ServeHTTP(responseWriter, request)
			})

			It("disables user account", func() {
				try.Until(func() bool { return saFakeClient.DisableUserCallCount() > 0 })
				Expect(saFakeClient.DisableUserCallCount()).To(Equal(1))
			})

			It("deletes user account", func() {
				try.Until(func() bool { return saFakeClient.DeleteUserCallCount() > 0 })
				Expect(saFakeClient.DeleteUserCallCount()).To(Equal(1))
			})

			It("returns appropriate response", func() {
				body, err := ioutil.ReadAll(responseWriter.Body)
				Expect(err).To(BeNil())
				Expect(string(body)).To(Equal(``))
			})

			It("calls dynamo get auth once", func() {
				Expect(dynamoFakeClient.GetCallCount()).To(Equal(1))
			})

			It("calls shopperauth sign in zero times", func() {
				Expect(saFakeClient.SignInCallCount()).To(Equal(0))
			})

			It("calls doomhammer once", func() {
				Expect(doomHammerFakeClient.CreateShopperCallCount()).To(Equal(1))
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

		Context("when DoomHammer does not return an error", func() {
			JustBeforeEach(func() {
				fakeLogger := &gologgerfakes.FakeLogger{}
				encryptor := crypto.NewEncryptor([]byte(crypto.RandomString(32)), []byte(crypto.RandomString(16)), fakeLogger)

				subject = NewPostAccountHandler(
					doomHammerFakeClient,
					saFakeClient,
					forterFakeClient,
					statsdClient,
					loginManager,
					dynamoFakeClient,
					encryptor,
				)
				doomHammerFakeClient.CreateShopperReturns(doomhammer.CreateShopperResponse{
					CustomerID: "877989880009",
					ShopperID:  "98880090909",
				},
					nil)
				dynamoFakeClient.GetFeatureFlagReturns(true)

				dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
					data, _ := json.Marshal(authEntity)
					json.Unmarshal(data, o)
					return nil
				}

				body = `
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

				closer = new(ioFakes.FakeCloser)
				request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
				request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

				Expect(err).ToNot(HaveOccurred())
				request = prepareApm(request)

				subject.ServeHTTP(responseWriter, request)
			})

			It("calls dynamo get auth", func() {
				Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically(">=", 2))
			})

			It("calls dynamo for Put", func() {
				try.Until(func() bool { return dynamoFakeClient.PutCallCount() == 2 })
				Expect(dynamoFakeClient.PutCallCount()).To(Equal(2))
			})

			It("Calls dynamo client for DeleteFromDynamo", func() {
				Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(1))
			})

			It("Calls create shopper to create a shopper", func() {
				Expect(doomHammerFakeClient.CreateShopperCallCount()).To(Equal(1))
			})

		})
	})

	Context("when DoomHammer finds error when creating shopper", func() {
		JustBeforeEach(func() {
			fakeLogger := &gologgerfakes.FakeLogger{}
			encryptor := crypto.NewEncryptor([]byte(crypto.RandomString(32)), []byte(crypto.RandomString(16)), fakeLogger)

			subject = NewPostAccountHandler(
				doomHammerFakeClient,
				saFakeClient,
				forterFakeClient,
				statsdClient,
				loginManager,
				dynamoFakeClient,
				encryptor,
			)
			doomHammerFakeClient.CreateShopperReturns(doomhammer.CreateShopperResponse{},
				authorizeconstants.ErrGenericBadRequest)
			dynamoFakeClient.GetFeatureFlagReturns(true)

			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			body = `
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

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("Does not call apigee client Exchange Token", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})
	})

	//New assertion with Doomhammer
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
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
		})

		It("calls shopperauth sign in once", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(1))
			username, password, headers, tc := saFakeClient.SignInArgsForCall(0)
			Expect(tc).NotTo(BeNil())
			Expect(headers).NotTo(BeNil())
			Expect(username).To(Equal("someemail@email.com"))
			Expect(password).To(Equal("ilikecookies"))
		})

		It("calls apigee token exchange once", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(1))
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

		Context("with no opt-in flag", func() {
			JustBeforeEach(func() {
				body = `
			{
				"code":"my_submit_code",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies"
			}`
				closer = new(ioFakes.FakeCloser)
				request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
				request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
				Expect(err).ToNot(HaveOccurred())
				request = prepareApm(request)

				request.Header.Set("True-Client-Ip", "1.2.3.4")
				request.Header.Set("X-Nor-Appiid", "fooappiid")

				subject.ServeHTTP(responseWriter, request)
			})
			It("sets the response code to 200", func() {
				Expect(responseWriter.Code).To(Equal(http.StatusOK))
			})
		})
	}) //

	Context("when apigee returns an error", func() {
		JustBeforeEach(func() {
			fakeSAToken.AccessToken = "MY_SA_TOKEN"
			fakeSAToken.ShopperID = "MY_SHOPPERID"
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
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

	Context("when a request with no first and last name is sent for X-Nor-Scope as WebRegistered", func() {
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
			request.Header.Set("X-Nor-Scope", "WebRegistered")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls doomhammer for shopper creation", func() {
			Expect(doomHammerFakeClient.CreateShopperCallCount()).To(Equal(1))
		})

		It("calls forter to get account signup results", func() {
			Expect(forterFakeClient.GetForterAccountSignupResultCallCount()).To(Equal(1))
		})

		It("calls dynamo for Put", func() {
			try.Until(func() bool { return dynamoFakeClient.PutCallCount() == 2 })
			Expect(dynamoFakeClient.PutCallCount()).To(Equal(2))
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically(">=", 1))
			table, key, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(key["InstallationId"]).To(Equal("fooappiid"))
			Expect(key["IPAddress"]).To(Equal("1.2.3.*"))
		})

		It("calls shopperauth sign in once", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(1))
			username, password, headers, tc := saFakeClient.SignInArgsForCall(0)
			Expect(tc).NotTo(BeNil())
			Expect(headers).NotTo(BeNil())
			Expect(username).To(Equal("someemail@email.com"))
			Expect(password).To(Equal("ilikecookies"))
		})

		It("calls apigee token exchange once", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(1))
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

		Context("with no opt-in flag", func() {
			JustBeforeEach(func() {
				body = `
			{
				"code":"my_submit_code",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies"
			}`
				closer = new(ioFakes.FakeCloser)
				request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
				request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
				Expect(err).ToNot(HaveOccurred())
				request = prepareApm(request)

				request.Header.Set("True-Client-Ip", "1.2.3.4")
				request.Header.Set("X-Nor-Appiid", "fooappiid")

				subject.ServeHTTP(responseWriter, request)
			})
			It("sets the response code to 200", func() {
				Expect(responseWriter.Code).To(Equal(http.StatusOK))
			})
		})
	})

	Context("when a request with no first and last name is sent for X-Nor-Scope as MobileRegistered", func() {
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
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies",
				"isOptIn": true
			}`

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
			request.Header.Set("X-Nor-Scope", "MobileRegistered")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls doomhammer for shopper creation", func() {
			Expect(doomHammerFakeClient.CreateShopperCallCount()).To(Equal(0))
		})

		It("calls forter to get account signup results", func() {
			Expect(forterFakeClient.GetForterAccountSignupResultCallCount()).To(Equal(0))
		})

		It("calls dynamo for Put", func() {
			try.Until(func() bool { return dynamoFakeClient.PutCallCount() == 2 })
			Expect(dynamoFakeClient.PutCallCount()).To(Equal(0))
		})

		It("calls dynamo get auth once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically(">=", 1))
			table, key, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(key["InstallationId"]).To(Equal("fooappiid"))
			Expect(key["IPAddress"]).To(Equal("1.2.3.*"))
		})

		It("calls shopperauth sign in once", func() {
			Expect(saFakeClient.SignInCallCount()).To(Equal(0))
		})

		It("calls apigee token exchange once", func() {
			Expect(apigeeFakeClient.ExchangeTokenCallCount()).To(Equal(0))
		})

		It("calls delete from dynamo once", func() {
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))
		})

		It("sets the response code to 200", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusBadRequest))
		})

		It("closes the request body", func() {
			Expect(closer.CloseCallCount()).To(Equal(1))
		})

		Context("with no opt-in flag", func() {
			JustBeforeEach(func() {
				body = `
			{
				"code":"my_submit_code",
				"verifier":"EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c",
				"first_name": "First",
				"last_name": "Last",
				"mobile_number": "206-555-1234",
				"email": "someemail@email.com",
				"password": "ilikecookies"
			}`
				closer = new(ioFakes.FakeCloser)
				request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
				request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")
				Expect(err).ToNot(HaveOccurred())
				request = prepareApm(request)

				request.Header.Set("True-Client-Ip", "1.2.3.4")
				request.Header.Set("X-Nor-Appiid", "fooappiid")

				subject.ServeHTTP(responseWriter, request)
			})
			It("sets the response code to 200", func() {
				Expect(responseWriter.Code).To(Equal(http.StatusBadRequest))
			})
		})
	})

})
