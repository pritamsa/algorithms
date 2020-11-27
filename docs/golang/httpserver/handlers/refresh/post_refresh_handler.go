package refresh

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gitlab.nordstrom.com/sentry/authorize/authorizeconstants"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter"
	"gitlab.nordstrom.com/sentry/authorize/crypto"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/router"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/authorize/stream"
	"gitlab.nordstrom.com/sentry/authorize/try"
	sw "gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
)

type postRefreshHandler struct {
	apigeeClient apigee.Client
	dynamoClient dynamo.Client
	forterClient forter.Client
	statsdClient sw.Client
}

// NewPostRefreshHandler is the handler for the refresh post.
func NewPostRefreshHandler(
	apigeeClient apigee.Client,
	dynamoClient dynamo.Client,
	forterClient forter.Client,
	statsdClient sw.Client,
) http.Handler {
	return &postRefreshHandler{
		apigeeClient: apigeeClient,
		dynamoClient: dynamoClient,
		forterClient: forterClient,
		statsdClient: statsdClient,
	}
}

var shopperCrypto = crypto.NewDesShopperIDDecryptorEncryptor()

// Refresh Token godoc
// @Id RefreshToken
// @Summary Refresh the access token
// @Description Call with the refresh token to generate a new access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param requestBody body model.RefreshEntity true "Refresh Token"
// @Param code query string true "Server side Cryptographic Nonce Code"
// @Param verifier query string true "Client side unsigned cryptographic nonce code used for generating code in AuthInit"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {object} model.RefreshResponse
// @Failure 418 {string} string "Returned if the caller uses same code again for same Appiid"
// @Failure 401 {string} string "Returned if refresh token is invalid or expired"
// @Failure 500 {string} string "Internal Error"
// @Router /v1/refresh [post]
func (handler *postRefreshHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logger, lc := logging.GetLoggerAndContext(request)
	if handlers.VerifyFirstAuth(writer, request, handler.dynamoClient, http.StatusTeapot) {
		return
	}
	tCtx := apm.FromRequest(request)

	webShopperID := "Unauthorized"
	trueClientIP := request.Header.Get("True-Client-Ip")
	installationID := request.Header.Get("X-Nor-Appiid")
	geoInformation := request.Header.Get("X-Akamai-Edgescape")
	scope := request.Header.Get("X-Nor-Scope")

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		tCtx.NoticeError(err)
		logger.Error(lc, "Failed to read refresh body", err, map[string]interface{}{"True-Client-Ip": trueClientIP, "appiid": installationID})
		stream.PutRecord(handler.dynamoClient, trueClientIP, installationID, http.StatusTeapot, webShopperID, geoInformation, tCtx)
		http.Redirect(writer, request, "/v1/authinit", http.StatusTeapot)
		return
	}

	refresh := model.RefreshEntity{}
	err = json.Unmarshal(requestBody, &refresh)
	if err != nil {
		tCtx.NoticeError(err)
		stream.PutRecord(handler.dynamoClient, trueClientIP, installationID, http.StatusBadRequest, webShopperID, geoInformation, tCtx)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(err.Error()))
		return
	}

	validationErr := refresh.Validate()
	if validationErr != nil {
		tCtx.NoticeError(validationErr)
		stream.PutRecord(handler.dynamoClient, trueClientIP, installationID, http.StatusBadRequest, webShopperID, geoInformation, tCtx)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(validationErr.Error()))
		return
	}

	// if feature flag to disable persistent signin for all users is true, return 401 immediately
	if request.URL.Path == authorizeconstants.APIVersion1_1+router.RefreshPostPath && handler.dynamoClient.GetFeatureFlag(authorizeconstants.DisableRefreshV11, tCtx) {
		stream.PutRecord(handler.dynamoClient, trueClientIP, installationID, http.StatusUnauthorized, webShopperID, geoInformation, tCtx)
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	aResponse, err := handler.apigeeClient.RefreshToken(refresh.Token, installationID, scope, tCtx)
	if err != nil {
		tCtx.NoticeError(err)
		logger.ErrorWithHeader(lc, "ApigeeError", err, request.Header)
		if err == apigee.ErrRefreshTokenExpired || err == apigee.ErrInvalidRefreshToken || err == apigee.ErrAccessTokenNotApproved || err == apigee.ErrRefreshTokenNotApproved {
			stream.PutRecord(handler.dynamoClient, trueClientIP, installationID, http.StatusUnauthorized, webShopperID, geoInformation, tCtx)
			http.Redirect(writer, request, "/v1/authinit", http.StatusUnauthorized)
			return
		}
		http.Redirect(writer, request, "/v1/authinit", http.StatusInternalServerError)
		return
	}

	err = handler.dynamoClient.DeleteAuth(installationID, try.GetPkceIP(trueClientIP), tCtx)
	if err != nil {
		tCtx.NoticeError(err)
		logger.ErrorWithHeader(lc, "DeleteAuthFromDynamo", err, request.Header)
	}

	// v1.1 will add checking for opt-in flag in dynamodb and call forter
	if request.URL.Path == authorizeconstants.APIVersion1_1+router.RefreshPostPath {
		webShopperID, err = shopperCrypto.Decrypt(aResponse.ShopperID)
		if err != nil {
			tCtx.NoticeError(err)
			logger.ErrorWithHeader(lc, "ShopperDecryptionError", err, request.Header)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		device := model.DeviceRecord{}
		keyMap := map[string]string{model.WebShopperID_Column: webShopperID, model.DeviceID_Column: installationID}
		err := handler.dynamoClient.Get(dynamo.ShopperDeviceTable, keyMap, &device, tCtx)
		if err == dynamo.ErrNotFound {
			logger.ErrorWithHeaderData(lc, "Refresh without a previous record", err, request.Header,
				map[string]interface{}{
					"webShopperId": webShopperID,
					"appiid":       installationID,
				})

			writer.WriteHeader(http.StatusConflict)
			return
		}
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Only allow if opt-in is true and it doesn't expire
		if !device.PersistentOptIn || device.PersistentOptInExpiration < time.Now().Unix() {
			logger.InfoWithHeaderData(lc, "Persistent opted out or expired", request.Header,
				map[string]interface{}{
					"webShopperId": webShopperID,
					"appiid":       installationID,
				})

			writer.WriteHeader(http.StatusConflict)
			return
		}

		if !handler.dynamoClient.GetFeatureFlag("bypassForterOnRefresh", tCtx) {
			forterClientModel := forter.ClientModel{
				ShopperID:       webShopperID,
				InstallationID:  installationID,
				TrueClientIP:    trueClientIP,
				LoginMethodType: forter.TokenRefresh,
				HTTPRequest:     request,
				TCtx:            tCtx,
			}

			forterDecision, forterError := handler.forterClient.GetForterAccountLoginResult(forterClientModel)
			if forterError != nil {
				handler.statsdClient.Increment("custom", map[string]string{"ForterLoginError": "error"})
				logger.Error(lc, "Error while calling Forter Account Login", forterError, nil)
				writer.WriteHeader(authorizeconstants.StatusChallenged)
				return
			}

			if forterDecision.ForterDecision != forter.Approve {
				writer.WriteHeader(authorizeconstants.StatusChallenged)
				return
			}
		} else {
			handler.statsdClient.Increment("custom", map[string]string{"ForterOnRefresh": "bypass"})
		}
	}

	logger.InfoWithHeaderData(lc, "Successful refresh",
		request.Header,
		map[string]interface{}{
			"shopperId": webShopperID,
		})

	stream.PutRecord(handler.dynamoClient, trueClientIP, installationID, http.StatusOK, webShopperID, geoInformation, tCtx)

	var rResponse interface{}
	if request.URL.Path == authorizeconstants.APIVersion1_1+router.RefreshPostPath {
		rResponse = model.RefreshResponseV11{
			AccessToken:           aResponse.AccessToken,
			RefreshToken:          aResponse.RefreshToken,
			ExpiresIn:             aResponse.ExpiresIn,
			RefreshTokenExpiresIn: aResponse.RefreshTokenExpiresIn,
		}
	} else {
		rResponse = model.RefreshResponse{
			AccessToken:           aResponse.AccessToken,
			RefreshToken:          aResponse.RefreshToken,
			ExpiresIn:             strconv.Itoa(aResponse.ExpiresIn),
			RefreshTokenExpiresIn: strconv.Itoa(aResponse.RefreshTokenExpiresIn),
		}
	}
	jResponse, err := json.Marshal(rResponse)
	if err != nil {
		logger.ErrorWithHeader(lc, "Refresh Response JSON Marshal failed", err, request.Header)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(jResponse)
}

// Refresh Token godoc
// @Id RefreshToken
// @Summary Refresh the access token
// @Description Call with the refresh token to generate a new access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param requestBody body model.RefreshEntity true "Refresh Token"
// @Param code query string true "Server side Cryptographic Nonce Code"
// @Param verifier query string true "Client side unsigned cryptographic nonce code used for generating code in AuthInit"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {object} model.RefreshResponseV11
// @Failure 250 {string} string "The user has to login"
// @Failure 401 {string} string "Returned if refresh token is invalid or expired"
// @Failure 418 {string} string "Returned if the caller uses same code again for same Appiid"
// @Failure 500 {string} string "Internal Error"
// @Router /v1.1/refresh [post]
func v11shell() {
	// This empty function is a shell for /v1.1/refresh docs
}
