package authorize

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/crypto"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/login"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
)

type postAuthorizeHandler struct {
	statsdClient statsd_wrapper.Client
	dynamoClient dynamo.Client
	loginManager login.Manager
	encryptor    *crypto.Encryptor
}

// NewPostAuthorizeHandler is the handler for the authorize post.
func NewPostAuthorizeHandler(
	statsdClient statsd_wrapper.Client,
	dynamoClient dynamo.Client,
	loginManager login.Manager,
	encryptor *crypto.Encryptor,
) http.Handler {
	return &postAuthorizeHandler{
		statsdClient: statsdClient,
		dynamoClient: dynamoClient,
		loginManager: loginManager,
		encryptor:    encryptor,
	}
}

// Authorize Post godoc
// @Id AuthorizePost
// @Summary Login with credentials and receive access token
// @Description Leg 3 of 3 legged Authorization calls. Submit user credentials.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param requestBody body model.AuthorizationRequest true "Log-in Credentials"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {object} shopperauth.SigninResponse
// @Failure 250 {string} string "The user has been MFA challenged and must complete the challenge before they can login"
// @Failure 401 {string} string "Username or password was invalid"
// @Failure 403 {string} string "Locked Out. Password Reset email will be sent, Username/password invalid, invalid PKCE params"
// @Failure 500 {string} string "Internal Error"
// @Router /v2/authorize [post]
func (h *postAuthorizeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	logger, lc := logging.GetLoggerAndContext(request)
	requestBody, err := handlers.VerifyAuth(writer, request, h.statsdClient, h.dynamoClient)
	if err != nil {
		return
	}

	model, err := parseBody(requestBody)
	if err != nil {
		tCtx := apm.FromRequest(request)
		tCtx.NoticeError(err)
		logger.Error(lc, "postAuthorizeHandler", err, nil)
		http.Redirect(writer, request, "/v1/authinit", http.StatusUnauthorized)
		return
	}

	username := model.Username
	password := model.Password

	username = strings.Trim(username, " ")
	if handlers.IsEmailInvalid(username) {
		enryptedUsername := h.encryptor.AESEncrypt(username, lc.TraceContext)
		tCtx := apm.FromRequest(request)
		tCtx.AddAttribute("username", enryptedUsername)
		logger.InfoWithHeader(lc, "Invalid username: ["+enryptedUsername+"]", request.Header)
		http.Redirect(writer, request, "/v1/authinit", http.StatusBadRequest)
		return
	}

	login := login.Request{
		Username:        username,
		Password:        password,
		PersistentOptIn: model.PersistentOptIn,
	}

	h.loginManager.Login(writer, request, login)
}

func parseBody(body []byte) (model.AuthorizationRequest, error) {
	var m model.AuthorizationRequest
	err := json.Unmarshal(body, &m)
	if err != nil {
		return model.AuthorizationRequest{}, err
	}

	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(body, &objmap)
	if err != nil {
		return model.AuthorizationRequest{}, err
	}

	code := m.Code

	m.Username, err = mapValue(objmap, "username__"+code)
	if err != nil {
		return model.AuthorizationRequest{}, err
	}

	m.Password, err = mapValue(objmap, "password__"+code)
	if err != nil {
		return model.AuthorizationRequest{}, err
	}

	return m, nil
}

func mapValue(objmap map[string]*json.RawMessage, key string) (string, error) {
	if val, ok := objmap[key]; ok {
		var result string
		err := json.Unmarshal(*val, &result)
		if err != nil {
			return "", err
		}
		return result, nil
	}
	return "", errors.New("Key " + key + " doesn't exist")
}
