package authinit

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"gitlab.nordstrom.com/sentry/authcrypto"
	"gitlab.nordstrom.com/sentry/authorize/authorizeconstants"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/concurrency"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/authorize/try"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
)

type authinitHandler struct {
	dynamoClient dynamo.Client
	client       statsd_wrapper.Client
}

func NewAuthinitHandler(dynamoClient dynamo.Client) http.Handler {
	return &authinitHandler{
		dynamoClient: dynamoClient,
	}
}

// AuthInit godoc
// @Id AuthInit
// @Summary Initialize the authentication flow
// @Description Leg 1 of 3 legged Authorization calls. This call starts the auth flow to get initial
// @Description server side cryptographic nonce code.
// @Tags Authentication
// @Produce json
// @Param code query string true "Cryptographic Nonce"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {string} string "{"code": "9nlgKYYsx..."}"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Blocked Request"
// @Failure 403 {string} string "Bad Client ID"
// @Failure 418 {string} string "This is returned if the caller uses the same code multiple times for the same app iid"
// @Failure 500 {string} string "Internal Error"
// @Failure 503 {string} string "Dependency Failure"
// @Router /v1/authinit [get]
func (handler *authinitHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logger, lc := logging.GetLoggerAndContext(request)
	tCtx := apm.FromRequest(request)
	//Information From Headers
	installationID := request.Header.Get("X-Nor-Appiid")
	if installationID == "" {
		writer.WriteHeader(http.StatusBadRequest)
		logger.Info(lc, "X-Nor-Appiid header is required", nil)
		return
	}
	if handler.blockOnAppiid(installationID, writer, logger, lc, tCtx) {
		return
	}

	geoInformation := request.Header.Get("X-Akamai-Edgescape")
	if handler.blockOnGeoinformation(geoInformation, writer, logger, lc, tCtx) {
		return
	}

	// Read pkce sync
	pkceCode := request.FormValue("code")
	if handler.blockOnPkceCode(pkceCode, writer, logger, lc, tCtx) {
		return
	}

	forwardedForIPAddresses := request.Header.Get("X-Forwarded-For")
	trueClientIP := request.Header.Get("True-Client-Ip")
	if trueClientIP == "" {
		writer.WriteHeader(http.StatusBadRequest)
		logger.Info(lc, "True-Client-Ip header is required", nil)
		return
	}

	clientID := request.Header.Get("X-Nor-Clientid")
	if _, ok := authorizeconstants.ClientIds[clientID]; !ok {
		writer.WriteHeader(http.StatusForbidden)
		logger.Info(lc, fmt.Sprintf("ClientId [%v] not supported", clientID), nil)
		return
	}

	scope := request.Header.Get("X-Nor-Scope")
	pkceCodeMethod := request.FormValue("method")
	redirectURI := request.FormValue("redirect_uri")
	authorizationKey := strings.Join([]string{trueClientIP, installationID}, ":")

	key, _ := authcrypto.NewSigningKey()
	_, signature, _ := authcrypto.GenerateAuthorization(key)
	authcode := base64.URLEncoding.EncodeToString(signature)
	logger.Info(lc, fmt.Sprintf("Authorization Key: "+authorizationKey), nil)

	err := handler.dynamoClient.PutWithTTL(dynamo.PkceTable, model.AuthorizationEntity{
		IPAddress:      try.GetPkceIP(trueClientIP),
		InstallationId: installationID,
		ClientId:       clientID,
		PKCE:           pkceCode,
		PKCEMethod:     pkceCodeMethod,
		AuthCode:       authcode,
		// Use Appiid for PublicKey for now
		PubKey:         installationID,
		RedirectURI:    redirectURI,
		GeoInformation: geoInformation,
		IPChain:        forwardedForIPAddresses,
		Scope:          scope,
	}, 300, tCtx)

	if err != nil {
		tCtx.NoticeError(err)
		writer.WriteHeader(http.StatusServiceUnavailable)
		logger.Error(lc, "Error in handler", err, nil)
		return
	}

	item := model.BlacklistItem{
		Id:   pkceCode,
		Type: dynamo.BlacklistTypePkce,
	}
	rtCtx := tCtx.NewGoRoutine()
	// Write pkce blacklist async
	concurrency.Run(func() {
		handler.dynamoClient.PutWithTTL(dynamo.BlacklistTable, item, dynamo.BlacklistTTL, rtCtx)
	}, lc.TraceContext)

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(fmt.Sprintf(`{"code":"%s"}`, authcode)))
}

