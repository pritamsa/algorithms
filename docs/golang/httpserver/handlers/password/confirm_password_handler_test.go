package password

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"

	"gitlab.nordstrom.com/sentry/authorize/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/clients"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	ioFakes "gitlab.nordstrom.com/sentry/authorize/clients/clientsfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo/dynamofakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter/forterfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth/shopperauthfakes"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/login/loginfakes"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"
)

var _ = Describe("The confirm password handler", func() {
	var (
		subject http.Handler

		responseWriter *httptest.ResponseRecorder

		fakeGologger     *gologgerfakes.FakeLogger
		saFakeClient     *shopperauthfakes.FakeClient
		dynamoFakeClient *dynamofakes.FakeClient
		dynamoSvc        *dynamofakes.FakeDynamoService
		forterFakeClient *forterfakes.FakeClient
		loginFakeManager *loginfakes.FakeManager
		closer           *ioFakes.FakeCloser
		body             string
		authEntity       model.AuthorizationEntity
		tCtx             *apmfakes.FakeTransactionContext
	)

	prepareApm := func(r *http.Request) *http.Request {
		ctx := context.WithValue(r.Context(), "tCtx", tCtx)
		return r.WithContext(ctx)
	}

	BeforeEach(func() {
		fakeGologger = &gologgerfakes.FakeLogger{}
		logging.CreateSingleLogger(fakeGologger, "", "")
		saFakeClient = &shopperauthfakes.FakeClient{}
		forterFakeClient = &forterfakes.FakeClient{}
		loginFakeManager = &loginfakes.FakeManager{}

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

		subject = NewConfirmPasswordHandler(saFakeClient, dynamoFakeClient, forterFakeClient, loginFakeManager)
		responseWriter = httptest.NewRecorder()
		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.SegmentStub = func(s string) apm.Segment {
			return &apmfakes.FakeSegment{}
		}
		tCtx.NewGoRoutineReturns(tCtx)
	})

	Context("when a valid confirm password request is sent", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			body = `
			{
               "confirmationCode":"898985",
				"email":"test@email.com",
               "shopperId":"90909067677",
               "password":"password#4"
    
			}`

			code := "my_auth_code"
			verifier := "EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			requestedPath := "/v1/password/confirm?code=" + code + "&verifier=" + verifier

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			dynamoSvc.DeleteItemReturns(&dynamodb.DeleteItemOutput{}, nil) //result, error
			dynamoFakeClient.GetFeatureFlagReturns(false)

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically("==", 1))
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(BeNumerically("==", 1))

			table, key, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(key["InstallationId"]).To(Equal("fooappiid"))
			Expect(key["IPAddress"]).To(Equal("1.2.3.*"))

			appiid, clientip, _ := dynamoFakeClient.DeleteAuthArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(appiid).To(Equal("fooappiid"))
			Expect(clientip).To(Equal("1.2.3.4"))
		})

		It("calls customer auth confirm password", func() {
			Expect(saFakeClient.ConfirmPasswordCallCount()).To(BeNumerically("==", 1))
			model, headers, _, _, _ := saFakeClient.ConfirmPasswordArgsForCall(0)
			Expect(model.Email).To(Equal("test@email.com"))
			Expect(model.ConfCode).To(Equal("898985"))
			Expect(model.ShopperId).To(Equal("90909067677"))
			Expect(model.Password).To(Equal("password#4"))
			Expect(headers).ToNot(BeNil())
			Expect(headers.Get("X-Nor-Appiid")).To(Equal("fooappiid"))
			Expect(headers.Get("True-Client-Ip")).To(Equal("1.2.3.4"))
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
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			body = `
			{
				"email":"tes‰`

			code := "my_auth_code"
			verifier := "EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			requestedPath := "/v1/password/confirm?code=" + code + "&verifier=" + verifier

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			dynamoSvc.DeleteItemReturns(&dynamodb.DeleteItemOutput{}, nil) //result, error

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(BeNumerically("==", 1))
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))

			table, key, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(key["InstallationId"]).To(Equal("fooappiid"))
			Expect(key["IPAddress"]).To(Equal("1.2.3.*"))
		})

		It("does not call customer auth confirm password", func() {
			Expect(saFakeClient.ConfirmPasswordCallCount()).To(Equal(0))
		})

		It("returns a bad request response", func() {
			responseCode := responseWriter.Code
			Expect(responseCode).To(Equal(400))
		})
	})

	Context("when a request missing email is sent ", func() {
		JustBeforeEach(func() {
			dynamoFakeClient.GetStub = func(t string, k map[string]string, o interface{}, c apm.TransactionContext) error {
				data, _ := json.Marshal(authEntity)
				json.Unmarshal(data, o)
				return nil
			}

			body = `
			{
				"email":""
			}`

			code := "my_auth_code"
			verifier := "EWKMF8RzLLcRR8ATvKDyVyY1iTEFcdU4m8imyipkH45HJotVZxeoJ80AbMPJ_s3c"
			requestedPath := "/v1/password/confirm?code=" + code + "&verifier=" + verifier

			closer = new(ioFakes.FakeCloser)
			request, err := http.NewRequest("POST", requestedPath, clients.FakeReadCloser(body, closer))
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			dynamoSvc.DeleteItemReturns(&dynamodb.DeleteItemOutput{}, nil) //result, error

			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-Clientid", "NINTERNALIOS")

			subject.ServeHTTP(responseWriter, request)
		})

		It("calls dynamo once", func() {
			Expect(dynamoFakeClient.GetCallCount()).To(Equal(1))
			Expect(dynamoFakeClient.DeleteAuthCallCount()).To(Equal(0))

			table, key, _, _ := dynamoFakeClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(key["InstallationId"]).To(Equal("fooappiid"))
			Expect(key["IPAddress"]).To(Equal("1.2.3.*"))
		})

		It("does not call customer auth forgot password", func() {
			Expect(saFakeClient.ConfirmPasswordCallCount()).To(Equal(0))
		})

		It("returns a bad request response", func() {
			responseCode := responseWriter.Code
			Expect(responseCode).To(Equal(400))
		})
	})

})
