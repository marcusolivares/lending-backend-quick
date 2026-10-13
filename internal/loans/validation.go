package loans

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// FieldError captures a specific field-level validation issue.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors groups multiple field-level validation failures.
type ValidationErrors struct {
	Errors []FieldError
}

func (ve ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}
	parts := make([]string, 0, len(ve.Errors))
	for _, err := range ve.Errors {
		parts = append(parts, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(parts, "; ")
}

func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, FieldError{Field: field, Message: message})
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

func ValidateSubmission(req LoanApplicationRequest) error {
	var errs ValidationErrors

	validateRequiredString(&errs, "customer.full_name", req.Customer.FullName)
	validateRequiredString(&errs, "customer.id_number", req.Customer.IDNumber)
	validateRequiredString(&errs, "customer.address.street", req.Customer.Address.Street)
	validateRequiredString(&errs, "customer.address.city", req.Customer.Address.City)
	validateRequiredString(&errs, "customer.address.postal_code", req.Customer.Address.PostalCode)
	validateRequiredString(&errs, "customer.phone_number", req.Customer.PhoneNumber)
	validateRequiredString(&errs, "customer.email", req.Customer.Email)
	validateRequiredString(&errs, "customer.marital_status", req.Customer.MaritalStatus)
	validateRequiredString(&errs, "car.brand", req.Car.Brand)
	validateRequiredString(&errs, "car.model", req.Car.Model)
	validateRequiredString(&errs, "car.license_plate", req.Car.LicensePlate)

	if req.Customer.MaritalStatus != MarriedStatus && req.Customer.MaritalStatus != UnmarriedStatus {
		errs.Add("customer.marital_status", "must be either MARRIED or UNMARRIED")
	}

	if _, err := time.Parse("2006-01-02", req.Customer.BirthDate); err != nil {
		errs.Add("customer.birth_date", "must be a real date in YYYY-MM-DD format")
	}
	if _, err := mail.ParseAddress(req.Customer.Email); err != nil {
		errs.Add("customer.email", "must be a syntactically valid email address")
	}

	if req.Customer.MonthlyIncome < 0 {
		errs.Add("customer.monthly_income", "must be non-negative")
	}
	if req.Customer.MaritalStatus == MarriedStatus {
		if req.Spouse == nil {
			errs.Add("spouse", "is required for married customers")
		} else {
			validateRequiredString(&errs, "spouse.full_name", req.Spouse.FullName)
			validateRequiredString(&errs, "spouse.id_card_number", req.Spouse.IDCardNumber)
			validateRequiredString(&errs, "spouse.phone_number", req.Spouse.PhoneNumber)
			if req.Spouse.MonthlyIncome < 0 {
				errs.Add("spouse.monthly_income", "must be non-negative")
			}
		}
	} else if req.Spouse != nil {
		errs.Add("spouse", "must be omitted for unmarried customers")
	}

	if req.Car.ManufacturingYear < 1000 || req.Car.ManufacturingYear > 9999 {
		errs.Add("car.manufacturing_year", "must be exactly four digits")
	} else if req.Car.ManufacturingYear > time.Now().Year() {
		errs.Add("car.manufacturing_year", "must not be later than the current calendar year")
	}

	if req.Loan.Amount <= 0 || int64(req.Loan.Amount)%100 != 0 {
		errs.Add("loan.amount", "must be a positive value divisible by 100")
	}
	if req.Loan.TenureMonths < 3 || req.Loan.TenureMonths > 48 {
		errs.Add("loan.tenure_months", "must be between 3 and 48 months inclusive")
	}

	if req.Customer.MaritalStatus == MarriedStatus && req.Spouse != nil {
		if req.Spouse.FullName == "" {
			errs.Add("spouse.full_name", "is required for married customers")
		}
		if req.Spouse.IDCardNumber == "" {
			errs.Add("spouse.id_card_number", "is required for married customers")
		}
		if req.Spouse.PhoneNumber == "" {
			errs.Add("spouse.phone_number", "is required for married customers")
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

func validateRequiredString(errs *ValidationErrors, field, value string) {
	if strings.TrimSpace(value) == "" {
		errs.Add(field, "is required")
	}
}

func ParseValidationErrors(err error) ([]FieldError, bool) {
	var validationErr ValidationErrors
	if errors.As(err, &validationErr) {
		return validationErr.Errors, true
	}
	return nil, false
}

func ValidateTimestamp(value string) error {
	if value == "" {
		return errors.New("timestamp is required")
	}
	_, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return nil
	}
	_, err = time.Parse("2006-01-02T15:04:05-07:00", value)
	if err == nil {
		return nil
	}
	return fmt.Errorf("timestamp must use RFC3339 format with timezone offset")
}

func isPlainEmail(value string) bool {
	return emailPattern.MatchString(value)
}

func isValidYear(value int) bool {
	return value >= 1000 && value <= 9999 && value <= time.Now().Year()
}

func ParseInt(value string) (int, error) {
	return strconv.Atoi(value)
}
