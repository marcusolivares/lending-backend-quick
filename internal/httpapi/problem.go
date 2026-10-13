package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/marcusolivares/lending-backend-quick/internal/loans"
)

type ProblemDetails struct {
	Type   string             `json:"type"`
	Title  string             `json:"title"`
	Status int                `json:"status"`
	Detail string             `json:"detail"`
	Errors []loans.FieldError `json:"errors,omitempty"`
}

func WriteProblem(w http.ResponseWriter, statusCode int, problemType, title, detail string, errs []loans.FieldError) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	payload := ProblemDetails{
		Type:   problemType,
		Title:  title,
		Status: statusCode,
		Detail: detail,
		Errors: errs,
	}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode problem response", http.StatusInternalServerError)
	}
}
