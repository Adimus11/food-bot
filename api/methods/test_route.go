package methods

import (
	"context"
	"fooder/api/utils"
	"net/http"
)

type TestRoute struct {
}

type TestRoutePayload struct {
}

type TestRouteResponse struct {
	Status string `json:"status"`
}

func NewTestRoute() *TestRoute {
	return &TestRoute{}
}

func (t *TestRoute) Payload() interface{} {
	return &TestRoutePayload{}
}

func (t *TestRoute) Do(ctx context.Context, vars map[string]string, payload interface{}) (*utils.ApiResponse, error) {
	return &utils.ApiResponse{
		StatusCode: http.StatusOK,
		Response: &TestRouteResponse{
			Status: "ok",
		},
	}, nil
}
