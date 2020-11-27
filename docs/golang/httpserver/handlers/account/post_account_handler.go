package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nu7hatch/gouuid"
	"gitlab.nordstrom.com/sentry/authorize/authorizeconstants"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/clients/doomhammer"
	"gitlab.nordstrom.com/sentry/authorize/clients/dynamo"
	"gitlab.nordstrom.com/sentry/authorize/clients/forter"
	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"
	"gitlab.nordstrom.com/sentry/authorize/concurrency"
	"gitlab.nordstrom.com/sentry/authorize/crypto"
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/login"
	"gitlab.nordstrom.com/sentry/authorize/model"
	"gitlab.nordstrom.com/sentry/gohttp/statsd_wrapper"
)

const (
	featureFlagName = "create-account-in-ICON"
)

type postAccountHandler struct {
	doomHammerClient  doomhammer.Client
	shopperAuthClient shopperauth.Client
	forterClient      forter.Client
	statsdClient      statsd_wrapper.Client
	loginManager      login.Manager
	dynamoClient      dynamo.Client
	encryptor         *crypto.Encryptor
}

// NewPostAccountHandler is the handler for the account post.
func NewPostAccountHandler(
	doomHammerClient doomhammer.Client,
	shopperAuthClient shopperauth.Client,
	forterClient forter.Client,
	statsdClient statsd_wrapper.Client,
	loginManager login.Manager,
	dynamoClient dynamo.Client,
	encryptor *crypto.Encryptor,
) http.Handler {
	return &postAccountHandler{
		doomHammerClient:  doomHammerClient,
		shopperAuthClient: shopperAuthClient,
		forterClient:      forterClient,
		statsdClient:      statsdClient,
		loginManager:      loginManager,
		dynamoClient:      dynamoClient,
		encryptor:         encryptor,
	}
}

