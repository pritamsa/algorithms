package gohealthcheck

import (
	"net/http"
	"os"

	"gitlab.nordstrom.com/sentry/gologger"
)

func RoutesForHealthchecks(checks ...Healthcheckable) Routes {
	logger := gologger.NewLogger(os.Stdout, "", "")
	return RoutesForHealthchecksWithLogger(logger, checks...)
}

func RoutesForHealthchecksWithLogger(logger gologger.Logger, checks ...Healthcheckable) Routes {
	return Routes{
		Basic: HealthcheckRoute{
			Name:    basicHandlerKey,
			Path:    basicHandlerPath,
			Handler: handlerForHealthchecks(logger),
		},
		Advanced: HealthcheckRoute{
			Name:    advancedHandlerKey,
			Path:    advancedHandlerPath,
			Handler: handlerForHealthchecks(logger, checks...),
		},
	}
}

type Routes struct {
	Basic    HealthcheckRoute
	Advanced HealthcheckRoute
}

type HealthcheckRoute struct {
	Handler http.Handler
	Name    string
	Path    string
}

const (
	advancedHandlerKey  = "AdvancedHealthcheck"
	advancedHandlerPath = "/status/advanced"

	basicHandlerKey  = "BasicHealthcheck"
	basicHandlerPath = "/status/basic"
)
