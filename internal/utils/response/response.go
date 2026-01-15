package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

// Create function to send general error
func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

// Function for create validation error
func ValidationError(errs validator.ValidationErrors) Response {
	// Create a slice to store all validation error messages
	var errMsgs []string

	// Running loop to store the error in errMsgs map.
	for _, err := range errs {
		// err.ActualTag() → gives the validation rule that failed (e.g., "required", "email", "gte")
		// err.Field()     → gives the struct field name that failed (e.g., "Name", "Email")
		switch err.ActualTag() {
		// If the failed rule is "required", create a specific error message
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", err.Field()))
		// Default case handles all other validation failures
		// Example: email format, number range, min/max, etc.
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	// Join all error messages into one string separated by commas
	// Example:
	//   errMsgs = []string{"field Name is required", "field Email is invalid"}
	//   After Join → "field Name is required, field Email is invalid"
	combinedErrors := strings.Join(errMsgs, ", ")

	// Return a Response struct with Status and Error message
	return Response{
		Status: StatusError,    // custom constant indicating "error" response
		Error:  combinedErrors, // full error message string
	}

}
