# Quickstart: Loan Submission

This guide describes the validation runs for the loan submission feature. It assumes
the Go service, PostgreSQL database, and SQL migrations have been set up according to
the project deployment configuration.

## Prerequisites

- A currently supported Go toolchain compatible with the version pinned in `go.mod`.
- PostgreSQL reachable through `DATABASE_URL`.
- Goose installed for applying the versioned SQL migrations.

## Run checks

From the repository root:

```powershell
go test ./...
```

Expected result: unit and HTTP contract tests pass. PostgreSQL integration tests also
pass when their database configuration is available.

Apply migrations to a disposable or designated development database:

```powershell
goose -dir migrations postgres "$env:DATABASE_URL" up
```

Start the service (default local address `http://localhost:8080`):

```powershell
$env:DATABASE_URL = "postgres://user:password@localhost:5432/lending?sslmode=disable"
go run ./cmd/api
```

The API listens at `http://localhost:8080` unless deployment configuration overrides
the address.

## End-to-end validation scenarios

1. Submit a valid unmarried applicant with customer, collateral, and loan data, with no
   `spouse` property. Expect `201 Created`, an application identifier, and
   `IN_PROGRESS`.
2. Submit a valid married applicant with all four spouse fields. Expect `201 Created`
   and persisted application, customer, spouse, and collateral data.
3. Omit a required customer, collateral, loan, or married-applicant spouse field.
   Expect a field-specific `422 Unprocessable Entity` response and no persisted rows.
4. Include spouse details for an unmarried applicant. Expect a field-specific
   `422 Unprocessable Entity` response and no persisted rows.
5. Submit a car year that is not four digits or is in the future; an empty car brand,
   model, or license plate; a loan amount that is non-positive or not divisible by
   100; or a tenure outside 3–48 months. Expect validation feedback and no persistence.
6. Force a database failure after an earlier insert in the transaction. Expect a
   server error and verify that the transaction leaves no partial application data.

See [the OpenAPI contract](./contracts/loan-submission.openapi.yaml) for request and
response schemas, and [the data model](./data-model.md) for persistence relationships
and invariants.
