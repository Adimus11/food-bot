package api

import (
	"fooder/config"
	"fooder/dispatcher"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(config *config.Config, dispatcher *dispatcher.Dispatcher, rc *dispatcher.RedisConsumer) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range GetRoutes(config, dispatcher, rc) {
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
