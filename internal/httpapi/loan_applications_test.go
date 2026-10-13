package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/marcusolivares/lending-backend-quick/internal/loans"
)

type fakeService struct {
	confirmation loans.LoanConfirmation
	err          error
}

func (f fakeService) Submit(_ context.Context, _ loans.LoanApplicationRequest) (loans.LoanConfirmation, error) {
	if f.err != nil {
		return loans.LoanConfirmation{}, f.err
	}
	return f.confirmation, nil
}

func TestLoanApplicationHandler_Success(t *testing.T) {
	service := fakeService{confirmation: loans.LoanConfirmation{ApplicationID: "abc-123", Status: loans.StatusInProgress}}
	router := NewRouter(service)

	payload := `{"customer":{"full_name":"Alice Example","id_number":"ABC123","address":{"street":"123 Main","city":"Metropolis","postal_code":"1010"},"birth_date":"1990-06-15","phone_number":"+123456789","email":"alice@example.com","monthly_income":5000,"marital_status":"UNMARRIED"},"car":{"brand":"Toyota","model":"Corolla","manufacturing_year":2022,"license_plate":"ABC-123"},"loan":{"amount":150000,"tenure_months":24}}`

	req := httptest.NewRequest(http.MethodPost, "/v1/loan-applications", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", res.Code, res.Body.String())
	}

	var confirmation loans.LoanConfirmation
	if err := json.Unmarshal(res.Body.Bytes(), &confirmation); err != nil {
		t.Fatalf("decode confirmation: %v", err)
	}
	if confirmation.Status != loans.StatusInProgress {
		t.Fatalf("expected status %q, got %q", loans.StatusInProgress, confirmation.Status)
	}
}

func TestLoanApplicationHandler_ValidationFailure(t *testing.T) {
	service := fakeService{err: loans.ValidationErrors{Errors: []loans.FieldError{{Field: "customer.email", Message: "must be a syntactically valid email address"}}}}
	router := NewRouter(service)

	payload := `{"customer":{"full_name":"Alice Example","id_number":"ABC123","address":{"street":"123 Main","city":"Metropolis","postal_code":"1010"},"birth_date":"1990-06-15","phone_number":"+123456789","email":"bad-email","monthly_income":5000,"marital_status":"UNMARRIED"},"car":{"brand":"Toyota","model":"Corolla","manufacturing_year":2022,"license_plate":"ABC-123"},"loan":{"amount":150000,"tenure_months":24}}`

	req := httptest.NewRequest(http.MethodPost, "/v1/loan-applications", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d: %s", res.Code, res.Body.String())
	}
	if !strings.Contains(res.Body.String(), "customer.email") {
		t.Fatalf("expected validation error payload to include field name: %s", res.Body.String())
	}
}

func TestLoanApplicationHandler_MalformedJSON(t *testing.T) {
	router := NewRouter(fakeService{})
	req := httptest.NewRequest(http.MethodPost, "/v1/loan-applications", strings.NewReader("{bad json}"))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", res.Code, res.Body.String())
	}
}

func TestLoanApplicationHandler_UnexpectedFailure(t *testing.T) {
	router := NewRouter(fakeService{err: errors.New("db timeout")})
	req := httptest.NewRequest(http.MethodPost, "/v1/loan-applications", strings.NewReader(`{"customer":{"full_name":"Alice","id_number":"A","address":{"street":"1","city":"X","postal_code":"1"},"birth_date":"1990-01-01","phone_number":"1","email":"a@example.com","monthly_income":10,"marital_status":"UNMARRIED"},"car":{"brand":"A","model":"B","manufacturing_year":2024,"license_plate":"C"},"loan":{"amount":100,"tenure_months":12}}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", res.Code, res.Body.String())
	}
}
