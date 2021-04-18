package api

import (
	"fooder/config"
	"fooder/repositories"
	"fooder/services"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type API struct {
	UsersRepository *repositories.UsersRepository
	ChatsRepository *repositories.ChatsRepository
	BotService      *services.BotService
}

func NewRouter(config *config.Config, api *API) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range GetRoutes(config, api) {
		var handler http.Handler
		log.Println(route.Name)
		handler = route.HandlerFunc

		for _, middleware := range route.Middlewares {
			handler = middleware(handler)
		}

		router.
			Methods(route.Method, "OPTIONS").
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
