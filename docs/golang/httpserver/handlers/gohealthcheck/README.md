[![pipeline status](https://gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/gohealthcheck/badges/master/pipeline.svg)](https://gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/gohealthcheck/commits/master)

# Go Healthcheck

This is a small library to add health checks to micro services.

## What is a health check?

A health check is any upstream dependency of your service that you might want to be able to check on a deployed instance.

## Types of health checks

- Basic -- this only verifies that your web server is up and can respond to HTTP requests. You may, for example, wish to use this as the health check for your load balancer.
- Advanced -- this verifies your web server and also any dependencies that you may have. You may implement the `Healthcheckable` interface to verify that upstream services are alive.

## Example

For example, if your service is a pizza ordering service, you may want to check that your client that communicates with other backend systems is alive and healthy.

```go
import (
  "net/http"

  "gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/gohealthcheck"
)

type pizzaclient struct {
  httpclient http.client
}

// this will be in the body of the healthcheck http response
func (pizzaclient pizzaclient) name() string {
   return "pizza-client"
}

func (pizzaclient pizzaclient) ishealthy() bool {
  _, err := pizzaclient.httpclient.do(pizzaclient.examplerequest)
  return err == nil
}

var _ gohealthcheck.healthcheckable = pizzaclient{}
```

## Usage

```go
package main

import (
  "net/http"

  "gitlab.nordstrom.com/sentry/authorize/httpserver/handlers/gohealthcheck"
  "gitlab.nordstrom.com/sentry/example"
)


func main() {
  var pizzaClient example.PizzaClient = example.NewPizzaClient()
  var handler http.Handler = example.NewOrderPizzaHandler(pizzaClient)

  healthcheckRoutes := gohealthcheck.RoutesForHealthchecks(pizzaClient)

  // create a multiplexor to route traffic to our endpoint and the healthcheck endpoints
  webmux := http.NewServeMux()

  // wire up the correct paths to our endpoints and http handlers
  webmux.Handle("/pizza/order", handler)
  webmux.Handle(healthcheckRoutes.Basic.Path, healthcheckRoutes.Basic.Handler)
  webmux.Handle(healthcheckRoutes.Advanced.Path, healthcheckRoutes.Advanced.Handler)

  // start listening for http requests
  http.ListenAndServe(":8080", webmux)
}
```

## Running tests

- `ginkgo -r`
