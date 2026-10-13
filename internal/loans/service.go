package loans

import (
	"context"
	"fmt"
)

type Repository interface {
	Save(ctx context.Context, request LoanApplicationRequest) (string, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Submit(ctx context.Context, request LoanApplicationRequest) (LoanConfirmation, error) {
	if err := ValidateSubmission(request); err != nil {
		return LoanConfirmation{}, err
	}

	applicationID, err := s.repository.Save(ctx, request)
	if err != nil {
		return LoanConfirmation{}, fmt.Errorf("save loan application: %w", err)
	}

	return LoanConfirmation{
		ApplicationID: applicationID,
		Status:        StatusInProgress,
	}, nil
}
