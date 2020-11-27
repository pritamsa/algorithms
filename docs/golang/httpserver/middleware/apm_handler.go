package apm_middleware

import (
	"context"
	"github.com/newrelic/go-agent"
	"gitlab.nordstrom.com/sentry/authorize/clients/apm"
	"gitlab.nordstrom.com/sentry/authorize/logging"
	"net/http"
)

type apmHandler struct {
	handler http.Handler
	client  newrelic.Application
	logger  logging.Logger
}

func NewApmHandler(
	handler http.Handler,
	client newrelic.Application,
	logger logging.Logger,
) http.Handler {
	return &apmHandler{
		handler: handler,
		client:  client,
		logger:  logger,
	}
}

func (h *apmHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tCtx := apm.StartWebTransaction(r, w, h.client, nil, h.logger)
	writer := writerWrapper(w, tCtx)
	ctx := context.WithValue(r.Context(), "tCtx", tCtx)
	h.handler.ServeHTTP(writer, r.WithContext(ctx))
	tCtx.End()
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	tCtx apm.TransactionContext
}

func (w *wrappedResponseWriter) Write(data []byte) (int, error) {
	return w.tCtx.Write(data)
}

func (w *wrappedResponseWriter) WriteHeader(code int) {
	w.tCtx.WriteHeader(code)
}

func writerWrapper(w http.ResponseWriter, tCtx apm.TransactionContext) *wrappedResponseWriter {
	return &wrappedResponseWriter{
		ResponseWriter: w,
		tCtx:           tCtx,
	}
}
