package models

import (
	"encoding/json"
	"fooder/objects"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ChatState string

const (
	ChatStarted  = "chat_started"
	ChatInactive = "chat_inactive"
)

type Event struct {
	gorm.Model
	ChatID     string          `json:"-" gorm:"foreignKey:ChatID"`
	Type       string          `json:"type"`
	AuthorID   string          `json:"author_id"`
	Body       json.RawMessage `json:"-"`
	ParsedBody interface{}     `sql:"-" json:"body"`

	CreationTimestamp uint64 `json:"creation_timestamp"`
}

func (e *Event) ParseEvent() error {
	var eventBody interface{}
	switch e.Type {
	case objects.CardEventType:
		eventBody = &objects.CardEvent{}
	case objects.ChatIdleEventType:
		eventBody = &objects.ChatIdleEvent{}
	case objects.MessageEventType:
		eventBody = &objects.MessageEvent{}
	case objects.RatingRequestedEventType, objects.RatingSetEventType:
		eventBody = &objects.RatingEvent{}
	}

	if err := json.Unmarshal(e.Body, &eventBody); err != nil {
		return err
	}

	e.ParsedBody = eventBody
	return nil
}

type Chat struct {
	gorm.Model
	DB     *gorm.DB `sql:"-"`
	UserID string
	ChatID string
	State  string
	Events []*Event
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

func (c *Chat) ParseEvents() error {
	for _, e := range c.Events {
		if err := e.ParseEvent(); err != nil {
			return err
		}
	}

	return nil
}
