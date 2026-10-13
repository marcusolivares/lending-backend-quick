package loans

import "time"

const (
	MarriedStatus    = "MARRIED"
	UnmarriedStatus  = "UNMARRIED"
	StatusInProgress = "IN_PROGRESS"
)

type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
}

type Customer struct {
	FullName      string  `json:"full_name"`
	IDNumber      string  `json:"id_number"`
	Address       Address `json:"address"`
	BirthDate     string  `json:"birth_date"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	MonthlyIncome float64 `json:"monthly_income"`
	MaritalStatus string  `json:"marital_status"`
}

type Spouse struct {
	FullName      string  `json:"full_name"`
	IDCardNumber  string  `json:"id_card_number"`
	PhoneNumber   string  `json:"phone_number"`
	MonthlyIncome float64 `json:"monthly_income"`
}

type CarCollateral struct {
	Brand             string `json:"brand"`
	Model             string `json:"model"`
	ManufacturingYear int    `json:"manufacturing_year"`
	LicensePlate      string `json:"license_plate"`
}

type Loan struct {
	Amount       float64 `json:"amount"`
	TenureMonths int     `json:"tenure_months"`
}

type LoanApplicationRequest struct {
	Customer  Customer      `json:"customer"`
	Spouse    *Spouse       `json:"spouse,omitempty"`
	Car       CarCollateral `json:"car"`
	Loan      Loan          `json:"loan"`
	CreatedAt time.Time     `json:"-"`
}

type LoanConfirmation struct {
	ApplicationID string `json:"application_id"`
	Status        string `json:"status"`
}
