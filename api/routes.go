package api

import (
	"fooder/api/methods"
	"fooder/api/middlewares"
	"fooder/api/utils"
	"fooder/config"
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

func NewBaseMiddlewares() []func(next http.Handler) http.Handler {
	meth := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origin := handlers.AllowedOrigins([]string{"*"})
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials", "Access-Control-Allow-Origin"})
	credentials := handlers.AllowCredentials()

	return []func(next http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return handlers.CombinedLoggingHandler(os.Stderr, next)
		},
		handlers.CORS(header, meth, origin, credentials),
	}
}

type Routes []*Route

func GetRoutes(c *config.Config, api *API) Routes {
	return Routes{
		&Route{
			Name:    "Test Route",
			Method:  http.MethodGet,
			Pattern: "/ping", //tested
			HandlerFunc: utils.HandlerWrapper(
				methods.NewTestRoute(),
			),
			Middlewares: NewBaseMiddlewares(),
		},
		&Route{
			Name:    "Auth Route",
			Method:  http.MethodPost,
			Pattern: "/auth",
			HandlerFunc: utils.HandlerWrapper(
				methods.NewAuthUserRoute(api.UsersRepository),
			),
			Middlewares: append(NewBaseMiddlewares(), middlewares.NewSessionMiddleware()),
		},
		&Route{
			Name:    "History Route",
			Method:  http.MethodGet,
			Pattern: "/history",
			HandlerFunc: utils.HandlerWrapper(
				methods.NewHistoryRoute(api.ChatsRepository),
			),
			Middlewares: append(NewBaseMiddlewares(), middlewares.NewAuthMiddleware(config.RetreiveConfig())),
		},
	}
}
