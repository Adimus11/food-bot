package models

import (
	"fooder/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserID string
	ChatID string
}

func (u *User) GetToken() string {
	tk := &Token{UserId: u.UserID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(config.RetreiveConfig().Auth.JWTPassword))

	return tokenString
}
