package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hexagonal/internal/app/users/business"
	"github.com/hexagonal/pkg/server"
)

type DeleteUserHandler struct {
	remover business.UserBusiness
}

func NewDeleteUserHandler(remover business.UserBusiness) (*DeleteUserHandler, error) {
	if remover == nil {
		return nil, errors.New("missing 'user business' dependency")
	}

	return &DeleteUserHandler{
		remover: remover,
	}, nil
}

func (p DeleteUserHandler) Delete() server.HttpHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		id := chi.URLParam(r, "user_id")

		userId, err := business.UserID(id).Validate()
		if err != nil {
			return err
		}

		if err := p.remover.Delete(ctx, &userId); err != nil {
			return err
		}

		return server.JSON(w, http.StatusOK, server.HttpResponse{
			Message: "deleted user",
		})
	}
}
