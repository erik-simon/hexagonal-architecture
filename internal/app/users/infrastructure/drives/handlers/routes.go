package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/hexagonal/internal/app/users/business"
	"github.com/hexagonal/pkg/server"
)

func SetUserRoutes(
	router *chi.Mux,
	postUser *PostUserHandler,
	deleteUser *DeleteUserHandler,
	searchAll *GetAllUserHandler,
) error {
	if postUser == nil {
		return errors.New("missing 'post user' dependency")
	}

	if router == nil {
		return errors.New("missing 'routes' dependency")
	}
	router.Group(func(r chi.Router) {
		r.Post("/create", UserErrorHandler(postUser.Post()))
		r.Delete("/delete/{user_id}", UserErrorHandler(deleteUser.Delete()))
		r.Get("/get-all", UserErrorHandler(searchAll.GetAll()))
	})

	return nil
}

func UserErrorHandler(apiFunc server.HttpHandlerFunc) http.HandlerFunc {
	if apiFunc == nil {
		panic("missing handler func")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		traceID := uuid.New()

		err := apiFunc(w, r)
		if err == nil {
			return
		}

		var syntaxErr *json.SyntaxError

		if errors.As(err, &syntaxErr) {
			_ = server.JSON(w, http.StatusBadRequest, server.HttpResponseError{
				Code:             business.ErrorUnknown,
				MessageToUser:    "Invalid JSON format. Please check the request body.",
				ErrorDescription: fmt.Sprintf("description: %[1]s, offset: %[2]d", syntaxErr.Error(), syntaxErr.Offset),
				TraceId:          traceID,
			})

			return
		}

		var unmarshalTypeErr *json.UnmarshalTypeError

		if errors.As(err, &unmarshalTypeErr) {
			_ = server.JSON(w, http.StatusUnprocessableEntity, server.HttpResponseError{
				Code:          business.ErrorInvalidArgument,
				MessageToUser: "invalid body types",
				ErrorDescription: fmt.Sprintf(
					"Struct: %[1]s, Field: %[2]s, Type: %[3]v, Value: %[4]s",
					unmarshalTypeErr.Struct,
					unmarshalTypeErr.Field,
					unmarshalTypeErr.Type,
					unmarshalTypeErr.Value,
				),
				TraceId: traceID,
			})

			return
		}

		asError := &business.Error{}

		if errors.As(err, &asError) {
			message := server.HttpResponseError{
				Code:             asError.Code(),
				MessageToUser:    asError.MessageToUser(),
				ErrorDescription: asError.Unwrap().Error(),
				TraceId:          traceID,
			}

			switch asError.Code() {
			case business.ErrorInvalidArgument,
				business.ErrorDuplicateRecord:

				_ = server.JSON(w, http.StatusBadRequest, message)
				return
			case business.ErrorRecordNotFound:
				_ = server.JSON(w, http.StatusNotFound, message)
				return
			default:
				_ = server.JSON(w, http.StatusInternalServerError, server.HttpResponseError{
					Code:             business.ErrorUnknown,
					MessageToUser:    "Failed user process",
					ErrorDescription: err,
					TraceId:          traceID,
				})
				return
			}
		}

		_ = server.JSON(w, http.StatusInternalServerError, server.HttpResponseError{
			Code:             business.ErrorUnknown,
			MessageToUser:    "Failed http request process",
			ErrorDescription: err,
			TraceId:          traceID,
		})
	}
}
