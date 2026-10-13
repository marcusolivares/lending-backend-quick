package postgres

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcusolivares/lending-backend-quick/internal/loans"
)

type LoanApplicationRepository struct {
	pool *pgxpool.Pool
}

func NewLoanApplicationRepository(pool *pgxpool.Pool) *LoanApplicationRepository {
	return &LoanApplicationRepository{pool: pool}
}

func (r *LoanApplicationRepository) Save(ctx context.Context, request loans.LoanApplicationRequest) (string, error) {
	if r == nil || r.pool == nil {
		return "", fmt.Errorf("database pool is not configured")
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	applicationID := newUUIDString()
	if _, err := tx.Exec(ctx, `
		INSERT INTO loan_applications (
			id,
			amount,
			tenure_months,
			status,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted,
			deleted_at
		) VALUES ($1, $2, $3, 'IN_PROGRESS', NOW(), 'system', NOW(), 'system', false, NULL)
	`, applicationID, request.Loan.Amount, request.Loan.TenureMonths); err != nil {
		return "", fmt.Errorf("insert loan application: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO customers (
			application_id,
			full_name,
			id_number,
			address_street,
			address_city,
			address_postal_code,
			birth_date,
			phone_number,
			email,
			monthly_income,
			marital_status,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted,
			deleted_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), 'system', NOW(), 'system', false, NULL)
	`, applicationID, request.Customer.FullName, request.Customer.IDNumber, request.Customer.Address.Street, request.Customer.Address.City, request.Customer.Address.PostalCode, request.Customer.BirthDate, request.Customer.PhoneNumber, request.Customer.Email, request.Customer.MonthlyIncome, request.Customer.MaritalStatus); err != nil {
		return "", fmt.Errorf("insert customer: %w", err)
	}

	if request.Customer.MaritalStatus == loans.MarriedStatus && request.Spouse != nil {
		if _, err := tx.Exec(ctx, `
			INSERT INTO spouses (
				application_id,
				full_name,
				id_card_number,
				phone_number,
				monthly_income,
				created_at,
				created_by,
				updated_at,
				updated_by,
				deleted,
				deleted_at
			) VALUES ($1, $2, $3, $4, $5, NOW(), 'system', NOW(), 'system', false, NULL)
		`, applicationID, request.Spouse.FullName, request.Spouse.IDCardNumber, request.Spouse.PhoneNumber, request.Spouse.MonthlyIncome); err != nil {
			return "", fmt.Errorf("insert spouse: %w", err)
		}
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO car_collaterals (
			application_id,
			brand,
			model,
			manufacturing_year,
			license_plate,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted,
			deleted_at
		) VALUES ($1, $2, $3, $4, $5, NOW(), 'system', NOW(), 'system', false, NULL)
	`, applicationID, request.Car.Brand, request.Car.Model, request.Car.ManufacturingYear, request.Car.LicensePlate); err != nil {
		return "", fmt.Errorf("insert car collateral: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", fmt.Errorf("commit transaction: %w", err)
	}

	return applicationID, nil
}

func newUUIDString() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	buf[6] = (buf[6] & 0x0f) | 0x40
	buf[8] = (buf[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:16])
}