func (handler *authinitHandler) blockOnAppiid(appiid string, writer http.ResponseWriter, logger logging.Logger, lc logging.LogContext, tCtx apm.TransactionContext) bool {
	return handler.blockOnValue(handler.isAppiidBlacklisted, appiid, writer, logger, lc, "appiid", tCtx)
}

func (handler *authinitHandler) blockOnPkceCode(pkceCode string, writer http.ResponseWriter, logger logging.Logger, lc logging.LogContext, tCtx apm.TransactionContext) bool {
	return handler.blockOnValue(handler.isPkceCodeBlacklisted, pkceCode, writer, logger, lc, "pkceCode", tCtx)
}

func (handler *authinitHandler) blockOnValue(isValueBlacklistedFunc func(value string, tCtx apm.TransactionContext) (bool, error), value string, writer http.ResponseWriter, logger logging.Logger, lc logging.LogContext, logKey string, tCtx apm.TransactionContext) bool {
	isValueBlacklisted, err := isValueBlacklistedFunc(value, tCtx)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "SerializationException" {
			writer.WriteHeader(http.StatusBadRequest)
		} else {
			writer.WriteHeader(http.StatusServiceUnavailable)
		}
		tCtx.NoticeError(err)
		logger.Error(lc, "Error in handler", err, nil)
		return true
	}

	if isValueBlacklisted {
		writer.WriteHeader(http.StatusTeapot)
		logger.Info(lc,
			fmt.Sprintf("Blocked %v %v", logKey, value),
			map[string]interface{}{
				logKey: value,
			})
		return true
	}

	return false
}

func (handler *authinitHandler) blockOnGeoinformation(geoInformation string, writer http.ResponseWriter, logger logging.Logger, lc logging.LogContext, tCtx apm.TransactionContext) bool {
	if geoInformation == "" {
		if authorizeconstants.SERVER_ENVIRONMENT == "int" {
			// allow Int to have fake Edgescape info
			geoInformation = "georegion=242,country_code=US,region_code=AL,city=BIRMINGHAM,dma=630,msa=1000,areacode=205,county=JEFFERSON+SHELBY,fips=01073+01117,lat=33.5208,long=-86.8027,timezone=CST,zip=35201-35224+35226+35228-35229+35231-35238+35242-35244+35246+35249+35253-35255+35259-35261+35266+35282-35283+35285+35287-35288+35290-35298,continent=NA,throughput=vhigh,bw=5000,asnum=3549,location_id=0"
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte("Insufficient information - 0001"))
			return true
		}
	}

	re, _ := regexp.Compile(`country_code=(\w\w)`)
	matches := re.FindAllStringSubmatch(geoInformation, -1)

	country := matches[0][1]
	isCountryBlacklisted, err := handler.isCountryBlacklisted(country, tCtx)

	if err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		logger.Error(lc, "Error in handler", err, nil)
		return true
	}

	if isCountryBlacklisted && country != "US" {
		writer.WriteHeader(http.StatusTeapot)
		logger.Info(lc,
			fmt.Sprintf("Blocked country %v", country),
			map[string]interface{}{
				"country": country,
			})
		return true
	}

	return false
}

func (handler *authinitHandler) isCountryBlacklisted(countryCode string, tCtx apm.TransactionContext) (bool, error) {
	return handler.isKeyBlacklisted(countryCode, tCtx)
}

func (handler *authinitHandler) isAppiidBlacklisted(appiid string, tCtx apm.TransactionContext) (bool, error) {
	return handler.isKeyBlacklisted(appiid, tCtx)
}

func (handler *authinitHandler) isPkceCodeBlacklisted(pkceCode string, tCtx apm.TransactionContext) (bool, error) {
	return handler.isKeyBlacklisted(pkceCode, tCtx)
}

func (handler *authinitHandler) isKeyBlacklisted(id string, tCtx apm.TransactionContext) (bool, error) {

	var item model.BlacklistItem
	err := handler.dynamoClient.Get(dynamo.BlacklistTable, map[string]string{"Id": id}, &item, tCtx)

	// not found = no blacklist
	if err == dynamo.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
