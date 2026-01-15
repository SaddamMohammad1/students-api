package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/SaddamMohammad1/students-api/internal/types"
	"github.com/SaddamMohammad1/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("creating a student")

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		// Send custom error message
		if errors.Is(err, io.EOF) {
			// response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err)) // Acctuall error send
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body"))) // Custom message in error send
			return
		}

		// Send real err message
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Request validation (User request validation that is send in request body by user)
		// Install this package for validation run in terminal -> go get github.com/go-playground/validator/v10
		// Validate the incoming student struct using go-playground/validator
		// validator.New()  → creates a new validator instance
		// .Struct(student) → validates the struct based on `validate:"..."` tags
		if err := validator.New().Struct(student); err != nil {
			// Convert the returned error into ValidationErrors type
			// because validator returns a generic 'error' interface,
			// but we need detailed information (field, rule, message)
			validateErrs := err.(validator.ValidationErrors)
			// Send a JSON response back to the client with status 400 (Bad Request)
			// response.ValidationError(validateErrs) → converts validation errors
			// into a proper readable JSON object for the client.
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
