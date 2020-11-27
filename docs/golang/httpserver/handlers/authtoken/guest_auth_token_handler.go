package authtoken

import (
	"encoding/base64"
	"encoding/json"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"net/http"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
	"gitlab.nordstrom.com/sentry/authorize/clients/apigee"
	"gitlab.nordstrom.com/sentry/authorize/crypto"
)

type guestAuthTokenHandler struct {
	apigeeClient apigee.Client
}

type CreateAuthTokenResponse struct {
	AccessToken  string
	ExpiresIn    string
	ShopperID    string `json:"ShopperId"`
	WebShopperID string `json:"WebShopperId"`
	TokenType    string
}

func NewGuestAuthTokenHandler(apigeeClient apigee.Client) http.Handler {
	return &guestAuthTokenHandler{apigeeClient}
}

func (h *guestAuthTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger, lc := logging.GetLoggerAndContext(r)
	authHeader := r.Header.Get("Authorization")
	strippedAuthHeader := strings.TrimSpace(strings.Replace(authHeader, "Basic", "", 1))
	decodedAuthHeader, err := base64.StdEncoding.DecodeString(strippedAuthHeader)
	tCtx := apm.FromRequest(r)

	if err != nil {
		tCtx.NoticeError(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	values := strings.Split(string(decodedAuthHeader), ":")
	if len(values) < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	clientID, clientSecret := values[0], values[1]

	shopperID := generateUUID()
	shopperCrypto := crypto.NewDesShopperIDDecryptorEncryptor()
	encryptedShopperID, err := shopperCrypto.Encrypt(shopperID)

	apigeeToken, err := h.apigeeClient.GetToken(apigee.CreateAuthTokenRequest{
		ShopperID:    encryptedShopperID,
		WebShopperID: shopperID,
		Appiid:       r.Header.Get("X-Nor-Appiid"),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TContext:     tCtx,
	})

	if err != nil {
		tCtx.NoticeError(err)
		logger.ErrorWithHeader(lc, "Error:Create guest token", err, r.Header)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authTokenResponse := CreateAuthTokenResponse{
		AccessToken:  apigeeToken.AccessToken,
		ExpiresIn:    apigeeToken.ExpiresIn,
		TokenType:    apigeeToken.TokenType,
		ShopperID:    encryptedShopperID,
		WebShopperID: shopperID,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(authTokenResponse)
}

func generateUUID() string {
	uuid, _ := uuid.NewV4()
	return strings.Replace(uuid.String(), "-", "", -1)
}
