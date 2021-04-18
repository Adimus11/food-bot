package services

import (
	"encoding/json"
	"fooder/objects"
	"fooder/repositories/models"
)

func (bs *BotService) actionForEvent(e *models.Event, state string) (*models.Event, string, error) {
	switch e.Type {
	case "message":
		return bs.actionForMessage(e, state)
	default:
		return nil, "", nil
	}
}

func (bs *BotService) actionForMessage(e *models.Event, state string) (*models.Event, string, error) {
	newState := ""
	response := &models.Event{}
	switch state {
	case models.ChatStarted:
		msg := &objects.MessageEvent{
			Message: "Hi! From which ingredients do you want to eat?",
		}

		data, err := json.Marshal(msg)
		if err != nil {
			return nil, "", err
		}

		response = &models.Event{
			Type:     objects.MessageEventType,
			AuthorID: "bot",
			Body:     data,
		}
		newState = models.WaitingForIngredients
	case models.WaitingForIngredients:
		card := &objects.CardEvent{
			DishID:      "tmp",
			Title:       "Scrambled Eggs",
			Description: "Simple but nutritious dish",
			Image:       "https://images.immediate.co.uk/production/volatile/sites/30/2020/08/recipe-image-legacy-id-1201452_12-7f7a0fa.jpg?quality=90&webp=true&resize=440,400",
			Link:        "https://www.bbcgoodfood.com/recipes/perfect-scrambled-eggs-recipe",
		}

		data, err := json.Marshal(card)
		if err != nil {
			return nil, "", err
		}

		response = &models.Event{
			Type:     objects.CardEventType,
			AuthorID: "bot",
			Body:     data,
		}
		newState = models.WaitingForReview
	case models.WaitingForReview:
		msg := &objects.MessageEvent{
			Message: "Thanks for you opinion, remember I'm always here for you to help",
		}

		data, err := json.Marshal(msg)
		if err != nil {
			return nil, "", err
		}

		response = &models.Event{
			Type:     objects.MessageEventType,
			AuthorID: "bot",
			Body:     data,
		}
		newState = models.ChatStarted
	}

	return response, newState, nil
}
