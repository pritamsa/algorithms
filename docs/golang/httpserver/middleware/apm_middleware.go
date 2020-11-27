package apm_middleware

import (
	"github.com/newrelic/go-agent"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"gitlab.nordstrom.com/sentry/gohttp/middleware"
	"net/http"
)

type apmMiddleware struct {
	client newrelic.Application
	logger logging.Logger
}

func NewApmMiddleware(client newrelic.Application, logger logging.Logger) middleware.Middleware {
	return apmMiddleware{
		client: client,
		logger: logger,
	}
}

func (m apmMiddleware) Wrap(handler http.Handler) http.Handler {
	return NewApmHandler(handler, m.client, m.logger)
}
