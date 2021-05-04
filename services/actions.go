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

func (bs *BotService) getMessageEvent(msg string) (*models.Event, error) {
	return bs.getMessageEventForAuthor(msg, "bot")
}

func (bs *BotService) getMessageEventForAuthor(msg, authorID string) (*models.Event, error) {
	msgObject := &objects.MessageEvent{
		Message: msg,
	}

	data, err := json.Marshal(msgObject)
	if err != nil {
		return nil, err
	}

	return &models.Event{
		Type:     objects.MessageEventType,
		AuthorID: authorID,
		Body:     data,
	}, nil
}

func (bs *BotService) actionForState(e *models.Event, state string) ([]*models.Event, string, error) {
	newState := ""
	response := make([]*models.Event, 0, 2)
	var err error

	switch state {
	case models.ChatStarted:
		if e.Type != objects.MessageEventType {
			err = errs.ErrWrongMsgTypeInState
			break
		}

		msg, err := bs.getMessageEvent("Hi! From which ingredients do you want to eat?")
		if err != nil {
			return nil, "", err
		}

		response = append(response, msg)
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

		if len(res.Ingridients) == 0 {
			msg, err := bs.getMessageEvent("Sorry, I couldn't recognize any ingredients in your response, could you repeat?")
			if err != nil {
				return nil, "", err
			}
			newState = state
			response = append(response, msg)
		} else {
			dishesh, err := bs.dishService.GetDishesForIngredients(res.Ingridients)
			if err != nil {
				return nil, "", err
			}

			msg, err := bs.getMessageEvent("This is what I found for you!")
			if err != nil {
				return nil, "", err
			}
			response = append(response, msg)

			dishSelection := &objects.DishSelection{
				Options: make([]*objects.Option, 0, len(dishesh)),
			}

			for index, dish := range dishesh {
				dishSelection.Options = append(dishSelection.Options, &objects.Option{
					OptionID:   index,
					OptionText: dish.Title,
					DishID:     dish.Title,
				})
			}

			data, err := json.Marshal(dishSelection)
			if err != nil {
				return nil, "", err
			}

			response = append(response, &models.Event{
				Type:     objects.SelectEventType,
				AuthorID: "bot",
				Body:     data,
			})
			newState = models.WaitingForChoosingDish
		}
	case models.WaitingForChoosingDish:
		if e.Type != objects.SelectEventType {
			err = errs.ErrWrongMsgTypeInState
			break
		}

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

		response = append(response, &models.Event{
			Type:     objects.CardEventType,
			AuthorID: "bot",
			Body:     data,
		})
		newState = models.DishSelected

	case models.DishSelected:
		//Todo Sent dish and opinion request

	case models.WaitingForReview:
		if e.Type != objects.MessageEventType {
			err = errs.ErrWrongMsgTypeInState
			break
		}

		msg, err := bs.getMessageEvent("Thanks for you opinion, remember I'm always here for you to help")
		if err != nil {
			return nil, "", err
		}
		response = append(response, msg)
		newState = models.ChatStarted
	}

	if err == errs.ErrWrongMsgTypeInState {
		switch state {
		case models.WaitingForChoosingDish:
			msg, err := bs.getMessageEvent("Please use buttons first, so I can send you dish you want")
			if err != nil {
				return nil, "", err
			}
			response = append(response, msg)

		default:
			msg, err := bs.getMessageEvent("Sorry I expected you to send me a message to understood you")
			if err != nil {
				return nil, "", err
			}
			response = append(response, msg)

		}
		newState = state
	} else if err != nil {
		return nil, "", err
	}

	return response, newState, nil
}
