package verify_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/flags"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/clients"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	ioFakes "gitlab.nordstrom.com/sentry/authorize/clients/clientsfakes"
	. "gitlab.nordstrom.com/sentry/authorize/clients/verify"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/tracecontext"
	"gitlab.nordstrom.com/sentry/gohttp/client/clientfakes"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"
)

var _ = Describe("VerifyClientAssociate", func() {
	var (
		httpClient         *clientfakes.FakeClient
		client             Client
		httpResponse       http.Response
		closer             *ioFakes.FakeCloser
		err                error
		fakeApmTransaction apmfakes.FakeTransactionContext
	)

	BeforeEach(func() {
		fakeGologger := &gologgerfakes.FakeLogger{}
		logging.CreateSingleLoggerForTest(fakeGologger, "", "")
		httpClient = new(clientfakes.FakeClient)
		fakeApmTransaction = apmfakes.FakeTransactionContext{}
		fakeApmTransaction.PrepareRequestStub = func(request *http.Request) *http.Request {
			return request
		}
		fakeApmTransaction.SegmentReturns(&apmfakes.FakeSegment{})
		fakeApmTransaction.ContextReturns(tracecontext.NewTraceContext())
		client = NewVerifyClient(httpClient,
			flags.VerifyArgs{
				"https://abc123.execute-api.us-west-2.amazonaws.com/env",
				"arn:aws:iam::744156072783:role/verify-int-execute-api-role"},
		)
		closer = new(ioFakes.FakeCloser)
		httpResponse = http.Response{Body: clients.FakeReadCloser("", closer)}
	})

	assertApm := func() {
		It("executed apm properly", func() {
			Expect(fakeApmTransaction.PrepareRequestCallCount()).To(Equal(1))
		})
	}

	Describe("Challenge Init Calls", func() {
		Context("That are successful", func() {
			var resp string
			BeforeEach(func() {
				httpClient.DoStub = func(ctx context.Context, r *http.Request) (resp *http.Response, err error) {
					httpResponse.StatusCode = http.StatusOK
					m := make(map[string]interface{})
					var body []byte
					body, _ = ioutil.ReadAll(r.Body)
					r.Body = ioutil.NopCloser(bytes.NewReader(body))
					json.Unmarshal(body, &m)
					httpResponse.Body = clients.FakeReadCloser(fmt.Sprintf("{\"sessionId\": \"%s\"}", m["sessionId"]), closer)
					return &httpResponse, nil
				}
				headers := map[string]string{
					"Via":                          "1.1 08f323deadbeefa7af34d5feb414ce27.cloudfront.net (CloudFront)",
					"Accept-Language":              "en-US,en;q=0.8",
					"CloudFront-Is-Desktop-Viewer": "true",
					"CloudFront-Is-SmartTV-Viewer": "false",
					"CloudFront-Is-Mobile-Viewer":  "false",
					"X-Forwarded-For":              "127.0.0.1, 127.0.0.2",
					"CloudFront-Viewer-Country":    "US",
					"Accept":                       "application/json",
					"Upgrade-Insecure-Requests":    "1",
					"X-Forwarded-Port":             "443",
					"Host":                         "1234567890.execute-api.us-east-1.amazonaws.com",
					"X-Forwarded-Proto":            "https",
					"X-Amz-Cf-Id":                  "cDehVQoZnx43VYQb9j2-nvCh-9z396Uhbp027Y2JvkCPNLmGJHqlaA==",
					"CloudFront-Is-Tablet-Viewer":  "false",
					"Cache-Control":                "max-age=0",
					"User-Agent":                   "Go-http-client/1.1",
					"CloudFront-Forwarded-Proto":   "https",
					"Accept-Encoding":              "gzip, deflate, sdch",
					"True-Client-Ip":               "10.1.1.1",
					"Tracecontext":                 "687350AE-6E3C-4ED2-9614-93C41FB1888E",
					"X-Nor-Scope":                  "WebRegistered",
				}
				req, _ := http.NewRequest(http.MethodPost, "https://test.net", bytes.NewReader([]byte("")))
				for key, value := range headers {
					req.Header.Set(key, value)
				}
				siResponse := shopperauth.SigninResponse{
					AccessToken:           "mQSS79Fk9UXSKaQI942nVebuIcrW",
					RefreshToken:          "LKrrGlrE5hhuHtOTLjr76qGs4GXQsBAd",
					WebShopperID:          "CEB3EDBDA13B4FE1912DE86531AC7A42",
					ShopperID:             "lgK2_ZXptGdsJEvY3tyEB_CPAs8anoqX1LfHDAiKtjx_Df7iBrFv-w2",
					ShopperToken:          "token",
					RefreshTokenExpiresIn: 1209599,
					ExpiresIn:             10000,
				}
				challengeInitRequest := shopperauth.ChallengeInitRequest{
					Appiid:              "3AD505CB-B77D-E19F-5A34-2A7337D06366",
					Email:               "test@test.com",
					SigninResponse:      siResponse,
					ForterCorrelationID: "HGJ7512345H3DE",
				}
				resp, err = client.ChallengeInit(challengeInitRequest, req.Header, &fakeApmTransaction)
			})
			It("Makes one outbound call", func() {
				Expect(httpClient.DoCallCount()).To(Equal(1))
			})

			It("Constructs the request properly", func() {
				tc, actualRequest := httpClient.DoArgsForCall(0)
				Expect(tc).ToNot(BeNil())
				Expect(actualRequest.URL.Scheme).To(Equal("https"))
				Expect(actualRequest.URL.Path).To(HavePrefix("/env/challenge/init"))
				Expect(actualRequest.Method).To(Equal(http.MethodPost))
				Expect(len(actualRequest.Header)).ToNot(Equal(0))
				Expect(actualRequest.Header.Get(XNorScope)).To(Equal("WebRegistered"))
				Expect(actualRequest.Header.Get(UserAgent)).To(Equal("Go-http-client/1.1"))
				Expect(actualRequest.Header.Get(TrueClientIp)).To(Equal("10.1.1.1"))
				chInitReq := shopperauth.ChallengeInitRequest{}
				body, _ := ioutil.ReadAll(actualRequest.Body)
				json.Unmarshal(body, &chInitReq)
				Expect(chInitReq.Appiid).To(Equal("3AD505CB-B77D-E19F-5A34-2A7337D06366"))
				Expect(chInitReq.Email).To(Equal("test@test.com"))
				Expect(chInitReq.ForterCorrelationID).To(Equal("HGJ7512345H3DE"))
				Expect(chInitReq.SessionID).NotTo(BeEmpty())
				Expect(chInitReq.SigninResponse.WebShopperID).To(Equal("CEB3EDBDA13B4FE1912DE86531AC7A42"))
				Expect(chInitReq.SigninResponse.AccessToken).To(Equal("mQSS79Fk9UXSKaQI942nVebuIcrW"))
				Expect(chInitReq.SigninResponse.RefreshToken).To(Equal("LKrrGlrE5hhuHtOTLjr76qGs4GXQsBAd"))
				Expect(chInitReq.SigninResponse.RefreshTokenExpiresIn).To(Equal(1209599))
				Expect(chInitReq.SigninResponse.ExpiresIn).To(Equal(10000))
				Expect(chInitReq.SigninResponse.ShopperID).To(Equal("lgK2_ZXptGdsJEvY3tyEB_CPAs8anoqX1LfHDAiKtjx_Df7iBrFv-w2"))
				Expect(chInitReq.SigninResponse.ShopperToken).To(Equal("token"))
			})

			It("Closes the body", func() {
				Expect(closer.CloseCallCount()).To(Equal(1))
			})

			It("Receives an appropriate response", func() {
				m := make(map[string]string)
				json.Unmarshal([]byte(resp), &m)
				Expect(err).To(BeNil())
				val, exists := m["sessionId"]
				Expect(exists).To(BeTrue())
				Expect(val).NotTo(BeEmpty())
			})

			assertApm()

		})

	})
})
