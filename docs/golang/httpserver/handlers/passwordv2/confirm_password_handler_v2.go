package passwordv2

import (
	"encoding/json"
	"net/http"

	"errors"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/login"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"io/ioutil"
	"strconv"
)

type confirmPasswordHandler struct {
	shopperAuthClient shopperauth.Client
	dynamoClient      dynamo.Client
	forterClient      forter.Client
	loginManager      login.Manager
}

// NewConfirmPasswordHandler is the handler for the forgot password post.
func NewConfirmPasswordHandler(shopperAuthClient shopperauth.Client, dynamoClient dynamo.Client, forterClient forter.Client, loginManager login.Manager) http.Handler {
	return &confirmPasswordHandler{
		shopperAuthClient: shopperAuthClient,
		dynamoClient:      dynamoClient,
		forterClient:      forterClient,
		loginManager:      loginManager,
	}
}

// Confirm Password godoc
// @Id ConfirmPasswordv2
// @Summary Validates the code sent to the user to reset their password
// @Description Validates the code sent to the user to reset their password
// @Tags Account
// @Accept json
// @Produce json
// @Param requestBody body model.ConfirmPasswordRequest true "Reset password code, email, and new password"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {object} shopperauth.SigninResponse
// @Failure 400 {string} string "Returned if the email/shopperId and/or new password is invalid, or wrong 6 digit code entered"
// @Failure 429 {string} string "Too Many Requests for a single email"
// @Failure 500 {string} string "Internal Error"
// @Router /v2/password/confirm [post]
func (handler *confirmPasswordHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	passwordServeHTTP(writer, request,
		func(body []byte) (model.Validatable, error) {
			if handler.dynamoClient.GetFeatureFlag("useNewConfirmPassword", apm.FromRequest(request)) {
				m := model.ConfirmPasswordV2Request{}
				err := json.Unmarshal(body, &m)
				return m, err
			}
			return nil, errors.New("Enable useNewConfirmPassword feature flag to make confirm password work.")

		},
		func(m interface{}, headers http.Header, tCtx apm.TransactionContext, logger logging.Logger, lc logging.LogContext) error {

			var username, password string
			var err error

			if handler.dynamoClient.GetFeatureFlag("useNewConfirmPassword", tCtx) {
				//Get email in the response from custauth to log the user in
				username, err = handler.shopperAuthClient.ConfirmPasswordV2(m.(model.ConfirmPasswordV2Request), headers, tCtx, logger, lc)
				password = m.(model.ConfirmPasswordV2Request).Password
			}

			if err != nil {
				return err
			}

			loginRequest := login.Request{
				Username:             username,
				Password:             password,
				LoginFromConfirmFlow: true,
			}

			handler.loginManager.Login(writer, request, loginRequest)
			return nil
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
