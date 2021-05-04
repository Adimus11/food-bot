package models

import (
	"encoding/json"
	"fooder/api/utils"
	"fooder/objects"
	"net/http"

	"github.com/jinzhu/gorm"
)

type Event struct {
	gorm.Model  `json:"-"`
	ChatID      string          `json:"-"`
	EventID     string          `json:"event_id"`
	Type        string          `json:"type"`
	AuthorID    string          `json:"author_id"`
	Body        json.RawMessage `json:"body"`
	ParsedEvent interface{}     `json:"-" sql:"-"`

	CreationTimestamp uint64 `json:"creation_timestamp"`
}

func (e *Event) ParseEvent() error {
	if e.ParsedEvent != nil {
		return nil
	}

	var eventBody interface{}
	switch e.Type {
	case objects.CardEventType:
		eventBody = &objects.CardEvent{}
	case objects.ChatIdleEventType:
		eventBody = &objects.ChatIdleEvent{}
	case objects.MessageEventType:
		eventBody = &objects.MessageEvent{}
	case objects.RatingEventType:
		eventBody = &objects.RatingEvent{}
	}

	if err := json.Unmarshal(e.Body, &eventBody); err != nil {
		return err
	}

	e.ParsedEvent = eventBody
	return nil
}

func (e *Event) ValidateEvent() error {
	if err := e.ParseEvent(); err != nil {
		return err
	}
	switch event := e.ParsedEvent.(type) {
	case *objects.MessageEvent:
		if event.Message == "" {
			return &utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Reason:     "`message` can't be empty",
			}
		}
	case *objects.RatingEvent:
		if event.DishID == "" {
			return &utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Reason:     "`dish_id` can't be empty",
			}
		}
		if event.Rating == nil {
			return &utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Reason:     "`rating` can't be empty",
			}
		}
		if *event.Rating <= 0 || *event.Rating > 5 {
			return &utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Reason:     "`rating` should be from [1-5]",
			}
		}
	case *objects.DishSelection:
		if event.SelectedOptionID == "" {
			return &utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Reason:     "`selected_option_id` can't be empty",
			}
		}
	default:
		return &utils.ApiError{
			StatusCode: http.StatusBadRequest,
			Reason:     "Unprocessable event",
		}
	}

	return nil
}