// Create Account godoc
// @Id CreateAccount
// @Summary Creates a new account
// @Description Creates a new account for a user
// @Tags Account
// @Accept json
// @Produce json
// @Param requestBody body model.AccountRequest true "Create Account Parameters"
// @Param X-Nor-Appiid header string true "App installation id for mobile"
// @Param X-Nor-Clientid header string true "Unique client id"
// @Param X-Nor-Scope header string true "WebRegistered or MobileRegistered"
// @Param Tracecontext header string false "Unique trace context id (UUID) for the request"
// @Success 200 {object} shopperauth.SigninResponse
// @Failure 250 {string} string "The user has been MFA challenged and must complete the challenge before they can login"
// @Failure 500 {string} string "Internal Error"
// @Failure 400 {string} string "Bad Request"
// @Router /v2/account [post]
func (handler *postAccountHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	logger, lc := logging.GetLoggerAndContext(request)
	tCtx := apm.FromRequest(request)
	tc := lc.TraceContext

	requestBody, err := handlers.VerifyAuth(writer, request, handler.statsdClient, handler.dynamoClient)
	if err != nil {
		tCtx.NoticeError(err)
		return
	}

	account := model.AccountEntity{}
	err = json.Unmarshal(requestBody, &account)
	if err != nil {
		tCtx.NoticeError(err)
		logger.ErrorWithHeader(lc, "", err, request.Header)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf(authorizeconstants.JSONErrorMessageTemplate, "Bad Request")))
		return
	}

	//check header before calling validate
	headers := request.Header
	sc := headers.Get("X-Nor-Scope")

	var validationErr error

	if strings.Contains(sc, "WebRegistered") {
		validationErr = account.ValidateWithEmailPassword()
	} else {
		validationErr = account.Validate()
	}

	if validationErr != nil {
		tCtx.NoticeError(validationErr)
		logger.ErrorWithHeader(lc, "", validationErr, request.Header)
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(fmt.Sprintf(authorizeconstants.JSONErrorMessageTemplate, validationErr.Error())))
		return
	}

	trimmedEmail := strings.Trim(account.Email, " ")
	if handlers.IsEmailInvalid(trimmedEmail) {
		encryptedEmail := handler.encryptor.AESEncrypt(trimmedEmail, tc)
		logger.InfoWithHeader(lc, "Invalid username: ["+encryptedEmail+"]", request.Header)
		http.Redirect(writer, request, "/v1/authinit", http.StatusBadRequest)
		return
	}

	shopperID := ""

	uuid, err := uuid.NewV4()
	generatedShopperID := strings.ToUpper(strings.Replace(uuid.String(), "-", "", -1))

	if err != nil {
		logger.Error(lc, "Error creating shopper id", err, nil)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = handler.shopperAuthClient.CreateUser(shopperauth.CreateUserRequest{
		Username: generatedShopperID, // in IDP, Username is the shopperId
		Password: account.Password,
		Email:    trimmedEmail,
	}, tCtx)
	if err != nil {
		switch err {
		case authorizeconstants.ErrEmailAlreadyExists,
			authorizeconstants.ErrInvalidEmailOrPassword,
			authorizeconstants.ErrGenericBadRequest:
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(fmt.Sprintf(authorizeconstants.JSONErrorMessageTemplate, err.Error())))
			return
		default:
			logger.Error(lc, "CustAuthError", err, nil)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	_, err = handler.doomHammerClient.CreateShopper(doomhammer.CreateShopperRequest{
		Email:           trimmedEmail,
		FirstName:       account.FirstName,
		LastName:        account.LastName,
		MobileNumber:    account.MobileNumber,
		SubscribeEmail:  account.IsOptIn,
		SubscribeMobile: account.IsOptIn,
	},
		request.Header,
		tCtx,
		generatedShopperID,
	)
	if err != nil {
		tCtx.NoticeError(err)
		satCtx := tCtx.NewGoRoutine()
		concurrency.Run(func() {
			// disable and then delete account if we can't create one in both cognito and doomhammer
			disableErr := handler.shopperAuthClient.DisableUser(generatedShopperID, satCtx)
			if disableErr != nil {
				satCtx.NoticeError(err)
				logger.Error(lc, "Error disabling user in cognito during account creation rollback", disableErr, nil)
			}
			deleteErr := handler.shopperAuthClient.DeleteUser(trimmedEmail, satCtx)
			if deleteErr != nil {
				satCtx.NoticeError(err)
				logger.Error(lc, "Error deleting user in cognito during account creation rollback", deleteErr, nil)
			}
		}, tc)

		switch err {
		case authorizeconstants.ErrEmailAlreadyExists,
			authorizeconstants.ErrInvalidEmailFormat,
			authorizeconstants.ErrMobileNumberAlreadyExists,
			authorizeconstants.ErrInvalidMobileNumberFormat,
			authorizeconstants.ErrInvalidFirstNameOrLastName,
			authorizeconstants.ErrGenericBadRequest:
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(fmt.Sprintf(authorizeconstants.JSONErrorMessageTemplate, err.Error())))
			return
		default:
			logger.Error(lc, "DoomHammerError", err, nil)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	shopperID = generatedShopperID
	trueClientIP := request.Header.Get("True-Client-Ip")
	installationID := request.Header.Get("X-Nor-Appiid")

	// TODO(qzaj): move this to before account creation
	ftCtx := tCtx.NewGoRoutine()
	concurrency.Run(func() {
		handler.forterClient.GetForterAccountSignupResult(forter.ClientModel{
			ShopperID:      shopperID,
			InstallationID: installationID,
			TrueClientIP:   trueClientIP,
			Username:       trimmedEmail,
			FirstName:      account.FirstName,
			LastName:       account.LastName,
			HTTPRequest:    request,
		}, ftCtx)
	}, tc)

	concurrency.Run(func() {
		if perr := handler.dynamoClient.Put(dynamo.ShopperIdTable, login.ShopperIDCache{
			Email:     trimmedEmail,
			ShopperID: shopperID,
		}, tCtx); perr != nil {
			tCtx.NoticeError(perr)
			logger.Error(lc, "Put shopperCache", perr, nil)
		}
	}, tc)

	login := login.Request{
		Username: trimmedEmail,
		Password: account.Password,
	}
	handler.loginManager.Login(writer, request, login)
}
