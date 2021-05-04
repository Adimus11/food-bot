package repositories

import (
	"fmt"
	"fooder/errs"
	"fooder/repositories/models"

	"github.com/jinzhu/gorm"
)

type ChatsRepository struct {
	db *gorm.DB
}

func NewChatsRepository(db *gorm.DB) *ChatsRepository {
	return &ChatsRepository{db}
}

func (cr *ChatsRepository) GetChat(userID string) (*models.Chat, error) {
	chat := &models.Chat{UserID: userID, DB: cr.db}
	err := cr.db.Preload("Events").First(&chat, "user_id = ?", userID).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Println("No chat 2")
		return chat, errs.ErrNotFound
	}

	events := []*models.Event{}
	if err := cr.db.Find(&events, "chat_id = ?", chat.ChatID).Error; err != nil {
		return nil, err
	}
	chat.Events = events

	return chat, err
}

func (cr *ChatsRepository) GetOrCreateChat(userID string) (*models.Chat, error) {
	chat, err := cr.GetChat(userID)
	if err == errs.ErrNotFound {
		fmt.Println("No chat")
		chat = models.NewChat(userID, cr.db)
		if err := cr.db.Save(chat).Error; err != nil {
			return nil, err
		}
	}

	return chat, nil
}
