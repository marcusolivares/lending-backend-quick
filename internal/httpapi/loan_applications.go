package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/marcusolivares/lending-backend-quick/internal/loans"
)

type LoanApplicationHandler struct {
	service ApplicationService
}

func NewLoanApplicationHandler(service ApplicationService) *LoanApplicationHandler {
	return &LoanApplicationHandler{service: service}
}

func (h *LoanApplicationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteProblem(w, http.StatusMethodNotAllowed, "https://example.com/problems/method-not-allowed", "Method not allowed", "Use POST to submit a loan application", nil)
		return
	}

	var req loans.LoanApplicationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		WriteProblem(w, http.StatusBadRequest, "https://example.com/problems/invalid-json", "Malformed JSON", "The request body is not valid JSON", nil)
		return
	}

	if r.Header.Get("Content-Type") != "" && !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		WriteProblem(w, http.StatusUnsupportedMediaType, "https://example.com/problems/unsupported-media-type", "Unsupported content type", "Content-Type must be application/json", nil)
		return
	}

	confirmation, err := h.service.Submit(r.Context(), req)
	if err != nil {
		var validationErr loans.ValidationErrors
		if errors.As(err, &validationErr) {
			WriteProblem(w, http.StatusUnprocessableEntity, "https://example.com/problems/validation-error", "Validation failed", "One or more fields are invalid", validationErr.Errors)
			return
		}
		WriteProblem(w, http.StatusInternalServerError, "https://example.com/problems/internal-error", "Internal server error", "The loan application could not be processed", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(confirmation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *Router) handleLoanApplication(w http.ResponseWriter, req *http.Request) {
	NewLoanApplicationHandler(r.service).ServeHTTP(w, req)
}
