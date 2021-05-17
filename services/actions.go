package services

import (
	"context"
	"encoding/json"
	"fmt"
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

func (bs *BotService) actionForState(e *models.Event, state string, user *models.User) ([]*models.Event, string, error) {
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
			dishesh, err := bs.dishService.GetDishesForIngredients(user, res.Ingridients)
			if err != nil {
				return nil, "", err
			}

			msg, err := bs.getMessageEvent("This is what I found for you!")
			if err != nil {
				return nil, "", err
			}
			response = append(response, msg)

			dishSelection := &objects.DishSelection{
				Options: make([]*objects.Option, 0, 3),
			}

			for index, dish := range dishesh {
				if index >= 3 {
					break
				}
				dishSelection.Options = append(dishSelection.Options, &objects.Option{
					OptionID:   index,
					OptionText: dish.Title,
					DishID:     dish.DishID,
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

		userEvent, ok := e.ParsedEvent.(*objects.DishSelection)
		if !ok {
			return nil, "", errs.WrongInterfaceError(e.ParsedEvent, "*objects.DishSelection")
		}

		dish, err := bs.dishService.GetDish(userEvent.GetDishIDForOption())
		if err != nil {
			return nil, "", err
		}

		userMsg, err := bs.getMessageEventForAuthor(fmt.Sprintf("Hi I would like `%s` for my dish", dish.Title), user.UserID)
		if err != nil {
			return nil, "", err
		}
		response = append(response, userMsg)

		data, err := json.Marshal(dish.ToCard())
		if err != nil {
			return nil, "", err
		}

		response = append(response, &models.Event{
			Type:     objects.CardEventType,
			AuthorID: "bot",
			Body:     data,
		})

		botMsg, err := bs.getMessageEvent("How you liked your meal?")
		if err != nil {
			return nil, "", err
		}
		response = append(response, botMsg)

		ratingEvent := &objects.RatingEvent{
			DishID: dish.DishID,
		}
		data, err = json.Marshal(ratingEvent)
		if err != nil {
			return nil, "", err
		}
		response = append(response, &models.Event{
			Type:     objects.RatingEventType,
			AuthorID: "bot",
			Body:     data,
		})

		newState = models.WaitingForReview

	case models.WaitingForReview:
		if e.Type != objects.RatingEventType {
			err = errs.ErrWrongMsgTypeInState
			break
		}

		userEvent, ok := e.ParsedEvent.(*objects.RatingEvent)
		if !ok {
			return nil, "", errs.WrongInterfaceError(e.ParsedEvent, "*objects.DishSelection")
		}

		if err := user.RateDish(userEvent.DishID, *userEvent.Rating); err != nil {
			return nil, "", err
		}

		msg, err := bs.getMessageEventForAuthor(fmt.Sprintf("I think it's solid `%d`", *userEvent.Rating), user.UserID)
		if err != nil {
			return nil, "", err
		}
		response = append(response, msg)

		msg, err = bs.getMessageEvent("Thanks for you opinion, remember I'm always here for you to help!")
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

		case models.WaitingForReview:
			msg, err := bs.getMessageEvent("Please rate dish first using rating selection!")
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
