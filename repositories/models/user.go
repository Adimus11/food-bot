package models

import (
	"fooder/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type RatedDish struct {
	gorm.Model
	UserID string
	DishID string
	Rate   int
}

type User struct {
	gorm.Model
	DB           *gorm.DB `sql:"-"`
	UserID       string
	ChatID       string
	RatedDishesh []*RatedDish
}

func (u *User) GetToken() string {
	tk := &Token{UserId: u.UserID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.RetreiveConfig().Auth.JWTPassword))

	return tokenString
}

func (u *User) RateDish(dishID string, rating int) error {
	dr := &RatedDish{
		UserID: u.UserID,
		DishID: dishID,
		Rate:   rating,
	}

	return u.DB.Save(dr).Error
}

func (u *User) GetLastRatings() ([]*RatedDish, error) {
	var ratings []*RatedDish

	if err := u.DB.Limit(5).Find(&ratings, "user_id = ?", u.UserID).Error; err != nil {
		return nil, err
	}

	return ratings, nil
}
