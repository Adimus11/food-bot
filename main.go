package main

import (
	"flag"
	"fooder/api"
	"fooder/config"
	"fooder/db"
	"fooder/dispatcher"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

func main() {
	time.Sleep(3 * time.Second)
	configFilePath := flag.String("config_file", "./config/config.json", "config.json file path")

	config, err := config.GetConfig(*configFilePath)
	if err != nil {
		log.Panic(err.Error())
	}

	_, err = db.InitDB(config)
	if err != nil {
		log.Panic(err.Error())
	}

	dis := dispatcher.NewDispatcher()
	rc := dispatcher.NewRedisConsumer(config.Services.Redis.URL, dis)
	router := api.NewRouter(config, dis, rc)
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
