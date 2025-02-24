package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/hexagonal/internal/app/users/business"
)

type (
	HttpHandlerFunc func(w http.ResponseWriter, r *http.Request) error

	HttpResponse struct {
		Data    any `json:"data"`
		Message any `json:"message"`
	}

	HttpResponseError struct {
		Code             business.ErrorCode `json:"code"`
		MessageToUser    any                `json:"message_to_user"`
		ErrorDescription any                `json:"error_description"`
		TraceId          uuid.UUID          `json:"trace_id"`
	}
)

func JSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}
