package models

import (
	"fooder/api/utils"
	"fooder/errs"
	"fooder/objects"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ChatState string

const (
	ChatStarted            = "chat_started"
	WaitingForIngredients  = "chat_ingredients"
	WaitingForChoosingDish = "chat_dish_choosing"
	WaitingForReview       = "chat_review"
	ChatInactive           = "chat_inactive"
)

type Chat struct {
	gorm.Model
	DB     *gorm.DB `sql:"-"`
	UserID string
	ChatID string
	State  string
	Events []*Event `gorm:"foreignKey:ChatID;references:ChatID"`
}

func NewChat(userID string, db *gorm.DB) *Chat {
	return &Chat{
		UserID: userID,
		ChatID: uuid.New().String(),
		Events: make([]*Event, 0),
		DB:     db,
		State:  ChatStarted,
	}
}

func (c *Chat) GetEvent(eventID string) *Event {
	for _, event := range c.Events {
		if event.EventID == eventID {
			return event
		}
	}

	return nil
}

func (c *Chat) CreateEventForUser(e *Event, userID string) (*Event, error) {
	switch e.Type {
	case objects.MessageEventType:
		e.CreationTimestamp = uint64(time.Now().Nanosecond())
		e.AuthorID = userID
		e.ChatID = c.ChatID
		e.EventID = uuid.New().String()

	case objects.SelectEventType:
		desiredEvent := c.GetEvent(e.EventID)
		if desiredEvent == nil {
			return nil, errs.ErrNotFound
		}
		if err := desiredEvent.ParseEvent(); err != nil {
			return nil, err
		}

		eventBody, ok := desiredEvent.ParsedEvent.(*objects.DishSelection)
		if !ok {
			return nil, errs.WrongInterfaceError(desiredEvent.ParsedEvent, "*objects.DishSelection")
		}

		userEventBody, ok := e.ParsedEvent.(*objects.DishSelection)
		if !ok {
			return nil, errs.WrongInterfaceError(e.ParsedEvent, "*objects.DishSelection")
		}

		found := false
		for _, option := range eventBody.Options {
			if option.OptionID == *userEventBody.SelectedOptionID {
				found = true
			}
		}

		if !found {
			return nil, &utils.ApiError{
				StatusCode: http.StatusBadRequest,
				Reason:     "`selected_option_id` isn't available in `options`",
			}
		}

		eventBody.SelectedOptionID = userEventBody.SelectedOptionID
		e.ParsedEvent = eventBody
		if err := e.UpdateEventBody(); err != nil {
			return nil, err
		}

	case objects.RatingEventType:
		desiredEvent := c.GetEvent(e.EventID)
		if desiredEvent == nil {
			return nil, errs.ErrNotFound
		}
		if err := desiredEvent.ParseEvent(); err != nil {
			return nil, err
		}

		eventBody, ok := desiredEvent.ParsedEvent.(*objects.RatingEvent)
		if !ok {
			return nil, errs.WrongInterfaceError(desiredEvent.ParsedEvent, "*objects.DishSelection")
		}

		userEventBody, ok := e.ParsedEvent.(*objects.RatingEvent)
		if !ok {
			return nil, errs.WrongInterfaceError(e.ParsedEvent, "*objects.DishSelection")
		}

		eventBody.Rating = userEventBody.Rating
		e.ParsedEvent = eventBody
		if err := e.UpdateEventBody(); err != nil {
			return nil, err
		}
	default:
		return nil, errs.ErrUnavailableTypeForUser
	}

	return e, nil
}

func (c *Chat) AddEvent(e *Event, userID string) error {
	e, err := c.CreateEventForUser(e, userID)
	if err != nil {
		return err
	}

	if err := c.DB.Save(e).Error; err != nil {
		return err
	}

	return nil
}

func (c *Chat) AddEventByBot(e *Event) error {
	e.CreationTimestamp = uint64(time.Now().Nanosecond())
	e.ChatID = c.ChatID
	e.EventID = uuid.New().String()

	return c.DB.Save(e).Error
}

func (c *Chat) SetState(state string) error {
	c.State = state
	return c.DB.Model(&Chat{}).Where("chat_id = ?", c.ChatID).Update("state", state).Error
}
