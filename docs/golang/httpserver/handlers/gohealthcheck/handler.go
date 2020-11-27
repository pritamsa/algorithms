package gohealthcheck

import (
	"net/http"

	"gitlab.nordstrom.com/sentry/gologger"
)

func handlerForHealthchecks(logger gologger.Logger, checks ...Healthcheckable) http.Handler {
	return NewHealthCheckHandler(logger, checks)
}
