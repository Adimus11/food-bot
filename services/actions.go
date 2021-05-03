package services

import (
	"context"
	"encoding/json"
	"fooder/errs"
	"fooder/objects"
	"fooder/repositories/models"
	conn_proto "fooder/services/proto"
	"time"
)

func (bs *BotService) actionForState(e *models.Event, state string) (*models.Event, string, error) {
	newState := ""
	response := &models.Event{}
	var err error

	switch state {
	case models.ChatStarted:
		if e.Type != objects.MessageEventType {
			err = errs.ErrWrongMsgTypeInState
			break
		}
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
		if e.Type != objects.MessageEventType {
			err = errs.ErrWrongMsgTypeInState
			break
		}

		// TODO: Add dish selection
		userResponse := &objects.MessageEvent{}
		if err = json.Unmarshal(e.Body, &userResponse); err != nil {
			return nil, "", err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		res, err := bs.nlpClient.GetIngridients(ctx, &conn_proto.UserInput{
			Text: userResponse.Message,
		})
		if err != nil {
			return nil, "", err
		}

		var (
			responseEvent     interface{}
			responseEventType string
		)
		if len(res.Ingridients) == 0 {
			responseEvent = &objects.MessageEvent{
				Message: "Sorry, I couldn't recognize any ingredients in your response, could you repeat?",
			}
			responseEventType = objects.MessageEventType
			newState = state
		} else {
			dishesh, err := bs.dishService.GetDishesForIngredients(res.Ingridients)
			if err != nil {
				return nil, "", err
			}

			dishSelection := &objects.DishSelection{
				Message: "This is what I found for you!",
				Options: make([]*objects.Option, 0, len(dishesh)),
			}

			for index, dish := range dishesh {
				dishSelection.Options = append(dishSelection.Options, &objects.Option{
					OptionID:   index,
					OptionText: dish.Title,
					DishID:     dish.Title,
				})
			}

			responseEvent = dishSelection
			responseEventType = objects.SelectEventType
			newState = models.WaitingForChoosingDish
		}

		data, err := json.Marshal(responseEvent)
		if err != nil {
			return nil, "", err
		}

		response = &models.Event{
			Type:     responseEventType,
			AuthorID: "bot",
			Body:     data,
		}

	case models.WaitingForChoosingDish:
		if e.Type != objects.SelectEventType {
			err = errs.ErrWrongMsgTypeInState
			break
		}

		//TODO: select dish and return it
		userEvent := &objects.DishSelection{}
		if err := json.Unmarshal(e.Body, userEvent); err != nil {
			return nil, "", err
		}

		dish, err := bs.dishService.GetDish(userEvent.SelectedOptionID)
		if err != nil {
			return nil, "", err
		}

		data, err := json.Marshal(dish.ToCard())
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
		if e.Type != objects.MessageEventType {
			err = errs.ErrWrongMsgTypeInState
			break
		}
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

	if err == errs.ErrWrongMsgTypeInState {
		switch state {
		case models.WaitingForChoosingDish:
			msg := &objects.MessageEvent{
				Message: "Please use buttons first, so I can send you dish you want",
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

		default:
			msg := &objects.MessageEvent{
				Message: "Sorry I expected you to send me a message to understood you",
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
		}
		newState = state
	} else if err != nil {
		return nil, "", err
	}

	return response, newState, nil
}
