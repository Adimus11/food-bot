package models

import (
	"fooder/objects"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Dish struct {
	gorm.Model
	DishID      string         `json:"dish_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Image       string         `json:"img"`
	Link        string         `json:"link"`
	Ingredients pq.StringArray `gorm:"type:text[]" json:"ingredients,omitempty"`
}

func NewDish(dish *Dish) *Dish {
	dish.DishID = uuid.New().String()
	return dish
}

func (d *Dish) ToCard() *objects.CardEvent {
	return &objects.CardEvent{
		DishID:      d.DishID,
		Title:       d.Title,
		Description: d.Description,
		Image:       d.Image,
		Link:        d.Link,
	}
}
