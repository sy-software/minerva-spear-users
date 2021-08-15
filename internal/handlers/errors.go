package handlers

import "net/http"

type ErrorCode int32

// TODO: Error codes could came from configuration server
const (
	InavalidBody          ErrorCode = 54000
	InvalidToken                    = 54001
	UserNotRegistered               = 54002
	UserAlreadyRegistered           = 54003
	InternalError                   = 54004
	InavalidRequest                 = 54005
)

var (
	InavalidBodyErr RestError = RestError{
		Code:       InavalidBody,
		Message:    "request body should a valid JSON",
		HTTPStatus: http.StatusBadRequest,
	}

	InavalidTokenErr RestError = RestError{
		Code:       InvalidToken,
		Message:    "invalid token",
		HTTPStatus: http.StatusBadRequest,
	}

	UserNotRegisteredErr RestError = RestError{
		Code:       UserNotRegistered,
		Message:    "user is not registered",
		HTTPStatus: http.StatusNotFound,
	}

	UserAlreadyRegisteredErr RestError = RestError{
		Code:       UserAlreadyRegistered,
		Message:    "user is already registered",
		HTTPStatus: http.StatusBadRequest,
	}

	InternalServerError RestError = RestError{
		Code:       InternalError,
		Message:    "internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}

	InvalidRequestError RestError = RestError{
		Code:       InavalidRequest,
		Message:    "ivalid request",
		HTTPStatus: http.StatusBadRequest,
	}
)

type RestError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	HTTPStatus int       `json:"-"`
}

func (e *RestError) Error() string {
	return e.Message
}
