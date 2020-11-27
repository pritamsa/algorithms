package verify

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.nordstrom.com/sentry/authorize/clients/shopperauth"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"gitlab.nordstrom.com/sentry/authcrypto"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/tracecontext"
)

type AssociateRequest struct {
	WebShopperID string `json:"webShopperId"`
	DeviceID     string `json:"deviceId"`
}

var sess, _ = session.NewSession(&aws.Config{Region: aws.String("us-west-2")})

const (
	XNorScope    = "X-Nor-Scope"
	UserAgent    = "User-Agent"
	TrueClientIp = "True-Client-Ip"
	XAkamaiEdge  = "X-Akamai-Edgescape"
	XNorClientId = "X-Nor-Clientid"
)

var CommonHeaders = []string{XNorScope, UserAgent, TrueClientIp, XAkamaiEdge, XNorClientId}

// IssueRequest makes a RESTful call to a given endpoint in the Verify service.
// reqBodyContent: Struct to be sent as request body (JSON)
// defaultReturn: Value to return if any errors occur
// method: HTTP Method
// path: URL endpoint
// endpointName: Name of endpoint for logging purposes
// expectedStatus: Expected HTTP Response status code for successful request
// responseHandler: Performs further validation on response body. Pass nil if no further validation is required
// request: Authorize HTTP request
// headers: headers to use in the request
// IssueRequest returns response value and error if it exists.
func (c *client) IssueRequest(reqBodyContent interface{}, method string, path string, endpointName string, expectedStatus int, originalRequestHeader http.Header, headers map[string]string, tCtx apm.TransactionContext) ([]byte, error) {
	var (
		reqBody []byte
		err     error
	)
	if reqBodyContent != nil {
		reqBody, err = json.Marshal(reqBodyContent)
	}
	logger, lc := logging.GetSingleLogger(), logging.CreateLogContext(tCtx.Context(), "", "")
	if err != nil {
		logger.ErrorWithHeader(lc, fmt.Sprintf("%s request creation error", endpointName), err, originalRequestHeader)
		return nil, err
	}

	var req *http.Request

	req, err = http.NewRequest(method, fmt.Sprintf("%s/%s", c.args.BaseURL, path), bytes.NewReader(reqBody))
	if err != nil {
		logger.ErrorWithHeader(lc, fmt.Sprintf("%s request creation error", endpointName), err, originalRequestHeader)
		return nil, err
	}

	tc := lc.TraceContext
	ctx := tracecontext.NewContext(context.Background(), tc)

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req = tCtx.PrepareRequest(req)
	resp, err := c.httpClient.Do(ctx, req)
	if err != nil {
		logger.ErrorWithHeader(lc, fmt.Sprintf("%s request do error", endpointName), err, originalRequestHeader)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorWithHeader(lc, fmt.Sprintf("%s Failed to read response", endpointName), err, originalRequestHeader)
		return nil, err
	}

	if resp.StatusCode != expectedStatus {
		errMessage := fmt.Sprintf("%s Non-%d Response %v;\n", endpointName, expectedStatus, resp.StatusCode)
		err = errors.New(fmt.Sprintf("%s ResponseBody: %s", errMessage, body))
		logger.ErrorWithHeader(lc, fmt.Sprintf("%s logging response body", endpointName), err, originalRequestHeader)
		return nil, err
	}

	return body, nil
}

func (c *client) ChallengeInit(challengeInitRequest shopperauth.ChallengeInitRequest, originalRequestHeader http.Header, tCtx apm.TransactionContext) (string, error) {
	defer tCtx.Segment("Verify: Challenge Init").End()
	nonce := generateNonce()

	challengeInitRequest.SessionID = nonce

	logger, lc := logging.GetSingleLogger(), logging.CreateLogContext(tCtx.Context(), "", "")

	// Load common headers
	headers := make(map[string]string)
	for _, item := range CommonHeaders {
		headers[item] = originalRequestHeader.Get(item)
	}

	body, err := c.IssueRequest(challengeInitRequest, http.MethodPost, "challenge/init", "Challenge Init", http.StatusOK, originalRequestHeader, headers, tCtx)
	if err != nil {
		return "", err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(body, &m)
	jsonBody := string(body)
	if err != nil {
		logger.ErrorWithHeader(lc, "Challenge Init response body error", err, originalRequestHeader)
		return "", err
	}
	sessionID := m["sessionId"].(string)
	if nonce != sessionID {
		errMessage := fmt.Sprintf("Nonce value %s != session id %s", nonce, sessionID)
		err := errors.New(fmt.Sprintf("%s ResponseBody: %s", errMessage, jsonBody))
		logger.ErrorWithHeader(lc, fmt.Sprintf("Challenge Init sessionId:%s != nonce: %s", sessionID, nonce), err, originalRequestHeader)
		return "", err
	}
	return jsonBody, nil
}

func generateNonce() string {
	key, _ := authcrypto.NewSigningKey()
	_, signature, _ := authcrypto.GenerateAuthorization(key)
	nonce := base64.URLEncoding.EncodeToString(signature)
	return strings.Replace(nonce, "==", "", -1)
}
