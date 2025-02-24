package command

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/hexagonal/internal/app/users/business"
	"github.com/spf13/cobra"
)

type (
	CommandFunc func(cmd *cobra.Command, args []string) error

	CommandResponseError struct {
		Code             business.ErrorCode `json:"code"`
		MessageToUser    any                `json:"message_to_user"`
		ErrorDescription any                `json:"error_description"`
		TraceId          uuid.UUID          `json:"trace_id"`
	}
)

func (e CommandResponseError) Error() string {
	jsonMessage, _ := json.MarshalIndent(e, "", "  ")
	return string(jsonMessage)
}
