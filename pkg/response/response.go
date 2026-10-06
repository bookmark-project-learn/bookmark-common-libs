package response

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type MessageResponse struct {
	Message string "json:\"message\""
}

type Response[T any] struct {
	Message string "json:\"message\""
	Data    *T     "json:\"data,omitempty\""
}

// Message represents the structure of a response message
type Message struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Common response messages
var (
	InternalErrResponse = Message{
		"Processing error",
		nil,
	}
	InputErrResponse = Message{
		"Input error",
		nil,
	}
)

func InputFieldError(err error) Message {
	if ok := errors.As(err, &validator.ValidationErrors{}); !ok {
		return InputErrResponse
	}

	var errs []string
	// cast err to validator.ValidationErrors for check model validation
	// https://pkg.go.dev/github.com/go-playground/validator/v10@v10.30.1#section-readme
	// need check twice cause can be not work in some cases. for details: read the link
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid ("+err.Tag()+")")
	}

	return Message{
		Message: "Input error",
		Details: errs,
	}
}

// ToMessageReposnse
func ToMessageReposnse(err error) MessageResponse {
	m := InputFieldError(err)
	details := []string{}
	if m.Details == nil {
		details = []string{err.Error()}
	} else {
		details = m.Details.([]string)
	}
	errs := details
	defaultMessage := m.Message
	message := strings.Join(errs, "; ")
	if message == "" {
		message = defaultMessage
	}
	return MessageResponse{
		Message: message,
	}

}

// ToDataResponse converts a Message struct to a Response struct.
// error only related to input validation, so the data field is set to nil.
func ToDataResponse[T any](err error) Response[T] {
	messRes := ToMessageReposnse(err)
	message := messRes.Message
	return Response[T]{
		Message: message,
		Data:    nil,
	}
}
