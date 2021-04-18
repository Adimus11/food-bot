package db

import (
	"fmt"
	"fooder/config"
	"log"

	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func InitDB(c *config.Config) (*gorm.DB, error) {
	username := c.Services.DB.User
	dbName := c.Services.DB.DBName
	dbHost := c.Services.DB.URL
	dbPassword := c.Services.DB.Password

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, username, dbName, dbPassword)
	log.Println(dbUri)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	db = conn
	return conn.Set("gorm:auto_preload", true), nil
}
