package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	statusOk    = "OK"
	statusError = "Error"
)

type Response struct {
	Status  string `json:"status"`
	Message string `json:"result"`
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json") // add header method to response
	// w.WriteHeader(status)
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	return json.NewEncoder(w).Encode(data)
}

func GeneralSuccess(status string) Response {
	return Response{
		Status:  statusOk,
		Message: status,
	}
}
func GeneralError(err error) Response {
	return Response{
		Status:  statusError,
		Message: err.Error(),
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid ", err.Field()))
		}
	}
	return Response{
		Status:  statusError,
		Message: strings.Join(errMsgs, ","),
	}
}
