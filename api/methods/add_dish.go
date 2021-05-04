package methods

import (
	"context"
	"fooder/api/utils"
	"fooder/errs"
	"fooder/objects"
	"fooder/repositories"
	"fooder/repositories/models"
	"net/http"
)

type AddDish struct {
	dishRepo *repositories.DishesRepository
}

func NewAddDish(dishRepo *repositories.DishesRepository) *AddDish {
	return &AddDish{dishRepo}
}

func (t *AddDish) Payload() interface{} {
	return &models.Dish{}
}

func (t *AddDish) Do(ctx context.Context, vars map[string]string, payload interface{}) (*utils.ApiResponse, error) {
	dish, ok := payload.(*models.Dish)
	if !ok {
		return nil, errs.WrongInterfaceError(payload, "*models.Dish")
	}
	if err := t.dishRepo.AddDish(dish); err != nil {
		return nil, err
	}

	return &utils.ApiResponse{
		StatusCode: http.StatusOK,
		Response: &objects.BasicResponse{
			Status: "ok",
		},
	}, nil
}
