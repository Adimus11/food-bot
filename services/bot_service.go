package services

import (
	"fooder/repositories"
	"fooder/repositories/models"
	pb "fooder/services/proto"
)

type BotService struct {
	chatsService *repositories.ChatsRepository
	nlpClient    pb.IngridientsServiceClient
	dishService  *repositories.DishesRepository
}

func NewBotService(chatsService *repositories.ChatsRepository, nlpClient pb.IngridientsServiceClient, dishService *repositories.DishesRepository) *BotService {
	return &BotService{chatsService: chatsService, nlpClient: nlpClient, dishService: dishService}
}

func (bs *BotService) RespondForEvent(c *models.Chat, e *models.Event) (*models.Event, error) {
	event, newState, err := bs.actionForState(e, c.State)
	if err != nil {
		return nil, err
	}

	if err := c.AddEvent(event, "bot"); err != nil {
		return nil, err
	}

	return event, c.SetState(newState)
}
