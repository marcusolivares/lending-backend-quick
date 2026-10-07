# Feature Specification: Loan Submission

**Feature Branch**: `001-submit-loan`

**Created**: 2026-10-07

**Status**: Draft

**Input**: User description: "Create a REST API for submitting a loan. The API will need customer data (full name, ID number, address, birth date, phone number, email, and current monthly income); car collateral data (brand, model, manufacturing year, and license plate); and proposed loan data (amount, a multiple of 100, and tenure between 3-48 months). If all inputs are valid, store the data in the database as IN_PROGRESS status."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Submit a Loan Application (Priority: P1)

A customer submits a loan application with their identity and contact information,
monthly income, car collateral details, and proposed loan amount and tenure. The service
validates the complete submission and, when valid, records the application with
`IN_PROGRESS` status and confirms the submission.

**Why this priority**: This is the primary value of the feature: capture a complete,
eligible-to-review loan application for subsequent processing.

**Independent Test**: Submit one request with valid customer, car, and loan details;
verify the confirmation identifies the newly recorded application and its `IN_PROGRESS`
status. Submit invalid requests and verify they are rejected without creating an
application.

**Acceptance Scenarios**:

1. **Given** all required customer, car, and proposed loan fields are valid, **When** the
   customer submits the application, **Then** the service records all submitted data,
   assigns `IN_PROGRESS` status, and returns a successful confirmation.
2. **Given** one or more required fields are missing or invalid, **When** the customer
   submits the application, **Then** the service rejects it with a clear validation
   response and creates no application.
3. **Given** a valid loan amount and tenure at their permitted boundaries, **When** the
   customer submits the application, **Then** the service accepts the values and records
   the application with `IN_PROGRESS` status.

### Edge Cases

- Reject a loan amount that is zero, negative, or not divisible by 100.
- Reject a tenure below 3 or above 48 months; accept both endpoints.
- Reject a malformed or impossible birth date, malformed email, empty required text,
  invalid monthly income, and malformed car manufacturing year.
- Reject a car manufacturing year later than the current calendar year.
- A validation failure in any section rejects the complete submission; partial
  application data MUST NOT be stored.
- If storage fails, do not return a successful submission confirmation.
- If a response includes a timestamp, it uses `YYYY-MM-DDTHH:MM:SS±HH:MM` with seconds
  and a numeric timezone offset.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The service MUST accept a loan submission containing customer full name,
  ID number, address (street, city, and ZIP/postal code), birth date, phone number,
  email, and current monthly income.
- **FR-002**: The service MUST accept car collateral details containing brand, model,
  manufacturing year, and license plate.
- **FR-003**: The service MUST accept a proposed loan amount and tenure in months.
- **FR-004**: The service MUST require every field in the customer, collateral, and
  proposed loan sections.
- **FR-005**: The service MUST validate birth date as a real calendar date in
  `YYYY-MM-DD` format, email as a syntactically valid email address, and the car
  manufacturing year as a valid year no later than the current calendar year.
- **FR-006**: The service MUST require monthly income to be a non-negative monetary
  amount and loan amount to be a positive monetary amount divisible by 100.
- **FR-007**: The service MUST require loan tenure to be a whole number from 3 through
  48 months, inclusive.
- **FR-008**: The service MUST reject submissions with any missing or invalid field,
  identify the invalid field or fields, and MUST NOT persist a partial application.
- **FR-009**: For each valid submission, the service MUST persist all submitted
  customer, car collateral, and proposed loan data as one application with status
  `IN_PROGRESS`.
- **FR-010**: After successful persistence, the service MUST return a confirmation that
  identifies the application and states its `IN_PROGRESS` status.
- **FR-011**: The service MUST NOT return a success confirmation when it cannot persist
  the application.
- **FR-012**: Any timestamp accepted or returned by the API MUST use
  `YYYY-MM-DDTHH:MM:SS±HH:MM` (for example, `2026-05-05T14:30:00+07:00`). The customer
  birth date is a calendar date, not a timestamp.
- **FR-013**: Persisted loan application data MUST support audit attribution and
  soft-deletion in accordance with project data governance.

### Key Entities

- **Customer**: The applicant, represented by full name, ID number, address, birth date,
  phone number, email, and current monthly income.
- **Car Collateral**: The vehicle securing the application, represented by brand, model,
  manufacturing year, and license plate.
- **Loan Application**: The submitted proposal linking customer information and car
  collateral to a loan amount and tenure; it has an application identifier and status.
  A successfully submitted application begins in `IN_PROGRESS`.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of submissions containing valid required data and within the loan
  amount and tenure limits are stored as one application with `IN_PROGRESS` status.
- **SC-002**: 100% of submissions containing missing or invalid data are rejected and
  create no application or partial application data.
- **SC-003**: 100% of successful submissions return a confirmation identifying the
  application and its status.
- **SC-004**: Customers receive field-specific validation feedback for every invalid
  field in a rejected submission.
- **SC-005**: Every timestamp present in an API request or response includes seconds
  and an explicit numeric timezone offset in the required format.

## Assumptions

- All listed customer, car, and proposed loan fields are required.
- Loan amount is strictly positive and monthly income may be zero; both are represented
  as monetary values in the application's configured currency.
- ZIP/postal code and phone number are required non-empty strings. Their detailed
  validation follows the application's supported customer locale rather than assuming
  one country-specific format.
- Birth date is a calendar date formatted as `YYYY-MM-DD`; no minimum age eligibility
  rule was specified.
- Email validation checks syntax only; verifying ownership or deliverability is outside
  this feature.
- A successful submission is recorded once per accepted request. Duplicate detection,
  credit assessment, approval, and later status transitions are outside this feature.
- Existing authentication, authorization, and privacy policies continue to apply.
- Database-created or updated timestamps, if exposed through the API, comply with the
  project's timestamp format rule.
