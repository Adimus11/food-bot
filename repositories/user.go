package repositories

import (
	"fmt"
	"fooder/repositories/models"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (ur *UserRepository) GetUser(id string) (*models.User, error) {
	if id == "" {
		fmt.Println("mis")
		id = uuid.New().String()
		fmt.Println(id)
	}
	u := &models.User{UserID: id}
	if err := ur.db.Table("users").FirstOrCreate(&u, "user_id = ?", u.UserID).Error; err != nil {
		return nil, err
	}
	return u, nil
}
