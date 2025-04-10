package response

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func Ok() Response {
	return Response{Status: StatusOk}
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errsMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errsMsgs = append(errsMsgs, err.Field()+" is required")
		case "url":
			errsMsgs = append(errsMsgs, err.Field()+" is not valid url")
		default:
			errsMsgs = append(errsMsgs, err.Field()+" is invalid")
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errsMsgs, ", "),
	}
}
