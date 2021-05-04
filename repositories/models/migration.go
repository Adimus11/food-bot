package models

import "github.com/jinzhu/gorm"

func DoMIgration(db *gorm.DB) {
	db.Debug().AutoMigrate(
		&RatedDish{},
		&User{},
		&Event{},
		&Chat{},
		&Dish{},
	)
}
