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

func (bs *BotService) RespondForEvent(c *models.Chat, e *models.Event, user *models.User) ([]*models.Event, error) {
	events, newState, err := bs.actionForState(e, c.State, user)
	if err != nil {
		return nil, err
	}

	for _, event := range events {
		if err := c.AddEvent(event, "bot"); err != nil {
			return nil, err
		}
	}

	return events, c.SetState(newState)
}
