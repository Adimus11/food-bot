package main

import (
	"flag"
	"fooder/api"
	"fooder/config"
	"fooder/db"
	"fooder/repositories"
	"fooder/repositories/models"
	"fooder/services"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"google.golang.org/grpc"

	pb "fooder/services/proto"
)

func main() {
	time.Sleep(3 * time.Second)
	configFilePath := flag.String("config_file", "./config/config.json", "config.json file path")

	config, err := config.GetConfig(*configFilePath)
	if err != nil {
		log.Panic(err.Error())
	}

	dbClient, err := db.InitDB(config)
	if err != nil {
		log.Panic(err.Error())
	}

	models.DoMigration(dbClient)

	conn, err := grpc.Dial(config.Services.NLPService.URL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Panic("did not connect: %v", err)
	}
	defer conn.Close()

	apiApp := &api.API{
		UsersRepository: repositories.NewUsersRepository(dbClient),
		ChatsRepository: repositories.NewChatsRepository(dbClient),
		DishRepository:  repositories.NewDishesRepository(dbClient),
	}

	nlpClient := pb.NewIngridientsServiceClient(conn)

	apiApp.BotService = services.NewBotService(apiApp.ChatsRepository, nlpClient, apiApp.DishRepository)

	router := api.NewRouter(config, apiApp)
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
