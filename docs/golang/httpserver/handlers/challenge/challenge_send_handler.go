package challenge

import (
	"net/http"

	"encoding/json"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"io/ioutil"
)

type challengeSendHandler struct {
	shopperAuthClient shopperauth.Client
}

// NewChallengeSendHandler is the handler for the mfa flows.
func NewChallengeSendHandler(shopperAuthClient shopperauth.Client) http.Handler {
	return &challengeSendHandler{
		shopperAuthClient: shopperAuthClient,
	}
}

// Challenge Send godoc
// @Id ChallengeSend
// @Summary Calls Customerauth challenge send.
// @Description Calls Customerauth challenge send that calls eccs challenge send. If there is an error from customerauth, the error is thrown to the client.
// @Tags Account
// @Accept json
// @Produce json
// @Param requestBody body model.ChallengeSendRequest true "User Email"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {string} string "Success"
// @Failure 500 {string} string "Internal Error"
// @Router /v1/challenge/send [post]
func (handler *challengeSendHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	getChallengeSendReq := func(body []byte) (model.Validatable, error) {
		m := model.ChallengeSendRequest{}
		err := json.Unmarshal(body, &m)
		return m, err
	}
	callChallengeSend := func(m interface{}, request *http.Request, tCtx apm.TransactionContext, logger logging.Logger, lc logging.LogContext) (*http.Response, error) {
		return handler.shopperAuthClient.ChallengeSend(m.(model.ChallengeSendRequest), request, tCtx, logger, lc)
	}

	challengeServeHTTP(writer, request, getChallengeSendReq, callChallengeSend)
}

func challengeServeHTTP(
	writer http.ResponseWriter,
	request *http.Request,
	jsonFn func([]byte) (model.Validatable, error),
	caFn func(interface{}, *http.Request, apm.TransactionContext, logging.Logger, logging.LogContext) (*http.Response, error),
) {
	logger, lc := logging.GetLoggerAndContext(request)
	logger.Info(lc, "in challengeServeHTTP", nil)
	tCtx := apm.FromRequest(request)

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		tCtx.NoticeError(err)
		logger.ErrorWithHeader(lc, "Failed to read body of "+request.RequestURI, err, request.Header)
		writer.WriteHeader(http.StatusBadRequest)
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

	//Challenge call to custauth
	resp, err := caFn(m, request, tCtx, logger, lc)

	if err != nil {
		tCtx.NoticeError(err)
		logger.ErrorWithHeader(lc, request.RequestURI+" Error", err, request.Header)
	}

	writer.Header().Set("Content-Type", "application/json")
	if resp != nil {
		data := make(map[string]interface{})
		respBody, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err == nil {
			data["responseBody"] = string(respBody)
		}
		data["responseCode"] = resp.StatusCode
		logger.Error(lc, "An error was received from shopper auth", nil, data)
		if resp.StatusCode == 0 {
			resp.StatusCode = http.StatusInternalServerError
			logger.Error(lc, "Zero status for challenge response from shopperAuth: ", err, nil)
		}
		writer.WriteHeader(resp.StatusCode)
	}
}
