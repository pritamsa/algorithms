package router

import (
	"gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/gohealthcheck"
	"gitlab.nordstrom.com/sentry/gohttp/middleware"
)

//go:generate counterfeiter . HealthcheckRouteBuilder
type HealthcheckRouteBuilder interface {
	RoutesForHealthchecks(...gohealthcheck.Healthcheckable) gohealthcheck.Routes
}

func NewHealthcheckRouteBuilder(middlewareWrapper middleware.MiddlewareWrapper) HealthcheckRouteBuilder {
	return healthcheckRouteBuilder{
		middlewareWrapper: middlewareWrapper,
	}
}

type healthcheckRouteBuilder struct {
	middlewareWrapper middleware.MiddlewareWrapper
}

func (b healthcheckRouteBuilder) RoutesForHealthchecks(checks ...gohealthcheck.Healthcheckable) gohealthcheck.Routes {
	routes := gohealthcheck.RoutesForHealthchecks(checks...)
	return gohealthcheck.Routes{
		Basic: gohealthcheck.HealthcheckRoute{
			Handler: b.middlewareWrapper.AddMiddlewareToHandler(routes.Basic.Handler),
			Name:    routes.Basic.Name,
			Path:    routes.Basic.Path,
		},
		Advanced: gohealthcheck.HealthcheckRoute{
			Handler: b.middlewareWrapper.AddMiddlewareToHandler(routes.Advanced.Handler),
			Name:    routes.Advanced.Name,
			Path:    routes.Advanced.Path,
		},
	}
}
