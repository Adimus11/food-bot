package models

import "github.com/jinzhu/gorm"

func DoMigration(db *gorm.DB) {
	db.Debug().AutoMigrate(
		&RatedDish{},
	)
	db.Debug().AutoMigrate(
		&User{},
	)
	db.Debug().AutoMigrate(
		&Event{},
	)
	db.Debug().AutoMigrate(
		&Chat{},
	)
	db.Debug().AutoMigrate(
		&Dish{},
	)
}
