package repositories

import (
	"fmt"
	"fooder/repositories/models"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db}
}

func (ur *UsersRepository) GetUser(id string) (*models.User, error) {
	if id == "" {
		fmt.Println("mis")
		id = uuid.New().String()
		fmt.Println(id)
	}
	u := &models.User{UserID: id, DB: ur.db}
	if err := ur.db.Table("users").FirstOrCreate(&u, "user_id = ?", u.UserID).Error; err != nil {
		return nil, err
	}
	return u, nil
}
