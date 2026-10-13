# Implementation Plan: Loan Submission

**Branch**: `001-submit-loan` | **Date**: 2026-10-13 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/001-submit-loan/spec.md`

## Summary

Implement a Go REST API that validates and records a complete loan application in
PostgreSQL, including the customer, the spouse when the customer is married, and car
collateral. Use `pgx` with explicit parameterized SQL and a single database transaction
for all related inserts. Do not use an ORM.

## Technical Context

**Language/Version**: Go; choose and pin a currently supported version compatible with
the selected pgx v5 release in `go.mod`.

**Primary Dependencies**: `github.com/jackc/pgx/v5` and `pgxpool`; no ORM. Use `goose`
with ordered plain-SQL migration files for schema changes.

**Storage**: PostgreSQL, authoritative system of record.

**Testing**: Go `testing` package; unit and HTTP contract tests run with `go test ./...`.
Persistence and migration integration tests run against PostgreSQL.

**Target Platform**: Go server process; deployment operating system and packaging are
not specified by the repository.

**Project Type**: REST web service.

**Performance Goals**: No request-latency or throughput targets are specified. Avoid
unnecessary external calls in submission handling; establish operational targets
separately if required.

**Constraints**: Resource-oriented REST contract; explicit request validation; plain
parameterized SQL through pgx; no ORM; all related records inserted atomically;
PostgreSQL constraints enforce persistent invariants; every table includes the mandated
audit and soft-delete columns; logs exclude customer and financial data.

**Scale/Scope**: One loan-submission endpoint and its persistence model. Traffic and
data-volume targets are not specified.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Plan |
|---|---|---|
| Contract-First REST API | PASS | Define a versioned resource endpoint, request/response schemas, status codes, and field validation errors in the OpenAPI contract; any timestamps later exposed by this API must follow the project format. |
| Idiomatic Go | PASS | Use focused Go packages, context propagation, explicit error handling, and pgx v5 without an ORM. |
| PostgreSQL Integrity | PASS | Use parameterized SQL, constraints, versioned SQL migrations, required audit/soft-delete columns, and a transaction for the application and related records. |
| Automated Testing | PASS | Add unit and HTTP contract tests plus PostgreSQL integration tests for constraints, migrations, and rollback behavior. |
| Secure and Observable Operations | PASS | Validate at the API boundary, rely on established authorization, and avoid logging credentials, PII, or financial values. |

No constitution violations or complexity exceptions are proposed.

## Project Structure

### Documentation (this feature)

```text
specs/001-submit-loan/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── loan-submission.openapi.yaml
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
└── api/
    └── main.go
internal/
├── config/
├── httpapi/
│   └── loan_applications.go
├── loans/
│   ├── model.go
│   ├── service.go
│   └── validation.go
└── postgres/
    └── loan_applications.go
migrations/
tests/
└── integration/
```

**Structure Decision**: The repository currently contains Spec Kit configuration and
feature documentation, but no application source or existing Go project layout. Start
with one Go API service using `cmd/api` for the entry point and focused `internal`
packages for configuration, HTTP handling, loan rules, and PostgreSQL persistence.
Keep versioned SQL migrations at the repository root. Unit tests should live alongside
their Go packages; database integration tests belong under `tests/integration`.

## Complexity Tracking

No constitution violations require justification.
