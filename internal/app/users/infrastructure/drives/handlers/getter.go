package handlers

import (
	"errors"
	"net/http"

	"github.com/hexagonal/internal/app/users/business"
	"github.com/hexagonal/pkg/server"
)

type GetAllUserHandler struct {
	getter business.UserBusiness
}

func NewGetAllUserHandler(getter business.UserBusiness) (*GetAllUserHandler, error) {
	if getter == nil {
		return nil, errors.New("missing 'user business' dependency")
	}

	return &GetAllUserHandler{
		getter: getter,
	}, nil
}

func (p GetAllUserHandler) GetAll() server.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		users, err := p.getter.SearchAll(ctx)
		if err != nil {
			return err
		}

		return server.JSON(w, http.StatusOK, server.HttpResponse{
			Data: NewUsersResponse(users),
		})
	}
}
