package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ChatState string

const (
	ChatStarted  = "chat_started"
	ChatInactive = "chat_inactive"
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

func (c *Chat) AddEvent(e *Event, userID string) error {
	e.CreationTimestamp = uint64(time.Now().Nanosecond())
	e.AuthorID = userID
	e.ChatID = c.ChatID

	if err := c.DB.Save(e).Error; err != nil {
		return err
	}

	return nil
}
