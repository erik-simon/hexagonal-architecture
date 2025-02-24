package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hexagonal/internal/app/users/business"
	"github.com/hexagonal/pkg/server"
)

type PostUserHandler struct {
	creator business.UserBusiness
}

func NewPostUserHandler(creator business.UserBusiness) (*PostUserHandler, error) {
	if creator == nil {
		return nil, errors.New("missing 'user business' dependency")
	}

	return &PostUserHandler{
		creator: creator,
	}, nil
}

func (p PostUserHandler) Post() server.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		userResponse := UserResponse{}

		if err := json.NewDecoder(r.Body).Decode(&userResponse); err != nil {
			return err
		}

		user, err := userResponse.ToBusiness()
		if err != nil {
			return err
		}

		if err := p.creator.Create(ctx, user); err != nil {
			return err
		}

		return server.JSON(w, http.StatusCreated, server.HttpResponse{
			Message: "created user",
			Data:    map[string]any{"user": NewUserResponse(user)},
		})
	}
}
