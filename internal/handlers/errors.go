package handlers

type ErrorCode int32

// TODO: Error codes could came from configuration server
const (
	InavalidBody          ErrorCode = 54000
	InvalidToken                    = 54001
	UserNotRegistered               = 54002
	UserAlreadyRegistered           = 54003
)

var (
	InavalidBodyErr RestError = RestError{
		Code:    InavalidBody,
		Message: "request body should a valid JSON",
	}

	InavalidTokenErr RestError = RestError{
		Code:    InvalidToken,
		Message: "invalid token",
	}

	UserNotRegisteredErr RestError = RestError{
		Code:    InvalidToken,
		Message: "user is not registered",
	}

	UserAlreadyRegisteredErr RestError = RestError{
		Code:    InvalidToken,
		Message: "user is already registered",
	}
)

type RestError struct {
	Code    ErrorCode
	Message string
}

func (e *RestError) Error() string {
	return e.Message
}
