package httpapi

import (
	"context"
	"net/http"

	"github.com/marcusolivares/lending-backend-quick/internal/loans"
)

type ApplicationService interface {
	Submit(ctx context.Context, request loans.LoanApplicationRequest) (loans.LoanConfirmation, error)
}

type Router struct {
	service ApplicationService
}

func NewRouter(service ApplicationService) http.Handler {
	r := &Router{service: service}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", r.healthHandler)
	mux.HandleFunc("/v1/loan-applications", r.handleLoanApplication)
	return mux
}

func (r *Router) healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
