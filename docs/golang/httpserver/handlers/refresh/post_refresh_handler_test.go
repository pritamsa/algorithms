package refresh_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/clients"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee/apigeefakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/clientsfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo/dynamofakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter/forterfakes"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/refresh"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper/statsd_wrapperfakes"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"
)

var _ = Describe("the refresh handler", func() {
	var tCtx *apmfakes.FakeTransactionContext

	prepareApm := func(r *http.Request) *http.Request {
		ctx := context.WithValue(r.Context(), "tCtx", tCtx)
		return r.WithContext(ctx)
	}
	Context("when a valid request is sent", func() {
		BeforeEach(func() {
			tCtx = &apmfakes.FakeTransactionContext{}
			tCtx.SegmentStub = func(s string) apm.Segment {
				return &apmfakes.FakeSegment{}
			}
			tCtx.NewGoRoutineReturns(tCtx)
		})

		It("is happy", func() {
			fakeGologger := &gologgerfakes.FakeLogger{}
			logging.CreateSingleLoggerForTest(fakeGologger, "", "")
			responseWriter := httptest.NewRecorder()
			apigeeFakeClient := &apigeefakes.FakeClient{}
			dynamoFakeClient := &dynamofakes.FakeClient{}
			statsdFakeClient := &statsd_wrapperfakes.FakeClient{}
			forterFakeClient := &forterfakes.FakeClient{}
			authEntity := model.AuthorizationEntity{
				PKCE:     "-S1l05-YI9a3yfaw5CcbxKedtiyPXkSwBBgCMzw14VQ*",
				AuthCode: "my_auth_code",
				PubKey:   "351A83D6-FFB0-48D9-A857-E1FE7DFC730C",
			}
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			body := `{ "token": "mytoken" }`
			subject := refresh.NewPostRefreshHandler(apigeeFakeClient, dynamoFakeClient, forterFakeClient, statsdFakeClient)
			closer := new(clientsfakes.FakeCloser)
			requestedPath := "/v1/refresh?code=my_auth_code&verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			request, _ := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request = prepareApm(request)
			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "myappiid")
			request.Header.Set("X-Akamai-Edgescape", "mygeoinfo")
			request.Header.Set("X-Nor-Scope", "webRegistered")
			subject.ServeHTTP(responseWriter, request)

			Expect(responseWriter.Code).To(Equal(http.StatusOK))
			resp := model.RefreshResponse{
				AccessToken:           "",
				RefreshToken:          "",
				ExpiresIn:             "0",
				RefreshTokenExpiresIn: "0",
			}
			json, _ := json.Marshal(resp)
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(string(json)))
		})

		It("responds normally when v1.1 refresh disable flag is true", func() {
			fakeGologger := &gologgerfakes.FakeLogger{}
			logging.CreateSingleLoggerForTest(fakeGologger, "", "")
			responseWriter := httptest.NewRecorder()
			apigeeFakeClient := &apigeefakes.FakeClient{}
			dynamoFakeClient := &dynamofakes.FakeClient{}
			dynamoFakeClient.GetFeatureFlagReturns(true)
			statsdFakeClient := &statsd_wrapperfakes.FakeClient{}
			forterFakeClient := &forterfakes.FakeClient{}
			authEntity := model.AuthorizationEntity{
				PKCE:     "-S1l05-YI9a3yfaw5CcbxKedtiyPXkSwBBgCMzw14VQ*",
				AuthCode: "my_auth_code",
				PubKey:   "351A83D6-FFB0-48D9-A857-E1FE7DFC730C",
			}
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}
			body := `{ "token": "mytoken" }`
			subject := refresh.NewPostRefreshHandler(apigeeFakeClient, dynamoFakeClient, forterFakeClient, statsdFakeClient)
			closer := new(clientsfakes.FakeCloser)
			requestedPath := "/v1/refresh?code=my_auth_code&verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			request, _ := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request = prepareApm(request)
			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "myappiid")
			request.Header.Set("X-Akamai-Edgescape", "mygeoinfo")
			subject.ServeHTTP(responseWriter, request)

			Expect(responseWriter.Code).To(Equal(http.StatusOK))
			resp := model.RefreshResponse{
				AccessToken:           "",
				RefreshToken:          "",
				ExpiresIn:             "0",
				RefreshTokenExpiresIn: "0",
			}
			json, _ := json.Marshal(resp)
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(string(json)))
		})

	})

	Context("v1.1", func() {
		var (
			dynamoFakeClient *dynamofakes.FakeClient
			forterFakeClient *forterfakes.FakeClient
			subject          http.Handler
			deviceRecord     model.DeviceRecord
			request          *http.Request
			responseWriter   *httptest.ResponseRecorder
		)
		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.SegmentStub = func(s string) apm.Segment {
			return &apmfakes.FakeSegment{}
		}
		tCtx.NewGoRoutineReturns(tCtx)
		fakeGologger := &gologgerfakes.FakeLogger{}
		logging.CreateSingleLoggerForTest(fakeGologger, "", "")
		apigeeFakeClient := &apigeefakes.FakeClient{}
		apigeeFakeClient.RefreshTokenReturns(apigee.Response{ShopperID: "YFedDA5om8F_HKo7SddADJrwOjDge0FqHarLHrZzyMf_pBhRZs_MmA2"}, nil)
		dynamoFakeClient = &dynamofakes.FakeClient{}
		statsdFakeClient := &statsd_wrapperfakes.FakeClient{}
		authEntity := model.AuthorizationEntity{
			PKCE:     "-S1l05-YI9a3yfaw5CcbxKedtiyPXkSwBBgCMzw14VQ*",
			AuthCode: "my_auth_code",
			PubKey:   "351A83D6-FFB0-48D9-A857-E1FE7DFC730C",
		}
		dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
			if t == dynamo.PkceTable {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			} else if t == dynamo.ShopperDeviceTable {
				data, _ := json.Marshal(deviceRecord)
				json.Unmarshal(data, o)
				return nil
			}
			return errors.New("Not Supported")
		}
		body := `{ "token": "mytoken" }`
		closer := new(clientsfakes.FakeCloser)
		requestedPath := "/v1.1/refresh?code=my_auth_code&verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"

		BeforeEach(func() {
			forterFakeClient = &forterfakes.FakeClient{}
			subject = refresh.NewPostRefreshHandler(apigeeFakeClient, dynamoFakeClient, forterFakeClient, statsdFakeClient)
			request, _ = http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request = prepareApm(request)
			responseWriter = httptest.NewRecorder()
		})

		It("is happy with v1.1", func() {
			// Arrange
			dynamoFakeClient.GetFeatureFlagReturns(false)
			forterFakeClient.GetForterAccountLoginResultReturns(forter.ForterAccountLoginResponse{ForterDecision: forter.Approve}, nil)
			deviceRecord = model.DeviceRecord{
				PersistentOptIn:           true,
				PersistentOptInExpiration: time.Now().AddDate(0, 0, 1).Unix(),
			}

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
			resp := model.RefreshResponseV11{}
			json, _ := json.Marshal(resp)
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(string(json)))
			clientModel := forterFakeClient.GetForterAccountLoginResultArgsForCall(0)
			Expect(clientModel.ShopperID).To(Equal("0ABCDEFABCDEFABCDEFABCDEFABCDEF0"))
			Expect(clientModel.MfaBypass).To(BeFalse())
			Expect(clientModel.LoginMethodType).To(Equal(forter.TokenRefresh))
		})

		It("returns a client error if flag is false", func() {
			// Arrange
			dynamoFakeClient.GetFeatureFlagReturns(false)
			deviceRecord = model.DeviceRecord{
				PersistentOptIn:           false,
				PersistentOptInExpiration: time.Now().AddDate(0, 0, 1).Unix(),
			}

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(http.StatusConflict))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})

		It("returns a client error if flag is expired", func() {
			// Arrange
			dynamoFakeClient.GetFeatureFlagReturns(false)
			deviceRecord = model.DeviceRecord{
				PersistentOptIn:           true,
				PersistentOptInExpiration: time.Now().AddDate(0, 0, -1).Unix(),
			}

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(http.StatusConflict))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})

		It("returns a client error if flag is false and expired", func() {
			// Arrange
			dynamoFakeClient.GetFeatureFlagReturns(false)
			deviceRecord = model.DeviceRecord{
				PersistentOptIn:           false,
				PersistentOptInExpiration: time.Now().AddDate(0, 0, -1).Unix(),
			}

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(http.StatusConflict))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})

		It("returns 250 with a Forter error", func() {
			// Arrange
			forterFakeClient.GetForterAccountLoginResultReturns(forter.ForterAccountLoginResponse{}, errors.New("Hi"))
			dynamoFakeClient.GetFeatureFlagReturns(false)
			deviceRecord = model.DeviceRecord{
				PersistentOptIn:           true,
				PersistentOptInExpiration: time.Now().AddDate(0, 0, 1).Unix(),
			}

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(250))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})

		It("returns challenged with a Forter recommendation", func() {
			// Arrange
			forterFakeClient.GetForterAccountLoginResultReturns(forter.ForterAccountLoginResponse{ForterDecision: forter.VerificationRequired}, nil)
			dynamoFakeClient.GetFeatureFlagReturns(false)
			deviceRecord = model.DeviceRecord{
				PersistentOptIn:           true,
				PersistentOptInExpiration: time.Now().AddDate(0, 0, 1).Unix(),
			}

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(250))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})

		It("returns 401 when the feature flag to disable v1.1 refresh is set to true", func() {
			// Arrange
			dynamoFakeClient.GetFeatureFlagReturns(true)

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(401))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})

		It("bypasses forter when bypass feature flag is enabled", func() {
			deviceRecord = model.DeviceRecord{
				PersistentOptIn:           true,
				PersistentOptInExpiration: time.Now().AddDate(0, 0, 1).Unix(),
			}
			dynamoFakeClient.GetFeatureFlagStub = func(id string, tCtx apm.TransactionContext) bool {
				return id == "bypassForterOnRefresh"
			}

			subject.ServeHTTP(responseWriter, request)

			Expect(responseWriter.Code).To(Equal(200))
			Expect(forterFakeClient.GetForterAccountLoginResultCallCount()).To(BeZero())
		})
	})

	Context("when apigee returns an error", func() {
		var (
			dynamoFakeClient *dynamofakes.FakeClient
			forterFakeClient *forterfakes.FakeClient
			subject          http.Handler
			deviceRecord     model.DeviceRecord
			request          *http.Request
			responseWriter   *httptest.ResponseRecorder
		)
		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.SegmentStub = func(s string) apm.Segment {
			return &apmfakes.FakeSegment{}
		}
		tCtx.NewGoRoutineReturns(tCtx)
		fakeGologger := &gologgerfakes.FakeLogger{}
		logging.CreateSingleLoggerForTest(fakeGologger, "", "")
		apigeeFakeClient := &apigeefakes.FakeClient{}
		dynamoFakeClient = &dynamofakes.FakeClient{}
		statsdFakeClient := &statsd_wrapperfakes.FakeClient{}
		authEntity := model.AuthorizationEntity{
			PKCE:     "-S1l05-YI9a3yfaw5CcbxKedtiyPXkSwBBgCMzw14VQ*",
			AuthCode: "my_auth_code",
			PubKey:   "351A83D6-FFB0-48D9-A857-E1FE7DFC730C",
		}
		dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
			if t == dynamo.PkceTable {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			} else if t == dynamo.ShopperDeviceTable {
				data, _ := json.Marshal(deviceRecord)
				json.Unmarshal(data, o)
				return nil
			}
			return errors.New("Not Supported")
		}
		body := `{ "token": "mytoken" }`
		closer := new(clientsfakes.FakeCloser)
		requestedPath := "/v1.1/refresh?code=my_auth_code&verifier=EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"

		BeforeEach(func() {
			forterFakeClient = &forterfakes.FakeClient{}
			subject = refresh.NewPostRefreshHandler(apigeeFakeClient, dynamoFakeClient, forterFakeClient, statsdFakeClient)
			request, _ = http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			request = prepareApm(request)
			responseWriter = httptest.NewRecorder()
		})

		It("responds with a 401 if apigee returns a Refresh Token Not Approved error", func() {
			// Arrange
			apigeeFakeClient.RefreshTokenReturns(apigee.Response{}, apigee.ErrRefreshTokenNotApproved)

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(401))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})

		It("responds with a 401 if apigee returns a Refresh Token Expired error", func() {
			// Arrange
			apigeeFakeClient.RefreshTokenReturns(apigee.Response{}, apigee.ErrRefreshTokenExpired)

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(401))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})

		It("responds with a 401 if apigee returns an Invalid Refresh Token error", func() {
			// Arrange
			apigeeFakeClient.RefreshTokenReturns(apigee.Response{}, apigee.ErrInvalidRefreshToken)

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(401))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})

		It("responds with a 401 if apigee returns an Access Token Not Approved error", func() {
			// Arrange
			apigeeFakeClient.RefreshTokenReturns(apigee.Response{}, apigee.ErrAccessTokenNotApproved)

			// Act
			subject.ServeHTTP(responseWriter, request)

			// Assert
			Expect(responseWriter.Code).To(Equal(401))
			exp := responseWriter.Body.String()
			Expect(exp).To(Equal(""))
		})
	})
})
