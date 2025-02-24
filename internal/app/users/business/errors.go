package business

import (
	"fmt"
)

type ErrorCode string

const (
	ErrorUnknown         ErrorCode = "user_error_unknown"
	ErrorInvalidArgument ErrorCode = "user_invalid_argument"
	ErrorRecordNotFound  ErrorCode = "user_record_not_found"
	ErrorDuplicateRecord ErrorCode = "user_duplicate_record"
)

type Error struct {
	orig      error
	code      ErrorCode
	msgToUser string
}

func (e Error) Error() string {
	if e.orig != nil {
		return e.orig.Error()
	}

	return e.msgToUser
}

func (e Error) MessageToUser() string {
	return e.msgToUser
}

func (e Error) Unwrap() error {
	return e.orig
}

func (e *Error) Code() ErrorCode {
	return e.code
}

func WrapErrorf(orig error, code ErrorCode, format string, arguments ...any) error {
	return &Error{
		orig:      orig,
		code:      code,
		msgToUser: fmt.Sprintf(format, arguments...),
	}
}

func NewErrorf(code ErrorCode, format string, a ...any) error {
	return WrapErrorf(nil, code, format, a...)
}
