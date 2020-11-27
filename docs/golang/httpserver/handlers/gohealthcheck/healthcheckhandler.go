package gohealthcheck

import (
	"encoding/json"
	"fmt"
	"gitlab.nordstrom.com/sentry/gologger"
	"gitlab.nordstrom.com/sentry/gologger/tracecontext"
	"net/http"
	"sync"
)

type ShutdownHandler interface {
	Shutdown()
}

type healthCheckHandler struct {
	shutdownReceived sync.RWMutex
	serviceDown      bool
	logger           gologger.Logger
	checks           []Healthcheckable
}

func (h *healthCheckHandler) normalServe(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", "application/json")
	healthy := []string{}
	unhealthy := []string{}
	for _, check := range h.checks {
		if ok, err := check.IsHealthy(); ok {
			healthy = append(healthy, check.Name())
		} else {
			h.logger.Error(
				tracecontext.GetFrom(request),
				fmt.Sprintf("%s.health-check-failed", check.Name()),
				err,
			)
			unhealthy = append(unhealthy, check.Name())
		}
	}

	healthCheck := map[string][]string{
		"healthy":   healthy,
		"unhealthy": unhealthy,
	}
	if len(unhealthy) > 0 {
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}

	encoder := json.NewEncoder(responseWriter)
	if err := encoder.Encode(healthCheck); err != nil {
		h.logger.Error(tracecontext.GetFrom(request), "Health check JSON encoding error", err, nil)
		responseWriter.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *healthCheckHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.shutdownReceived.RLock()
	if h.serviceDown {
		responseWriter.WriteHeader(http.StatusServiceUnavailable)
	} else {
		h.normalServe(responseWriter, request)
	}
	h.shutdownReceived.RUnlock()
}

func (h *healthCheckHandler) Shutdown() {
	h.shutdownReceived.Lock()
	h.serviceDown = true
	h.shutdownReceived.Unlock()
}

func NewHealthCheckHandler(logger gologger.Logger, checks []Healthcheckable) http.Handler {
	return &healthCheckHandler{
		sync.RWMutex{},
		false,
		logger,
		checks,
	}
}
