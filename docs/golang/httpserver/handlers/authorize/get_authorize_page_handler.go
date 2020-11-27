package authorize

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"gitlab.nordstrom.com/sentry/authcrypto"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/try"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
)

type getAuthorizePageHandler struct {
	dynamoClient      dynamo.Client
	statsdClient      statsd_wrapper.Client
	serverEnvironment string
}

// NewGetAuthorizePageHandler is the handler for the authorize page.
func NewGetAuthorizePageHandler(dynamoClient dynamo.Client, statsdClient statsd_wrapper.Client, serverEnvironment string) http.Handler {
	return &getAuthorizePageHandler{
		dynamoClient:      dynamoClient,
		statsdClient:      statsdClient,
		serverEnvironment: serverEnvironment,
	}
}

// Authorize Get godoc
// @Id AuthorizeGet
// @Summary Get Nonce Code From Server
// @Description Leg 2 of 3 legged Authorization calls. This call gets the  server side cryptographic nonce code.
// @Tags Authentication
// @Produce json
// @Param code query string true "Server Cryptographic Nonce Code"
// @Param verifier query string true "Client side unsigned cryptographic nonce code used for generating code in AuthInit"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {string} string "{"code": "9nlgKYYsx..."}"
// @Failure 403 {string} string "This is returned if invalid parameters are passed"
// @Failure 500 {string} string "Internal Error"
// @Router /v2/authorize [get]
func (handler *getAuthorizePageHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logger, lc := logging.GetLoggerAndContext(request)
	tCtx := apm.FromRequest(request)
	if handlers.VerifyFirstAuth(writer, request, handler.dynamoClient, http.StatusForbidden) {
		return
	}

	trueClientIP := request.Header.Get("True-Client-Ip")
	installationID := request.Header.Get("X-Nor-Appiid")

	key, _ := authcrypto.NewSigningKey()
	_, signature, _ := authcrypto.GenerateAuthorization(key)
	submitcode := base64.URLEncoding.EncodeToString(signature)

	uerr := handler.dynamoClient.UpdateAuth(installationID, try.GetPkceIP(trueClientIP), map[string]string{
		"SubmitCode":       submitcode,
		"AuthGetTimestamp": time.Now().Format(time.RFC3339Nano),
	}, tCtx)
	if uerr != nil {
		tCtx.NoticeError(uerr)
		logger.Error(lc, "failed to set submitcode & record timestamp for authorize/get: "+submitcode, uerr, nil)
		http.Redirect(writer, request, "/v1/authinit", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write([]byte(fmt.Sprintf(`{"code":"%s"}`, submitcode)))

}
