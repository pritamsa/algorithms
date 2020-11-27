package authinit_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm/apmfakes"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo/dynamofakes"
	authinitHandler "gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/authinit"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/authorize/try"
	"gitlab.nordstrom.com/sentry/gologger/gologgerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func PrepareApm(r *http.Request, tCtx apm.TransactionContext) *http.Request {
	ctx := context.WithValue(r.Context(), "tCtx", tCtx)
	return r.WithContext(ctx)
}

var _ = Describe("the authinit handler", func() {
	var (
		subject http.Handler

		requestedPath  string
		responseWriter *httptest.ResponseRecorder

		fakeDynamoClient *dynamofakes.FakeClient
		request          *http.Request
		tCtx             *apmfakes.FakeTransactionContext
		err              error
	)

	const (
		edgescape = "georegion=242,country_code=US,region_code=AL,city=BIRMINGHAM,dma=630,msa=1000,areacode=205,county=JEFFERSON+SHELBY,fips=01073+01117,lat=33.5208,long=-86.8027,timezone=CST,zip=35201-35224+35226+35228-35229+35231-35238+35242-35244+35246+35249+35253-35255+35259-35261+35266+35282-35283+35285+35287-35288+35290-35298,continent=NA,throughput=vhigh,bw=5000,asnum=3549,location_id=0"
	)

	BeforeEach(func() {
		logging.CreateSingleLoggerForTest(&gologgerfakes.FakeLogger{}, "", "")
		fakeDynamoClient = &dynamofakes.FakeClient{}

		subject = authinitHandler.NewAuthinitHandler(fakeDynamoClient)
		responseWriter = httptest.NewRecorder()
		tCtx = &apmfakes.FakeTransactionContext{}
		tCtx.NewGoRoutineReturns(tCtx)
	})

	assertApm := func(success bool) {
		var subCalls int
		if success {
			subCalls = 1
		} else {
			subCalls = 0
		}
		It("Makes the appropriate calls", func() {
			Expect(tCtx.NewGoRoutineCallCount()).To(Equal(subCalls))
		})
	}

	prepareApm := func(r *http.Request) *http.Request {
		return PrepareApm(r, tCtx)
	}

	Context("when dynamo client returns no error", func() {
		BeforeEach(func() {
			requestedPath = "/v1/authinit"
			request, err = http.NewRequest("GET", requestedPath, nil)
			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Akamai-Edgescape", edgescape)
			request.Header.Set("X-Forwarded-For", "2.3.4.5")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-ClientId", "acc8d2c1-6f3e-4990-b363-d43ec8b68cac")
			request.Header.Set("X-Nor-Scope", "testscope")

			request.Form = map[string][]string{
				"code":         {"testpkcecode"},
				"method":       {"testmethod"},
				"redirect_uri": {"testuri"},
			}

			request = prepareApm(request)

			fakeDynamoClient.GetReturns(dynamo.ErrNotFound)

			Expect(err).ToNot(HaveOccurred())

		})

		JustBeforeEach(func() {
			subject.ServeHTTP(responseWriter, request)
		})

		It("sets the response code to 200", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})

		It("sets correct headers", func() {
			Expect(responseWriter.Header()["Content-Type"][0]).To(Equal("application/json"))
		})

		It("calls dynamo client", func() {
			try.Until(func() bool {
				return fakeDynamoClient.GetCallCount() == 3
			})
			Expect(fakeDynamoClient.GetCallCount()).To(Equal(3))

			table, in, out, context := fakeDynamoClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("fooappiid"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())

			table, in, out, context = fakeDynamoClient.GetArgsForCall(1)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("US"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())

			table, in, out, context = fakeDynamoClient.GetArgsForCall(2)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("testpkcecode"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())

			try.Until(func() bool {
				return fakeDynamoClient.PutWithTTLCallCount() == 2
			})
			Expect(fakeDynamoClient.PutWithTTLCallCount()).To(Equal(2))

			table, item, ttl, context := fakeDynamoClient.PutWithTTLArgsForCall(0)
			ae := item.(model.AuthorizationEntity)
			Expect(table).To(Equal(dynamo.PkceTable))
			Expect(ttl).To(Equal(300))
			Expect(ae.IPAddress).To(Equal("1.2.3.*"))
			Expect(ae.InstallationId).To(Equal("fooappiid"))
			Expect(ae.ClientId).To(Equal("acc8d2c1-6f3e-4990-b363-d43ec8b68cac"))
			Expect(ae.PKCE).To(Equal("testpkcecode"))
			Expect(ae.PKCEMethod).To(Equal("testmethod"))
			Expect(len(ae.AuthCode)).To(Equal(88))
			Expect(ae.SubmitCode).To(Equal(""))
			Expect(ae.PostVersion).To(Equal(""))
			Expect(ae.PubKey).To(Equal("fooappiid"))
			Expect(ae.RedirectURI).To(Equal("testuri"))
			Expect(ae.AuthorizationTimeout).To(Equal(""))
			Expect(ae.Scope).To(Equal("testscope"))
			Expect(ae.GeoInformation).To(Equal("georegion=242,country_code=US,region_code=AL,city=BIRMINGHAM,dma=630,msa=1000,areacode=205,county=JEFFERSON+SHELBY,fips=01073+01117,lat=33.5208,long=-86.8027,timezone=CST,zip=35201-35224+35226+35228-35229+35231-35238+35242-35244+35246+35249+35253-35255+35259-35261+35266+35282-35283+35285+35287-35288+35290-35298,continent=NA,throughput=vhigh,bw=5000,asnum=3549,location_id=0"))
			Expect(ae.IPChain).To(Equal("2.3.4.5"))
			Expect(ae.AuthorizationAttempts).To(BeEmpty())
			Expect(ae.AuthGetTimestamp).To(Equal(""))

			table, item, ttl, context = fakeDynamoClient.PutWithTTLArgsForCall(1)
			be := item.(model.BlacklistItem)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(be).ToNot(BeNil())
			Expect(be.Id).To(Equal("testpkcecode"))
			Expect(be.Type).To(Equal("pkce"))
			Expect(ttl).To(Equal(3600))
			Expect(context).ToNot(BeNil())

		})

		It("returns auth code", func() {
			Expect(responseWriter.Body.String()).To(ContainSubstring(`{"code":"`))
			// URLEncoding will not use "+"
			Expect(responseWriter.Body.String()).NotTo(ContainSubstring("+"))
		})

		assertApm(true)
	})

	Context("when missing geo info", func() {
		JustBeforeEach(func() {
			fakeDynamoClient.GetReturns(dynamo.ErrNotFound)

			requestedPath = "/v1/authinit"
			request, err := http.NewRequest("GET", requestedPath, nil)
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("sets the response code to 401", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusUnauthorized))
		})

		It("calls dynamo client zero times", func() {
			Expect(fakeDynamoClient.PutCallCount()).To(Equal(0))
		})

		It("returns client error", func() {
			Expect(responseWriter.Body.String()).To(Equal("Insufficient information - 0001"))
		})

		It("calls dynamo client once", func() {
			Expect(fakeDynamoClient.GetCallCount()).To(Equal(1))
			table, in, out, context := fakeDynamoClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("fooappiid"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())
		})

		assertApm(false)
	})

	Context("when client is not in the list", func() {
		JustBeforeEach(func() {
			fakeDynamoClient.GetReturns(dynamo.ErrNotFound)

			requestedPath = "/v1/authinit"
			request, err = http.NewRequest("GET", requestedPath, nil)
			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Akamai-Edgescape", edgescape)
			request.Header.Set("X-Forwarded-For", "2.3.4.5")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-ClientId", "Any values")
			request.Header.Set("X-Nor-Scope", "testscope")

			request.Form = map[string][]string{
				"code":         {"testpkcecode"},
				"method":       {"testmethod"},
				"redirect_uri": {"testuri"},
			}

			request = prepareApm(request)
			subject.ServeHTTP(responseWriter, request)
		})

		It("returns client error", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusForbidden))
		})
	})

	Context("when client is in the list(MWPClientID)", func() {
		JustBeforeEach(func() {
			fakeDynamoClient.GetReturns(dynamo.ErrNotFound)

			requestedPath = "/v1/authinit"
			request, err = http.NewRequest("GET", requestedPath, nil)
			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Header.Set("X-Akamai-Edgescape", edgescape)
			request.Header.Set("X-Forwarded-For", "2.3.4.5")
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Nor-ClientId", "7f6f2a5a-229e-47a9-8f54-145451c836a6")
			request.Header.Set("X-Nor-Scope", "testscope")

			request.Form = map[string][]string{
				"code":         {"testpkcecode"},
				"method":       {"testmethod"},
				"redirect_uri": {"testuri"},
			}

			request = prepareApm(request)
			subject.ServeHTTP(responseWriter, request)
		})

		It("returns no error", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})
	})

	Context("when dynamo client Get an error", func() {
		JustBeforeEach(func() {
			fakeDynamoClient.GetReturns(errors.New("hii"))

			requestedPath = "/v1/authinit"
			request, err := http.NewRequest("GET", requestedPath, nil)
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Akamai-Edgescape", edgescape)
			Expect(err).ToNot(HaveOccurred())
			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("sets the response code to 500", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusServiceUnavailable))
		})

		It("calls dynamo client once", func() {
			Expect(fakeDynamoClient.PutCallCount()).To(Equal(0))
			Expect(fakeDynamoClient.GetCallCount()).To(Equal(1))
			table, in, out, context := fakeDynamoClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("fooappiid"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())
		})

		assertApm(false)
	})

	Context("when a country is blacklisted", func() {
		JustBeforeEach(func() {
			// first call does not find an appiid blacklist item
			fakeDynamoClient.GetReturnsOnCall(0, dynamo.ErrNotFound)
			// second call finds country blacklist item
			fakeDynamoClient.GetReturnsOnCall(1, nil)

			requestedPath = "/v1/authinit"
			request, err := http.NewRequest("GET", requestedPath, nil)
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Akamai-Edgescape", "georegion=242,country_code=RU,region_code=AL")
			Expect(err).ToNot(HaveOccurred())

			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("sets the response code to 418", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusTeapot))
		})

		It("calls dynamo client twice", func() {
			Expect(fakeDynamoClient.PutWithTTLCallCount()).To(Equal(0))
			Expect(fakeDynamoClient.GetCallCount()).To(Equal(2))

			table, in, out, context := fakeDynamoClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("fooappiid"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())

			table, in, out, context = fakeDynamoClient.GetArgsForCall(1)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("RU"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())
		})

		assertApm(false)
	})

	Context("when US is never blocked even if US is in the black list", func() {
		JustBeforeEach(func() {
			// first call does not find an appiid blacklist item
			fakeDynamoClient.GetReturnsOnCall(0, dynamo.ErrNotFound)
			// second call finds country blacklist item
			fakeDynamoClient.GetReturnsOnCall(1, nil)
			// third call does not find a pkce blacklist item
			fakeDynamoClient.GetReturnsOnCall(2, dynamo.ErrNotFound)

			requestedPath = "/v1/authinit"
			request, err := http.NewRequest("GET", requestedPath, nil)
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Akamai-Edgescape", "georegion=242,country_code=US,region_code=AL")
			request.Header.Set("X-Nor-ClientId", "acc8d2c1-6f3e-4990-b363-d43ec8b68cac")
			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Form = map[string][]string{"code": {"testpkcecode"}}
			Expect(err).ToNot(HaveOccurred())

			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("sets the response code to 200", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusOK))
		})

		It("calls dynamo client", func() {
			try.Until(func() bool { return fakeDynamoClient.GetCallCount() == 3 })
			Expect(fakeDynamoClient.GetCallCount()).To(Equal(3))

			table, in, out, context := fakeDynamoClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("fooappiid"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())

			table, in, out, context = fakeDynamoClient.GetArgsForCall(1)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("US"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())

			table, in, out, context = fakeDynamoClient.GetArgsForCall(2)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("testpkcecode"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())
		})

		It("calls dynamo client", func() {
			try.Until(func() bool { return fakeDynamoClient.PutWithTTLCallCount() > 1 })
			Expect(fakeDynamoClient.PutWithTTLCallCount()).To(Equal(2))
		})

		assertApm(true)
	})

	Context("when an appiid is blacklisted", func() {
		JustBeforeEach(func() {

			//var item = model.BlacklistItem{Id: "fooappiid", Type: "appiid"}
			fakeDynamoClient.GetReturns(nil)

			requestedPath = "/v1/authinit"
			request, err := http.NewRequest("GET", requestedPath, nil)
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Akamai-Edgescape", "georegion=242,country_code=US,region_code=AL")
			request.Header.Set("X-Nor-ClientId", "acc8d2c1-6f3e-4990-b363-d43ec8b68cac")
			Expect(err).ToNot(HaveOccurred())

			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("sets the response code to 418", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusTeapot))
		})

		It("calls dynamo client once", func() {
			Expect(fakeDynamoClient.PutCallCount()).To(Equal(0))
			Expect(fakeDynamoClient.GetCallCount()).To(Equal(1))
			table, in, out, context := fakeDynamoClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("fooappiid"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())
		})

		assertApm(false)
	})

	Context("when dynamo client PutAuthWithTTLReturns an error", func() {
		JustBeforeEach(func() {
			// first call does not find an appiid blacklist item
			fakeDynamoClient.GetReturnsOnCall(0, dynamo.ErrNotFound)
			// second call does not find country blacklist item
			fakeDynamoClient.GetReturnsOnCall(1, dynamo.ErrNotFound)
			// third call does not find a pkce blacklist item
			fakeDynamoClient.GetReturnsOnCall(2, dynamo.ErrNotFound)

			fakeDynamoClient.PutWithTTLReturns(errors.New("hii"))

			requestedPath = "/v1/authinit"
			request, err := http.NewRequest("GET", requestedPath, nil)
			request.Header.Set("X-Nor-Appiid", "fooappiid")
			request.Header.Set("X-Akamai-Edgescape", edgescape)
			request.Header.Set("X-Nor-ClientId", "acc8d2c1-6f3e-4990-b363-d43ec8b68cac")
			request.Header.Set("True-Client-Ip", "1.2.3.4")
			request.Form = map[string][]string{"code": {"testpkcecode"}}
			Expect(err).ToNot(HaveOccurred())

			request = prepareApm(request)

			subject.ServeHTTP(responseWriter, request)
		})

		It("sets the response code to 500", func() {
			Expect(responseWriter.Code).To(Equal(http.StatusServiceUnavailable))
		})

		It("calls dynamo client", func() {
			try.Until(func() bool { return fakeDynamoClient.GetCallCount() == 3 })
			Expect(fakeDynamoClient.GetCallCount()).To(Equal(3))

			table, in, out, context := fakeDynamoClient.GetArgsForCall(0)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("fooappiid"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())

			table, in, out, context = fakeDynamoClient.GetArgsForCall(1)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("US"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())

			table, in, out, context = fakeDynamoClient.GetArgsForCall(2)
			Expect(table).To(Equal(dynamo.BlacklistTable))
			Expect(in).ToNot(BeNil())
			Expect(in["Id"]).To(Equal("testpkcecode"))
			Expect(out).ToNot(BeNil())
			Expect(context).ToNot(BeNil())
		})

		It("calls dynamo client", func() {
			Expect(fakeDynamoClient.PutWithTTLCallCount()).To(Equal(1))
		})

		assertApm(false)
	})
})
