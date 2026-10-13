package loans

import "testing"

func TestValidateSubmission_ValidMarriedRequest(t *testing.T) {
	req := validMarriedRequest()
	if err := ValidateSubmission(req); err != nil {
		t.Fatalf("expected valid married request to pass: %v", err)
	}
}

func TestValidateSubmission_ValidUnmarriedRequest(t *testing.T) {
	req := validUnmarriedRequest()
	if err := ValidateSubmission(req); err != nil {
		t.Fatalf("expected valid unmarried request to pass: %v", err)
	}
}

func TestValidateSubmission_RejectsMissingSpouseForMarried(t *testing.T) {
	req := validMarriedRequest()
	req.Spouse = nil
	if err := ValidateSubmission(req); err == nil {
		t.Fatal("expected married submission to require spouse")
	}
}

func TestValidateSubmission_RejectsSpouseForUnmarried(t *testing.T) {
	req := validUnmarriedRequest()
	req.Spouse = &Spouse{
		FullName:      "Jane Doe",
		IDCardNumber:  "12345",
		PhoneNumber:   "555-0000",
		MonthlyIncome: 4000,
	}
	if err := ValidateSubmission(req); err == nil {
		t.Fatal("expected unmarried submission to reject spouse details")
	}
}

func TestValidateSubmission_RejectsInvalidFields(t *testing.T) {
	req := validUnmarriedRequest()
	req.Customer.Email = "not-an-email"
	req.Car.ManufacturingYear = timeNowPlusOneYear()
	req.Loan.Amount = 250
	req.Customer.MonthlyIncome = -1

	if err := ValidateSubmission(req); err == nil {
		t.Fatal("expected validation failure for invalid request")
	}
}

func validMarriedRequest() LoanApplicationRequest {
	return LoanApplicationRequest{
		Customer: Customer{
			FullName:      "Alice Example",
			IDNumber:      "ABC123",
			Address:       Address{Street: "123 Main", City: "Metropolis", PostalCode: "1010"},
			BirthDate:     "1990-06-15",
			PhoneNumber:   "+123456789",
			Email:         "alice@example.com",
			MonthlyIncome: 5000,
			MaritalStatus: MarriedStatus,
		},
		Spouse: &Spouse{
			FullName:      "Bob Example",
			IDCardNumber:  "XYZ456",
			PhoneNumber:   "+123456780",
			MonthlyIncome: 4000,
		},
		Car: CarCollateral{
			Brand:             "Toyota",
			Model:             "Corolla",
			ManufacturingYear: 2022,
			LicensePlate:      "ABC-123",
		},
		Loan: Loan{Amount: 150000, TenureMonths: 24},
	}
}

func validUnmarriedRequest() LoanApplicationRequest {
	req := validMarriedRequest()
	req.Customer.MaritalStatus = UnmarriedStatus
	req.Spouse = nil
	return req
}

func timeNowPlusOneYear() int {
	return 2100
}
