package repositories

import (
	"fooder/errs"
	"fooder/repositories/models"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"
)

type pair struct {
	score int
	dish  *models.Dish
}

type dishPair []*pair

func (d dishPair) Len() int {
	return len(d)
}

func (d dishPair) Less(i, j int) bool {
	return d[i].score < d[j].score
}

func (d dishPair) Swap(i, j int) {
	tmp := d[i]
	d[i] = d[j]
	d[j] = tmp
}

func (d dishPair) getTop() []*models.Dish {
	sort.Sort(d)
	topValues := make([]*models.Dish, 0, 4)

	for _, pair := range d {
		topValues = append(topValues, pair.dish)
	}

	return topValues
}

type DishesRepository struct {
	db *gorm.DB
}

func NewDishesRepository(db *gorm.DB) *DishesRepository {
	return &DishesRepository{db}
}

func (dr *DishesRepository) AddDish(dish *models.Dish) error {
	return dr.db.Save(models.NewDish(dish)).Error
}

func (dr *DishesRepository) GetDishes() ([]*models.Dish, error) {
	var dishes []*models.Dish
	err := dr.db.Find(&dishes).Error
	return dishes, err
}

func (dr *DishesRepository) GetDishesForIngredients(ingredients []string) ([]*models.Dish, error) {
	dishesh, err := dr.GetDishes()
	if err != nil {
		return nil, err
	}

	filtered := make(dishPair, 0, len(dishesh))

	for _, dish := range dishesh {
		score := 0
		dishIngrid := strings.Join(dish.Ingredients, ",")
		for _, userIngrd := range ingredients {
			if strings.Contains(dishIngrid, userIngrd) {
				score++
			}
		}

		filtered = append(filtered, &pair{
			score: score,
			dish:  dish,
		})
	}

	return filtered.getTop(), nil
}

func (dr *DishesRepository) GetDish(dishID string) (*models.Dish, error) {
	dish := &models.Dish{}
	err := dr.db.First(&dish, "dish_id = ?", dishID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, errs.ErrNotFound
	}

	return dish, err
}
