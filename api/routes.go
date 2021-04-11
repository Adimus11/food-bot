package api

import (
	"fooder/api/methods"
	"fooder/api/middlewares"
	"fooder/api/utils"
	"fooder/config"
	"fooder/dispatcher"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middlewares []func(next http.Handler) http.Handler
}

var baseMiddlewares []func(next http.Handler) http.Handler = []func(next http.Handler) http.Handler{
	func(next http.Handler) http.Handler {
		return handlers.CombinedLoggingHandler(os.Stderr, next)
	},
	middlewares.NewCORSMiddleware(),
}

type Routes []Route

func GetRoutes(c *config.Config, dis *dispatcher.Dispatcher, rc *dispatcher.RedisConsumer) Routes {
	return Routes{
		Route{
			Name:    "Test Route",
			Method:  "GET",
			Pattern: "/ping", //tested
			HandlerFunc: utils.HandlerWrapper(
				methods.NewTestRoute(),
			),
			Middlewares: append(baseMiddlewares, []func(next http.Handler) http.Handler{
				middlewares.NewAuthMiddleware(c, "admin"),
			}...),
		},
	}
}
