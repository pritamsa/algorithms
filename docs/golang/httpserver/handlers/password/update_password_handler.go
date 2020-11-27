package password

import (
	"encoding/json"
	"net/http"

	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"

	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/model"
)

type updatePasswordHandler struct {
	shopperAuthClient shopperauth.Client
	forterClient      forter.Client
	dynamoClient      dynamo.Client
}

// NewUpdatePasswordHandler is the handler for the update password PUT endpoint.
func NewUpdatePasswordHandler(shopperAuthClient shopperauth.Client, dynamoClient dynamo.Client, forterClient forter.Client) http.Handler {
	return &updatePasswordHandler{
		shopperAuthClient: shopperAuthClient,
		forterClient:      forterClient,
		dynamoClient:      dynamoClient,
	}
}

// Update Password godoc
// @Id UpdatePassword
// @Summary Updates a user's password
// @Description Updates a user's password
// @Tags Account
// @Accept json
// @Produce json
// @Param requestBody body model.UpdatePasswordRequest true "Current password and proposed new password"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200
// @Failure 400 {string} string "Returned if the new password is invalid"
// @Failure 401 {string} string "Returned if the old password is incorrect"
// @Failure 404 {string} string "Returned if the email address is not found"
// @Failure 500 {string} string "Internal Error"
// @Router /v1/password/update [put]
func (handler *updatePasswordHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	validateInput := func(body []byte) (model.Validatable, error) {
		m := model.UpdatePasswordRequest{}
		err := json.Unmarshal(body, &m)
		return m, err
	}

	updatePasswordLogic := func(m interface{}, headers http.Header, tCtx apm.TransactionContext, logger logging.Logger, lc logging.LogContext) error {
		return handler.shopperAuthClient.UpdatePassword(m.(model.UpdatePasswordRequest), headers, tCtx, logger, lc)
	}

	passwordServeHTTP(writer, request,
		validateInput,
		updatePasswordLogic,
		handler.dynamoClient)
}
