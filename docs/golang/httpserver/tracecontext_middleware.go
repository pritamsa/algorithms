package httpserver

import (
	"net/http"

	"gitlab.nordstrom.com/sentry/authorize/tracecontext"
	"gitlab.nordstrom.com/sentry/gohttp/middleware"
)

type traceContextMiddleware struct{}

//NewTraceContextMiddleware returns a new trace context middleware
func NewTraceContextMiddleware() middleware.Middleware {
	return &traceContextMiddleware{}
}

func (m traceContextMiddleware) Wrap(h http.Handler) http.Handler {
	return NewTraceContextHandler(h)
}

type traceContextHandler struct {
	handler http.Handler
}

//NewTraceContextHandler returns a new trace context handler with the provided inner handler
func NewTraceContextHandler(h http.Handler) traceContextHandler {
	return traceContextHandler{h}
}

func (h traceContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tc := tracecontext.FromRequest(r)
	tc.Set(r)
	h.handler.ServeHTTP(w, r)
}
