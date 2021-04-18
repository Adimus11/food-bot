package methods

import (
	"context"
	"fmt"
	"fooder/api/consts"
	"fooder/api/utils"
	"fooder/repositories"
	"net/http"
)

type AuthUserRoute struct {
	ur *repositories.UsersRepository
}

type AuthUserPayload struct{}

type AuthUserResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

func NewAuthUserRoute(ur *repositories.UsersRepository) *AuthUserRoute {
	return &AuthUserRoute{ur}
}

func (r *AuthUserRoute) Payload() interface{} {
	return &AuthUserPayload{}
}

func (r *AuthUserRoute) Do(ctx context.Context, vars map[string]string, payload interface{}) (*utils.ApiResponse, error) {
	sessionData := utils.GetSessionDataFromCtx(ctx)

	fmt.Printf("user_id: %s", sessionData.UserID)

	user, err := r.ur.GetUser(sessionData.UserID)
	if err != nil {
		return nil, err
	}

	return &utils.ApiResponse{
		StatusCode: http.StatusOK,
		Response: &AuthUserResponse{
			Token:  user.GetToken(),
			UserID: user.UserID,
		},
		Cookies: []*http.Cookie{
			{
				Name:     consts.SessionCookie,
				Value:    user.UserID,
				SameSite: http.SameSiteNoneMode,
			},
		},
	}, nil
}
