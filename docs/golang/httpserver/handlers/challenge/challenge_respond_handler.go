package challenge

import (
	"encoding/json"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"net/http"
)

type challengeRespondHandler struct {
	shopperAuthClient shopperauth.Client
}

// NewChallengeSendHandler is the handler for the mfa flows.
func NewChallengeRespondHandler(shopperAuthClient shopperauth.Client) http.Handler {
	return &challengeRespondHandler{
		shopperAuthClient: shopperAuthClient,
	}
}

// Challenge Respond godoc
// @Id ChallengeRespond
// @Summary Calls Customerauth challenge respond
// @Description Calls Customerauth challenge respond that calls eccs challenge respond flow
// @Tags Account
// @Accept json
// @Produce json
// @Param requestBody body model.ChallengeRespondRequest true "User Email"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {string} string "Success"
// @Failure 500 {string} string "Internal Error"
// @Router /v1/challenge/respond [post]
func (handler *challengeRespondHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	getChallengeRespondRequest := func(body []byte) (model.Validatable, error) {
		m := model.ChallengeRespondRequest{}
		err := json.Unmarshal(body, &m)
		return m, err
	}
	callChallengeRespond := func(m interface{}, request *http.Request, tCtx apm.TransactionContext, logger logging.Logger, lc logging.LogContext) (*http.Response, error) {
		return handler.shopperAuthClient.ChallengeRespond(m.(model.ChallengeRespondRequest), request, tCtx, logger, lc)
	}
	challengeServeHTTP(writer, request, getChallengeRespondRequest, callChallengeRespond)
}
