package password

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
)

type forgotPasswordHandler struct {
	shopperAuthClient shopperauth.Client
	dynamoClient      dynamo.Client
	forterClient      forter.Client
}

// NewForgotPasswordHandler is the handler for the forgot password post.
func NewForgotPasswordHandler(shopperAuthClient shopperauth.Client, dynamoClient dynamo.Client, forterClient forter.Client) http.Handler {
	return &forgotPasswordHandler{
		shopperAuthClient: shopperAuthClient,
		dynamoClient:      dynamoClient,
		forterClient:      forterClient,
	}
}

// Forgot Password godoc
// @Id ForgotPassword
// @Summary Sends code to user to reset their password
// @Description Sends an email to the user with a code the user enters to reset their password
// @Tags Account
// @Accept json
// @Produce json
// @Param requestBody body model.ForgotPasswordRequest true "User Email"
// @Param code query string true "Server Cryptographic Nonce Code"
// @Param verifier query string true "Verifier for Nonce Code"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "User not found"
// @Failure 410 {string} string "Gone returned if user's account is disabled"
// @Failure 418 {string} string "Returned if the caller uses same code again for same Appiid"
// @Failure 429 {string} string "Too Many Requests for a single email"
// @Failure 500 {string} string "Internal Error"
// @Router /v1/password/forgot [post]
func (handler *forgotPasswordHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	passwordServeHTTP(writer, request,
		func(body []byte) (model.Validatable, error) {
			m := model.ForgotPasswordRequest{}
			err := json.Unmarshal(body, &m)
			return m, err
		},
		func(m interface{}, headers http.Header, tCtx apm.TransactionContext, logger logging.Logger, lc logging.LogContext) error {
			return handler.shopperAuthClient.ForgotPassword(m.(model.ForgotPasswordRequest), headers, tCtx, logger, lc)
		},
		handler.dynamoClient)
}

func passwordServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
	jsonFn func([]byte) (model.Validatable, error),
	caFn func(interface{}, http.Header, apm.TransactionContext, logging.Logger, logging.LogContext) error,
	dynamoClient dynamo.Client,
) {
	logger, lc := logging.GetLoggerAndContext(request)
	logger.Info(lc, "in passwordServeHTTP", nil)
	if handlers.VerifyFirstAuth(writer, request, dynamoClient, http.StatusTeapot) {
		return
	}
	tCtx := apm.FromRequest(request)

	trueClientIP := request.Header.Get("True-Client-Ip")
	installationID := request.Header.Get("X-Nor-Appiid")

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		tCtx.NoticeError(err)
		logger.ErrorWithHeader(lc, "Failed to read body of "+request.RequestURI, err, request.Header)
		writer.WriteHeader(http.StatusTeapot)
		return
	}

	m, err := jsonFn(requestBody)
	if err != nil {
		tCtx.NoticeError(err)
		logger.ErrorWithHeader(lc, "Failed to unmarshal body of "+request.RequestURI, err, request.Header)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	validationErr := m.Validate()
	if validationErr != nil {
		tCtx.NoticeError(validationErr)
		logger.Error(lc, "Validation error: "+validationErr.Error(), validationErr, nil)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(validationErr.Error()))
		return
	}

	if derr := dynamoClient.DeleteAuth(installationID, trueClientIP, tCtx); derr != nil {
		tCtx.NoticeError(derr)
		logger.Error(lc, "DeleteFromDynamo", derr, nil)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = caFn(m, request.Header, tCtx, logger, lc)
	if err != nil {
		tCtx.NoticeError(err)
		logger.ErrorWithHeader(lc, request.RequestURI+" Error", err, request.Header)
		errorCode, parseErr := strconv.Atoi(err.Error())

		if parseErr != nil {
			logger.Error(lc, "Error parsing status code to int: "+err.Error(), parseErr, nil)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if errorCode > 0 {
			writer.WriteHeader(errorCode)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
}
