package methods

import (
	"context"
	"fmt"
	"fooder/api/utils"
	"fooder/errs"
	"fooder/objects"
	"fooder/repositories"
	"fooder/repositories/models"
	"fooder/services"
	"net/http"
	"strings"
)

type SendEventRoute struct {
	cr         *repositories.ChatsRepository
	botService *services.BotService
	ur         *repositories.UsersRepository
}

type SendEventRoutePayload struct {
	*models.Event
}

func NewSendEventRoute(cr *repositories.ChatsRepository, botService *services.BotService, ur *repositories.UsersRepository) *SendEventRoute {
	return &SendEventRoute{cr, botService, ur}
}

func (t *SendEventRoute) Payload() interface{} {
	return &SendEventRoutePayload{}
}

func (t *SendEventRoute) validate(payload *SendEventRoutePayload) error {
	if payload.Event == nil {
		return &utils.ApiError{
			StatusCode: http.StatusBadRequest,
			Reason:     "Missing payload",
		}
	}

	foundType := false
	for _, eventType := range objects.EventTypes {
		if eventType == payload.Type {
			foundType = true
			break
		}
	}

	if !foundType {
		return &utils.ApiError{
			StatusCode: http.StatusBadRequest,
			Reason:     fmt.Sprintf("`type` should be on of: `%s`", strings.Join(objects.EventTypes, ",")),
		}
	}

	if err := payload.ValidateEvent(); err != nil {
		return err
	}

	return nil
}

func (t *SendEventRoute) Do(ctx context.Context, vars map[string]string, payload interface{}) (*utils.ApiResponse, error) {
	token, ok := utils.GetTokenFromContext(ctx).(*models.Token)
	if !ok {
		return nil, errs.ErrTokenNotInCtx
	}

	event, ok := payload.(*SendEventRoutePayload)
	if !ok {
		return nil, errs.WrongInterfaceError(payload, "*SendEventRoutePayload")
	}

	if err := t.validate(event); err != nil {
		return nil, err
	}

	user, err := t.ur.GetUser(token.UserId)
	if err != nil {
		return nil, err
	}

	chat, err := t.cr.GetOrCreateChat(token.UserId)
	if err != nil {
		return nil, err
	}

	if err := chat.AddEvent(event.Event, token.UserId); err != nil {
		return nil, err
	}

	responseEvents, err := t.botService.RespondForEvent(chat, event.Event, user)
	if err != nil {
		return nil, err
	}

	return &utils.ApiResponse{
		StatusCode: http.StatusOK,
		Response:   responseEvents,
	}, nil
}
