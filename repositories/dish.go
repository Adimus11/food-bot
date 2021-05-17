package repositories

import (
	"fmt"
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

type dishPair struct {
	dishesh   []*pair
	penalties map[string]int
}

func (d *dishPair) Len() int {
	return len(d.dishesh)
}

func (d *dishPair) Less(i, j int) bool {
	iPenalty := d.penalties[d.dishesh[i].dish.DishID]
	jPenalty := d.penalties[d.dishesh[j].dish.DishID]

	if jPenalty < iPenalty {
		return false
	}

	return d.dishesh[i].score < d.dishesh[j].score
}

func (d *dishPair) Swap(i, j int) {
	tmp := d.dishesh[i]
	d.dishesh[i] = d.dishesh[j]
	d.dishesh[j] = tmp
}

func (d *dishPair) getTop() []*models.Dish {
	sort.Sort(d)
	topValues := make([]*models.Dish, 0, 4)

	for _, pair := range d.dishesh {
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

func (dr *DishesRepository) GetDishesForIngredients(user *models.User, ingredients []string) ([]*models.Dish, error) {
	dishesh, err := dr.GetDishes()
	if err != nil {
		return nil, err
	}

	result := &dishPair{
		dishesh:   make([]*pair, 0, len(dishesh)),
		penalties: make(map[string]int, 5),
	}

	for _, dish := range dishesh {
		score := 0
		dishIngrid := strings.Join(dish.Ingredients, ",")
		for _, userIngrd := range ingredients {
			if strings.Contains(dishIngrid, userIngrd) {
				score++
			}
		}

		result.dishesh = append(result.dishesh, &pair{
			score: score,
			dish:  dish,
		})
	}

	ratings, err := user.GetLastRatings()
	if err != nil {
		return nil, err
	}

	for i, rating := range ratings {
		if _, exists := result.penalties[rating.DishID]; !exists {
			result.penalties[rating.DishID] = 0
		}

		result.penalties[rating.DishID] += (len(ratings) - i)
	}

	for dish, penalty := range result.penalties {
		fmt.Printf("Dish %s with penalty: %d\n", dish, penalty)
	}

	return result.getTop(), nil
}

func (dr *DishesRepository) GetDish(dishID string) (*models.Dish, error) {
	dish := &models.Dish{}
	err := dr.db.First(&dish, "dish_id = ?", dishID).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Println("No dish")
		return nil, errs.ErrNotFound
	}

	return dish, err
}
