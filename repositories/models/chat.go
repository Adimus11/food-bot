package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type ChatState string

const (
	ChatStarted            = "chat_started"
	WaitingForIngredients  = "chat_ingredients"
	WaitingForChoosingDish = "chat_dish_choosing"
	DishSelected           = "chat_dish_selected"
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

func (c *Chat) AddEvent(e *Event, userID string) error {
	e.CreationTimestamp = uint64(time.Now().Nanosecond())
	e.AuthorID = userID
	e.ChatID = c.ChatID
	e.EventID = uuid.New().String()

	if err := c.DB.Save(e).Error; err != nil {
		return err
	}

	return nil
}

func (c *Chat) SetState(state string) error {
	c.State = state
	return c.DB.Model(&Chat{}).Where("chat_id = ?", c.ChatID).Update("state", state).Error
}
