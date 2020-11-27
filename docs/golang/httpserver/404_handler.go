package httpserver

import (
	"net/http"

	"gitlab.nordstrom.com/sentry/authorize/authorizeconstants"

	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/authorize/tracecontext"
)

func NewNotFoundHandler(logger logging.Logger) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.WriteHeader(http.StatusNotFound)

		tc := tracecontext.FromRequest(request)
		logContext := logging.CreateLogContext(tc, "", authorizeconstants.AppName)
		logger.Info(logContext, "404_not_found", map[string]interface{}{"requested_path": request.URL.Path})
	})
}
