package methods

import (
	"context"
	"fooder/api/utils"
	"fooder/errs"
	"fooder/repositories"
	"fooder/repositories/models"
	"net/http"
)

type HistoryRoute struct {
	cr *repositories.ChatsRepository
}

type HistoryRoutePayload struct {
}

func NewHistoryRoute(cr *repositories.ChatsRepository) *HistoryRoute {
	return &HistoryRoute{cr}
}

func (t *HistoryRoute) Payload() interface{} {
	return &HistoryRoutePayload{}
}

func (t *HistoryRoute) Do(ctx context.Context, vars map[string]string, payload interface{}) (*utils.ApiResponse, error) {
	token, ok := utils.GetTokenFromContext(ctx).(*models.Token)
	if !ok {
		return nil, errs.ErrTokenNotInCtx
	}

	chat, err := t.cr.GetChat(token.UserId)
	if err == errs.ErrNotFound {
		return &utils.ApiResponse{
			StatusCode: http.StatusOK,
			Response:   []*models.Event{},
		}, nil
	} else if err != nil {
		return nil, err
	}

	if err := chat.ParseEvents(); err != nil {
		return nil, err
	}

	return &utils.ApiResponse{
		StatusCode: http.StatusOK,
		Response:   chat.Events,
	}, nil
}
