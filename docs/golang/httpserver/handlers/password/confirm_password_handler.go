package password

import (
	"encoding/json"
	"net/http"

	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/login"
	"gitlab.nordstrom.com/sentry/authorize/model"
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
// @Id ConfirmPassword
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
// @Router /v1/password/confirm [post]
func (handler *confirmPasswordHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	passwordServeHTTP(writer, request,
		func(body []byte) (model.Validatable, error) {
			if handler.dynamoClient.GetFeatureFlag("useNewConfirmPassword", apm.FromRequest(request)) {
				m := model.ConfirmPasswordV2Request{}
				err := json.Unmarshal(body, &m)
				return m, err
			}
			m := model.ConfirmPasswordRequest{}
			err := json.Unmarshal(body, &m)
			return m, err

		},
		func(m interface{}, headers http.Header, tCtx apm.TransactionContext, logger logging.Logger, lc logging.LogContext) error {

			var email, username, password string
			var err error

			if handler.dynamoClient.GetFeatureFlag("useNewConfirmPassword", tCtx) {
				//Get email in the response from custauth to log the user in
				email, err = handler.shopperAuthClient.ConfirmPasswordV2(m.(model.ConfirmPasswordV2Request), headers, tCtx, logger, lc)
				password = m.(model.ConfirmPasswordV2Request).Password
			} else {
				email, err = handler.shopperAuthClient.ConfirmPassword(m.(model.ConfirmPasswordRequest), headers, tCtx, logger, lc)
				if m.(model.ConfirmPasswordRequest).Email == "" {
					username = email
				} else {
					username = m.(model.ConfirmPasswordRequest).Email
				}
				password = m.(model.ConfirmPasswordRequest).Password
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
