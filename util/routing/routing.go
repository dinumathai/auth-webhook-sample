package routing

import (
	"net/http"
	"strings"
	"time"

	"github.com/dinumathai/auth-webhook-sample/log"

	"github.com/gorilla/mux"
)

//Route ...
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// HTTP method names

//GET verb
const GET string = "GET"

//POST verb
const POST string = "POST"

//PUT verb
const PUT string = "PUT"

//PATCH verb
const PATCH string = "PATCH"

//DELETE verb
const DELETE string = "DELETE"

//Routes paths
type Routes []Route

//BuildRouter Builds a Mux router from the given route definitions
func BuildRouter(routes Routes) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.Use(loggingMiddleware(router))

	for _, route := range routes {
		router.
			Methods(strings.Split(route.Method, ",")...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

//GetPathVariables returns map of url path parametres
func GetPathVariables(r *http.Request) map[string]string {
	return mux.Vars(r)
}

//loggingMiddleware Improves traceability by performing request logging before and after the main handler
func loggingMiddleware(router *mux.Router) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				match mux.RouteMatch

				routeName = "NotMatchedRoute"
				start     = time.Now()
			)

			if router.Match(r, &match) {
				routeName = match.Route.GetName()
			}

			log.Debugf("Request received: [%s] %s %s", routeName, r.Method, r.RequestURI)
			defer log.Debugf("Request handled in: %s", time.Since(start))

			next.ServeHTTP(w, r)
		})
	}
}
