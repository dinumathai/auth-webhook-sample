package server

import (
	"github.com/dinumathai/auth-webhook-sample/api"
	"github.com/dinumathai/auth-webhook-sample/auth"
	"github.com/dinumathai/auth-webhook-sample/types"
	"github.com/dinumathai/auth-webhook-sample/util/health"
	"github.com/dinumathai/auth-webhook-sample/util/routing"
)

//BuildRoutes builds routes for this service
func BuildRoutes(config *types.ConfigMap) []routing.Route {
	var routes = routing.Routes{
		routing.Route{
			Name:        "HealthCheck",
			Method:      "GET",
			Pattern:     "/health",
			HandlerFunc: health.PongHandler,
		},
		routing.Route{
			Name:        "V0-Login",
			Method:      "POST",
			Pattern:     "/v0/login",
			HandlerFunc: api.LoginV0Handler(config),
		},
		routing.Route{
			Name:        "V0-Validate",
			Method:      "POST",
			Pattern:     "/v0/authenticate",
			HandlerFunc: api.ValidationHandler(auth.V0),
		},
		routing.Route{
			Name:        "V0-Authorize",
			Method:      "POST",
			Pattern:     "/v0/authorize",
			HandlerFunc: api.AuthorizeV0Handler(auth.V0),
		},
	}
	return routes
}
