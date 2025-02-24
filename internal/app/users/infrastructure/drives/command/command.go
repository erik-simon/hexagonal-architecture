package command

import (
	"errors"

	"github.com/google/uuid"
	"github.com/hexagonal/internal/app/users/business"
	"github.com/hexagonal/pkg/command"
	"github.com/spf13/cobra"
)

const (
	userIDFlag   = "id"
	usernameFlag = "username"
)

func UserErrorCommand(commandFunc command.CommandFunc) func(cmd *cobra.Command, args []string) {
	if commandFunc == nil {
		panic("missing command func")
	}

	return func(cmd *cobra.Command, args []string) {
		traceID := uuid.New()

		err := commandFunc(cmd, args)

		if err == nil {
			return
		}

		// TODO: HANDLE COBRA ERRORS

		asError := &business.Error{}

		if errors.As(err, &asError) {
			message := command.CommandResponseError{
				Code:             asError.Code(),
				MessageToUser:    asError.MessageToUser(),
				ErrorDescription: asError.Unwrap().Error(),
				TraceId:          traceID,
			}

			switch asError.Code() {
			case business.ErrorInvalidArgument,
				business.ErrorDuplicateRecord,
				business.ErrorRecordNotFound:

				cmd.PrintErr(message)
				return
			default:
				cmd.PrintErr(command.CommandResponseError{
					Code:             business.ErrorUnknown,
					MessageToUser:    "Failed user process",
					ErrorDescription: err,
					TraceId:          traceID,
				})
				return
			}
		}

		cmd.PrintErr(command.CommandResponseError{
			Code:             business.ErrorUnknown,
			MessageToUser:    "Failed http request process",
			ErrorDescription: err,
			TraceId:          traceID,
		})
	}
}
