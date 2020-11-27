package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"gitlab.nordstrom.com/sentry/authorize/clients/apigee"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/authorize/stream"
	"gitlab.nordstrom.com/sentry/authorize/try"
	sw "gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
)

// AuthRequest is the incoming auth request body
type AuthRequest struct {
	Code     string `json:"code"`
	Verifier string `json:"verifier"`
}

// SignInRequest is the sign in request body
type SignInRequest struct {
	ShopperAuthClient shopperauth.Client
	Username          string
	Password          string
}

// ApigeeSignOutRequest is the sign out request
type ApigeeSignOutRequest struct {
	ApigeeClient apigee.Client
	AccessToken  string
}

// IsEmailInvalid returns whether an email is invalid
func IsEmailInvalid(email string) bool {
	return strings.Contains(email, " ") || strings.Contains(email, ",") || !strings.Contains(email, "@") || !strings.Contains(email, ".") || strings.HasSuffix(email, ".")
}

// VerifyFirstAuth will validate and return whether the caller should continue processing the request or return early
func VerifyFirstAuth(writer http.ResponseWriter, request *http.Request, dynamoClient dynamo.Client, errorStatusCode int) bool {
	tCtx := apm.FromRequest(request)
	defer tCtx.Segment("VerifyFirstAuth").End()
	logger, lc := logging.GetLoggerAndContext(request)

	webShopperID := "Unauthorized"
	trueClientIP := request.Header.Get("True-Client-Ip")
	installationID := request.Header.Get("X-Nor-Appiid")
	clientID := request.Header.Get("X-Nor-ClientId")
	geoInformation := request.Header.Get("X-Akamai-Edgescape")

	queryVals := request.URL.Query()
	authCode := queryVals.Get("code")
	if authCode == "" {
		logger.InfoWithHeaderData(lc, "Error 0001 - Missing code query param", request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return true
	}

	authVerifier := queryVals.Get("verifier")

	if authVerifier == "" {
		logger.InfoWithHeaderData(lc, "Error 0002 - Missing code verifier param", request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return true
	}
	authorizationKey := strings.Join([]string{trueClientIP, installationID}, ":")

	authEntity, rerr := try.GetAuthEntity(dynamoClient, installationID, trueClientIP,
		func(entity model.AuthorizationEntity, ierr error) bool { return ierr == dynamo.ErrNotFound },
		tCtx)
	if rerr != nil {
		if rerr == dynamo.ErrNotFound {
			logger.InfoWithHeaderData(lc, fmt.Sprintf("Error 0003 - No cache entry for %s", authorizationKey), request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
			stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
			http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
			return true
		}
		logger.ErrorWithHeaderData(lc, fmt.Sprintf("Error getting cache entry for %s", authorizationKey), rerr, request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		http.Redirect(writer, request, "/v1/authinit", http.StatusInternalServerError)
		return true
	}

	if clientID != authEntity.ClientId {
		statusCode := http.StatusForbidden
		logger.InfoWithHeaderData(lc, fmt.Sprintf("Client ID [%v] does not match [%v]", clientID, authEntity.ClientId), request.Header, map[string]interface{}{"errorStatusCode": statusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, statusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", statusCode)
		return true
	}

	if authCode != authEntity.AuthCode {
		logger.InfoWithHeaderData(lc, fmt.Sprintf("Error 0004 - code query param:%s != cache entry:%s", authCode, authEntity.AuthCode), request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return true
	}

	if !authEntity.IsVerifierParamValid(authVerifier) {
		logger.InfoWithHeaderData(lc, "Error 0005 - verifier is not valid", request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return true
	}
	return false
}

// VerifyAuth will ...
func VerifyAuth(writer http.ResponseWriter, request *http.Request, statsdClient sw.Client, dynamoClient dynamo.Client) ([]byte, error) {
	tCtx := apm.FromRequest(request)
	defer tCtx.Segment("VerifyAuth").End()
	errorStatusCode := http.StatusUnauthorized
	logger, lc := logging.GetLoggerAndContext(request)

	webShopperID := "Unauthorized"
	trueClientIP := request.Header.Get("True-Client-Ip")
	installationID := request.Header.Get("X-Nor-Appiid")
	clientID := request.Header.Get("X-Nor-ClientId")
	geoInformation := request.Header.Get("X-Akamai-Edgescape")

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logger.ErrorWithHeaderData(lc, "Error 0006 - verifyAuth", err, request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return []byte{}, err
	}

	var authRequest AuthRequest
	err = json.Unmarshal(body, &authRequest)
	if err != nil {
		logger.ErrorWithHeaderData(lc, "Error 0007 - verifyAuth", err, request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return []byte{}, err
	}

	submitCode, authVerifier := authRequest.Code, authRequest.Verifier
	if submitCode == "" || authVerifier == "" {
		logger.InfoWithHeaderData(lc, "Error 0008 - code or verifier is empty", request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return []byte{}, errors.New("Unauthorized")
	}
	authorizationKey := strings.Join([]string{trueClientIP, installationID}, ":")

	authEntity, err := try.GetAuthEntity(dynamoClient, installationID, trueClientIP,
		func(entity model.AuthorizationEntity, ierr error) bool { return entity.SubmitCode == "" },
		tCtx)
	if err != nil {
		if err == dynamo.ErrNotFound {
			logger.InfoWithHeaderData(lc, fmt.Sprintf("Error 0009 - No cache entry for %s", authorizationKey), request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
			stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
			http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
			return []byte{}, err
		}
		logger.ErrorWithHeaderData(lc, fmt.Sprintf("Error getting cache entry for %s", authorizationKey), err, request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		http.Redirect(writer, request, "/v1/authinit", http.StatusInternalServerError)
		return []byte{}, err
	}

	if clientID != authEntity.ClientId {
		statusCode := http.StatusForbidden
		logger.InfoWithHeaderData(lc, fmt.Sprintf("Client ID [%v] does not match [%v]", clientID, authEntity.ClientId), request.Header, map[string]interface{}{"errorStatusCode": statusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, statusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", statusCode)
		return []byte{}, errors.New("Forbidden")
	}

	userAgent := request.Header.Get("User-Agent")
	duration, err := getAuthDuration(authEntity)
	if err != nil {
		logger.Error(lc, fmt.Sprintf("error getting duration between Authorize GET and POST for %s", authorizationKey), err, nil)
	}

	addDelayForBots(logger, lc, duration, dynamoClient, tCtx)

	err = recordAuthPostLatency(duration, userAgent, clientID, request.URL.Path, statsdClient)

	if err != nil {
		logger.Error(lc, fmt.Sprintf("error recording auth post latency for %s", authorizationKey), err, nil)
	}

	if submitCode != authEntity.SubmitCode {
		logger.InfoWithHeaderData(lc, fmt.Sprintf("Error 0010 - code in request body:%s != cache entry:%s", submitCode, authEntity.SubmitCode), request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return []byte{}, errors.New("Unauthorized")
	}

	if authVerifier == "" {
		logger.InfoWithHeaderData(lc, "Error 0011 - verifier is empty", request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return []byte{}, errors.New("Unauthorized")
	}

	if !authEntity.IsVerifierParamValid(authVerifier) {
		logger.InfoWithHeaderData(lc, "Error 0012 - verifier is not valid", request.Header, map[string]interface{}{"errorStatusCode": errorStatusCode})
		stream.PutRecord(dynamoClient, trueClientIP, installationID, errorStatusCode, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", errorStatusCode)
		return []byte{}, errors.New("Unauthorized")
	}

	return body, nil
}

func addDelayForBots(logger logging.Logger, lc logging.LogContext, d time.Duration, c dynamo.Client, tCtx apm.TransactionContext) {
	if d > 0 && d.Seconds() < 1.2 {
		t := c.GetSigninDelay(lc.TraceContext, tCtx)
		if t > 0 {
			time.Sleep(time.Duration(t) * time.Millisecond)
			logger.Info(lc, fmt.Sprintf("Added a delay of %v milliseconds because bot took %v milliseconds", t, d.Seconds()*1000), nil)
		}
	}
}

func getAuthDuration(e model.AuthorizationEntity) (time.Duration, error) {
	t := e.AuthGetTimestamp
	if t == "" {
		return 0, errors.New("auth get timestamp is empty")
	}
	before, err := time.Parse(time.RFC3339Nano, t)
	if err != nil {
		return 0, err
	}
	return time.Since(before), nil
}

func recordAuthPostLatency(d time.Duration, userAgent string, norClient string, url string, client sw.Client) error {
	if d == 0 {
		return errors.New("auth get timestamp is empty")
	}
	client.Timing("auth_post_latency", d, map[string]string{
		sw.Endpoint:  url,
		sw.UserAgent: userAgent,
		sw.NorClient: norClient,
	})
	return nil
}
