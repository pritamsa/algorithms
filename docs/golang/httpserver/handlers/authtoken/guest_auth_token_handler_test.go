package authtoken

import (
	"context"
	"encoding/json"
	"errors"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	"net/http"
	"net/http/httptest"
	"regexp"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee/apigeefakes"
	"gitlab.nordstrom.com/sentry/authorize/crypto"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"
)

var _ = Describe("GuestAuthTokenHandler", func() {
	var (
		subject             http.Handler
		fakeGologger        *gologgerfakes.FakeLogger
		apigeeFakeClient    *apigeefakes.FakeClient
		responseWriter      *httptest.ResponseRecorder
		clientRequest       *http.Request
		authorizationHeader string
		shopperCrypto       = crypto.NewDesShopperIDDecryptorEncryptor()
		tCtx                *apmfakes.FakeTransactionContext
	)

	BeforeEach(func() {
		tCtx = &apmfakes.FakeTransactionContext{}
		fakeGologger = &gologgerfakes.FakeLogger{}
		logging.CreateSingleLoggerForTest(fakeGologger, "", "")
		apigeeFakeClient = &apigeefakes.FakeClient{}
		var err error
		clientRequest, err = http.NewRequest("POST", "/", nil)
		clientRequest.Header.Add("TraceContext", "abcdefg")
		clientRequest.Header.Add("X-Nor-Appiid", "myappiid")
		Expect(err).NotTo(HaveOccurred())
		prepareApm := func(r *http.Request) *http.Request {
			ctx := context.WithValue(r.Context(), "tCtx", tCtx)
			return r.WithContext(ctx)
		}
		clientRequest = prepareApm(clientRequest)
		authorizationHeader = "Basic dGVzdGlkOnNlY3JldDEyMzQ="
	})
	JustBeforeEach(func() {
		if authorizationHeader != "" {
			clientRequest.Header.Add("Authorization", authorizationHeader)
		}
		responseWriter = httptest.NewRecorder()
		subject = NewGuestAuthTokenHandler(apigeeFakeClient)
		subject.ServeHTTP(responseWriter, clientRequest)
	})

	Context("when the request is valid", func() {
		BeforeEach(func() {
			apigeeFakeClient.GetTokenReturns(apigee.BearerToken{
				AccessToken: "abcd",
				ExpiresIn:   "60",
				TokenType:   "test",
			}, nil)
		})

		It("calls token provider client", func() {
			Expect(apigeeFakeClient.GetTokenCallCount()).To(Equal(1))

			args := apigeeFakeClient.GetTokenArgsForCall(0)
			shopperID := args.ShopperID

			Expect(shopperID).NotTo(BeEmpty())
			Expect(len(shopperID)).To(Equal(55))
			decryptedShopperID, err := shopperCrypto.Decrypt(shopperID)
			Expect(err).To(BeNil())
			matched, err := regexp.MatchString("[a-z0-9]{32}", decryptedShopperID)
			Expect(err).To(BeNil())
			Expect(matched).To(BeTrue())

			Expect(args.Appiid).To(Equal("myappiid"))
			Expect(args.ClientID).To(Equal("testid"))
			Expect(args.ClientSecret).To(Equal("secret1234"))
		})

		It("returns a valid response", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
			Expect(responseWriter.Header().Get("Content-Type")).To(Equal("application/json; charset=utf-8"))
			var tokenResponse CreateAuthTokenResponse
			err := json.Unmarshal(responseWriter.Body.Bytes(), &tokenResponse)
			Expect(err).To(BeNil())
			Expect(tokenResponse.AccessToken).To(Equal("abcd"))
			Expect(tokenResponse.ExpiresIn).To(Equal("60"))
			Expect(tokenResponse.TokenType).To(Equal("test"))

			webShopperID := tokenResponse.WebShopperID
			Expect(webShopperID).NotTo(BeEmpty())
			Expect(len(webShopperID)).To(Equal(32))
			matched, err := regexp.MatchString("[a-z0-9]{32}", webShopperID)
			Expect(err).To(BeNil())
			Expect(matched).To(BeTrue())

			shopperID := tokenResponse.ShopperID
			Expect(shopperID).NotTo(BeEmpty())
			Expect(len(shopperID)).To(Equal(55))
			decryptedShopperID, err := shopperCrypto.Decrypt(shopperID)
			Expect(err).To(BeNil())
			matched, err = regexp.MatchString("[a-z0-9]{32}", decryptedShopperID)
			Expect(err).To(BeNil())
			Expect(matched).To(BeTrue())
		})
	})

	Context("when the request is missing Authorization header", func() {
		BeforeEach(func() {
			authorizationHeader = ""
		})

		It("does not call the AuthTokenUseCase", func() {
			Expect(apigeeFakeClient.GetTokenCallCount()).To(Equal(0))
		})

		It("returns Unauthorized", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("when the request has malformed Authorization header", func() {
		BeforeEach(func() {
			authorizationHeader = "HelloWorld"
		})

		It("does not call the AuthTokenUseCase", func() {
			Expect(apigeeFakeClient.GetTokenCallCount()).To(Equal(0))
		})

		It("returns Unauthorized", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})
	})

	Context("when token provider returns an error", func() {
		BeforeEach(func() {
			apigeeFakeClient.GetTokenReturns(apigee.BearerToken{}, errors.New("upstreamerror"))
		})

		It("calls the AuthTokenUseCase", func() {
			Expect(apigeeFakeClient.GetTokenCallCount()).To(Equal(1))
		})

		It("returns InternalServerError with empty body", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusInternalServerError))
			Expect(responseWriter.Body.String()).To(BeEmpty())
		})
	})
})
