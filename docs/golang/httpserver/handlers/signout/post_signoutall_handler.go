package signout

import (
	"encoding/json"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/concurrency"
	"gitlab.nordstrom.com/sentry/authorize/crypto"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	sw "gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
	"net/http"
)

type postSignOutAllHandler struct {
	postSignOutHandler
}

func NewPostSignOutAllHandler(shopperauthClient shopperauth.Client, apigeeClient apigee.Client, forterClient forter.Client, statsdClient sw.Client) http.Handler {
	return &postSignOutHandler{
		shopperauthClient: shopperauthClient,
		apigeeClient:      apigeeClient,
		forterClient:      forterClient,
		statsdClient:      statsdClient,
	}
}

// SignoutAll godoc
// @Id SignoutAll
// @Summary Signout all the user
// @Description Signout all the user
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
// @Router /v1/signout/all [post]
func (handler *postSignOutAllHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

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
		handleSignOutAll(request, *handler, signoutRequest.AccessToken, stCtx)
	}, tCtx.Context())

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(""))
}

func handleSignOutAll(request *http.Request, handler postSignOutAllHandler, token string, tCtx apm.TransactionContext) {
	logger, lc := logging.GetLoggerAndContext(request)

	var signOutErr error
	var shopperID string

	trueClientIP := request.Header.Get("True-Client-Ip")
	installationID := request.Header.Get("X-Nor-Appiid")

	shopperID, signOutErr = getApigeeSignOutAllResponse(handler, request, tCtx)

	if signOutErr != nil {
		tCtx.NoticeError(signOutErr)
		handler.statsdClient.Increment("custom", map[string]string{"ApigeeSignOutAllError": "error"})

		logger.ErrorWithHeader(lc, "SignOut all Error", signOutErr, request.Header)
	} else {
		tCtx.NoticeError(signOutErr)
		sendForterLogout(shopperID, installationID, trueClientIP, request, handler.postSignOutHandler, tCtx)
		logger.InfoWithHeader(lc, "Successful signout", request.Header)
	}
}

func getApigeeSignOutAllResponse(handler postSignOutAllHandler, request *http.Request, tCtx apm.TransactionContext) (string, error) {
	logger, lc := logging.GetLoggerAndContext(request)
	user := request.Header.Get("User")
	deleteAllTokenRequest := apigee.DeleteAllTokensRequest{
		ShopperID: user,
		TxnCtx:    tCtx,
	}

	signOutErr := deleteAllTokenRequest.ApigeeClient.DeleteAllTokens(deleteAllTokenRequest)
	if signOutErr != nil {
		logger.ErrorWithHeader(lc, "Failed to delete  apigee all tokens for user: "+user, signOutErr, request.Header)
	}

	shopperID, err := crypto.NewDesShopperIDDecryptorEncryptor().Decrypt(user)
	if err != nil {
		logger.InfoWithHeader(lc, "Failed to decrypt shopper id: "+user, request.Header)
	}
	return shopperID, signOutErr
}
