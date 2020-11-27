package httpserver

import (
	"net/http"
	"time"

	"gitlab.nordstrom.com/sentry/gohttp/middleware"
)

type timeoutMiddleware struct {
	timeout time.Duration
}

func NewTimeoutMiddleware(timeout time.Duration) middleware.Middleware {
	return &timeoutMiddleware{
		timeout: timeout,
	}
}

func (m timeoutMiddleware) Wrap(handler http.Handler) http.Handler {
	return http.TimeoutHandler(handler, m.timeout, RequestCausedATimeoutMessage)
}

const RequestCausedATimeoutMessage = "Your request took too long waiting for the server to respond"
