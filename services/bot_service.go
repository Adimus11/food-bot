package services

import (
	"fooder/repositories"
	"fooder/repositories/models"
)

type BotService struct {
	chatsService *repositories.ChatsRepository
}

func NewBotService(chatsService *repositories.ChatsRepository) *BotService {
	return &BotService{chatsService: chatsService}
}

func (bs *BotService) RespondForEvent(c *models.Chat, e *models.Event) (*models.Event, error) {
	event, newState, err := bs.actionForEvent(e, c.State)
	if err != nil {
		return nil, err
	}

	if err := c.AddEvent(event, "bot"); err != nil {
		return nil, err
	}

	return event, c.SetState(newState)
}
