package signout

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.nordstrom.com/sentry/authorize/clients/apigee"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/concurrency"
	"gitlab.nordstrom.com/sentry/authorize/crypto"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	sw "gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
)

type postSignOutHandler struct {
	shopperauthClient shopperauth.Client
	apigeeClient      apigee.Client
	forterClient      forter.Client
	statsdClient      sw.Client
}

// NewPostSignOutHandler is the handler for the refresh post.
func NewPostSignOutHandler(shopperauthClient shopperauth.Client, apigeeClient apigee.Client, forterClient forter.Client, statsdClient sw.Client) http.Handler {
	return &postSignOutHandler{
		shopperauthClient: shopperauthClient,
		apigeeClient:      apigeeClient,
		forterClient:      forterClient,
		statsdClient:      statsdClient,
	}
}

type SignOutRequest struct {
	AccessToken string `json:"access_token"`
}

// Signout godoc
// @Id Signout
// @Summary Signout the user
// @Description Signout the user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param requestBody body model.SignoutRequest true "Access Token"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {string} string "Success"
// @Failure 401 {string} string "Missing request body"
// @Failure 500 {string} string "Internal Error"
// @Router /v1/signout [post]
func (handler *postSignOutHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	tCtx := apm.FromRequest(request)
	body, err := verifySignOut(writer, request, tCtx)

	if err != nil {
		tCtx.NoticeError(err)
		return
	}

	signoutRequest := SignOutRequest{}
	err = json.Unmarshal(body, &signoutRequest)

	if err != nil {
		tCtx.NoticeError(err)
		return
	}

	stCtx := tCtx.NewGoRoutine()
	concurrency.Run(func() {
		handleSignOut(request, *handler, signoutRequest.AccessToken, stCtx)
	}, tCtx.Context())

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(""))
}

func handleSignOut(request *http.Request, handler postSignOutHandler, token string, tCtx apm.TransactionContext) {
	logger, lc := logging.GetLoggerAndContext(request)

	var signOutErr error
	var shopperID string

	trueClientIP := request.Header.Get("True-Client-Ip")
	installationID := request.Header.Get("X-Nor-Appiid")
	shopperID, signOutErr = getApigeeSignOutResponse(handler, request, token, trueClientIP, installationID, tCtx)

	if signOutErr != nil {
		tCtx.NoticeError(signOutErr)
		handler.statsdClient.Increment("custom", map[string]string{"ApigeeSignOutError": "error"})

		logger.ErrorWithHeader(lc, "SignOut Error", signOutErr, request.Header)
	} else {
		tCtx.NoticeError(signOutErr)
		sendForterLogout(shopperID, installationID, trueClientIP, request, handler, tCtx)
		logger.InfoWithHeader(lc, "Successful signout", request.Header)
	}
}

func getApigeeSignOutResponse(handler postSignOutHandler, request *http.Request, token string, trueClientIP string, installationID string, tCtx apm.TransactionContext) (string, error) {
	logger, lc := logging.GetLoggerAndContext(request)
	user := request.Header.Get("User")
	clientID := request.Header.Get("X-Nor-ClientId")
	deleteTokenRequest := apigee.DeleteTokenRequest{
		ApigeeClient:   handler.apigeeClient,
		AccessToken:    token,
		User:           user,
		TrueClientIP:   trueClientIP,
		InstallationID: installationID,
		ClientID:       clientID,
		TxnCtx:         tCtx,
	}

	signOutErr := deleteTokenRequest.ApigeeClient.DeleteToken(deleteTokenRequest)
	if signOutErr != nil {
		logger.ErrorWithHeader(lc, "Failed to delete apigee token for user: "+user, signOutErr, request.Header)
	}

	shopperID, err := crypto.NewDesShopperIDDecryptorEncryptor().Decrypt(user)
	if err != nil {
		logger.InfoWithHeader(lc, "Failed to decrypt shopper id: "+user, request.Header)
	}
	return shopperID, signOutErr
}

// verifySignOut will validate whether we were sent the correct info and return the body content
func verifySignOut(writer http.ResponseWriter, request *http.Request, tCtx apm.TransactionContext) ([]byte, error) {
	defer tCtx.Segment("verifySignout").End()
	logger, lc := logging.GetLoggerAndContext(request)

	trueClientIP := request.Header.Get("True-Client-Ip")
	installationID := request.Header.Get("X-Nor-Appiid")
	clientID := request.Header.Get("X-Nor-Clientid")
	scope := request.Header.Get("X-Nor-Scope")

	defer request.Body.Close()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		logger.Error(lc, "Error 0013 - verifySignOut", err, map[string]interface{}{"True-Client-Ip": trueClientIP, "appiid": installationID})
		writer.WriteHeader(http.StatusUnauthorized)
		return []byte{}, err
	}
	if clientID == "" {
		logger.InfoWithHeader(lc, "Info - verifySignOut - missing client id", request.Header)
	}
	if scope == "" {
		logger.InfoWithHeader(lc, "Info - verifySignOut - missing scope", request.Header)
	}

	return body, nil
}

func sendForterLogout(shopperID string, installationID string, trueClientIP string, request *http.Request, handler postSignOutHandler, tCtx apm.TransactionContext) {
	logger, lc := logging.GetLoggerAndContext(request)
	if shopperID != "" {
		err := handler.forterClient.SendLogout(shopperID, installationID, trueClientIP, request, tCtx)
		if err != nil {
			logger.Error(lc, "Failed to call Forter for Logout", err, nil)
		}
	} else {
		logger.InfoWithHeader(lc, "Can not call Forter without a shopper id", request.Header)
	}
}
