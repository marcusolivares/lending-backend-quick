# Tasks: Loan Submission

**Input**: Design documents from `/specs/001-submit-loan/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/](./contracts/), [quickstart.md](./quickstart.md)

**Tests**: Required by the project constitution. Add unit, HTTP contract, and PostgreSQL integration tests for the specified behavior.

**Organization**: Tasks are grouped by the sole user story in `spec.md` so the MVP can be implemented and validated independently.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel with other marked tasks because it uses different files and has no dependency on incomplete tasks.
- **[Story]**: User story label for traceability.
- Every task includes the exact file path(s) to create or change.

## Phase 1: Setup

**Purpose**: Initialize the Go module and repository-level development configuration.

- [ ] T001 Initialize the Go module as `github.com/marcusolivares/lending-backend-quick` in `go.mod`, pin a supported Go version and compatible `github.com/jackc/pgx/v5` and `github.com/pressly/goose/v3` dependencies, and generate `go.sum`.
- [ ] T002 [P] Add Go build, test, and local environment artifacts to the root `.gitignore`.

## Phase 2: Foundational

**Purpose**: Provide configuration, database connectivity, and shared HTTP infrastructure required by the loan submission story.

- [ ] T003 Implement required environment configuration loading and validation for `DATABASE_URL` and the HTTP listen address in `internal/config/config.go`.
- [ ] T004 Implement the PostgreSQL `pgxpool` lifecycle, connectivity check, context-aware close, and explicit error propagation in `internal/postgres/pool.go`.
- [ ] T005 Implement a reusable RFC 9457 problem-details response type and JSON writer with safe client-visible messages in `internal/httpapi/problem.go`.
- [ ] T006 Implement the Go HTTP server entry point in `cmd/api/main.go`, wiring configuration, the PostgreSQL pool, router, server timeouts, and graceful shutdown.
- [ ] T007 Implement the shared HTTP router and health endpoint in `internal/httpapi/router.go`, keeping route registration separate from server startup.

**Checkpoint**: Go module, configuration, PostgreSQL pool, problem responses, and a runnable HTTP server foundation are ready for the loan-submission story.

## Phase 3: User Story 1 - Submit a Loan Application (Priority: P1) — MVP

**Goal**: Accept and validate complete applications from married or unmarried customers and atomically persist the customer, applicable spouse, collateral, and loan with `IN_PROGRESS` status.

**Independent Test**: Against PostgreSQL with migrations applied, submit one valid married request and one valid unmarried request and verify each returns `201 Created`, an application identifier, and `IN_PROGRESS`, with the expected related rows. Submit invalid and conditional-spouse requests and verify field-specific errors and no persisted partial data.

- [ ] T008 [P] [US1] Define Go domain/request/confirmation types for application, customer, spouse, car collateral, and proposed loan in `internal/loans/model.go`, preserving all fields and conditional spouse relationship in `data-model.md`.

### Tests for User Story 1

> Implement the tests before validation, persistence, and endpoint logic; verify expected failures after T008 defines the shared types.

- [ ] T009 [P] [US1] Add table-driven validation tests for required non-empty customer, spouse, and collateral strings; real `YYYY-MM-DD` customer birth dates; syntactically valid email; non-negative customer/spouse income; valid four-digit non-future car year; positive loan amount divisible by 100; and tenure 3–48, including the required timestamp format rule for any returned or accepted timestamps in `internal/loans/validation_test.go`.
- [ ] T010 [P] [US1] Add HTTP contract tests for married and unmarried success, missing/invalid fields, spouse required only for `MARRIED`, spouse forbidden for `UNMARRIED`, `201` confirmation shape, returned `IN_PROGRESS` status, RFC 9457 problem details, field-specific validation errors, and required timestamp formatting in `internal/httpapi/loan_applications_test.go`.
- [ ] T011 [P] [US1] Add PostgreSQL integration tests using a dedicated test database configuration to apply migrations, persist married and unmarried submissions, enforce database constraints, verify audit and soft-delete columns, and prove a later insert failure rolls back all earlier inserts in `tests/integration/loan_applications_test.go`.

### Implementation for User Story 1

- [ ] T012 [US1] Implement request validation in `internal/loans/validation.go`: every required string must be non-empty; customer birth date must be a real date in `YYYY-MM-DD`; email must be syntactically valid; customer and applicable spouse monthly income must be non-negative; car brand, model, and license plate must be non-empty; manufacturing year must be exactly four digits and not later than the current calendar year; loan amount must be positive and divisible by 100; tenure must be a whole number from 3 through 48; spouse fields are required only for married applicants and forbidden for unmarried applicants.
- [ ] T013 [P] [US1] Create the ordered Goose SQL migration `migrations/00001_create_loan_applications.sql` for application, customer, spouse, and car-collateral tables with foreign keys, required non-null/check constraints, exact-decimal monetary columns, application status default/check, and `created_at`, `created_by`, `updated_at`, `updated_by`, `deleted`, and `deleted_at` on every table.
- [ ] T014 [US1] Implement loan persistence using explicit parameterized pgx SQL in `internal/postgres/loan_applications.go`; insert the application, customer, applicable spouse, and collateral through one transaction, use `RETURNING` for the application ID, roll back every failed insert, and return success only after commit.
- [ ] T015 [US1] Implement the loan-application service in `internal/loans/service.go` to validate the submission, call the persistence boundary, and propagate request context and errors without success-shaped fallbacks.
- [ ] T016 [US1] Implement `POST /v1/loan-applications` in `internal/httpapi/loan_applications.go` to decode the OpenAPI request shape, reject malformed/unsupported JSON, return field-specific `422` errors for validation failures, return `201` with application ID and `IN_PROGRESS` after commit, and return safe problem details on persistence failure.
- [ ] T017 [US1] Register `POST /v1/loan-applications` in `internal/httpapi/router.go` and wire the concrete pgx-backed service in `cmd/api/main.go`.
- [ ] T018 [US1] Align `specs/001-submit-loan/contracts/loan-submission.openapi.yaml` with the implemented request fields, married/unmarried spouse conditional, response status and shape, RFC 9457 problem details, required timestamp formatting, and validation/error responses.

**Checkpoint**: User Story 1 is independently functional. Run its validation and HTTP tests, then its PostgreSQL integration tests against the migrated test database.

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Verify the complete feature and keep its run guide consistent with the implementation.

- [ ] T019 [P] Review `internal/httpapi/loan_applications.go`, `internal/loans/service.go`, and `internal/postgres/loan_applications.go` to ensure request data, PII, and financial values are not written to logs, database errors do not expose SQL or driver details to clients, and all persisted tables include the required audit/soft-delete columns and timestamp semantics required by the project constitution.
- [ ] T020 Update `specs/001-submit-loan/quickstart.md` with the actual configuration names, migration command, server address, test prerequisites, expected outcomes, and timestamp-format expectations implemented by `cmd/api/main.go` and `migrations/`.
- [ ] T021 Run `go test ./...`, `go vet ./...`, Goose migration up/down validation against a disposable PostgreSQL database, and all end-to-end scenarios in `specs/001-submit-loan/quickstart.md`, including timestamp-format and audit-field verification; fix failures in the corresponding Go, migration, contract, or documentation files.

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies; T001 and T002 can start immediately.
- **Foundational (Phase 2)**: Depends on T001; complete T003–T007 before the story endpoint is wired.
- **User Story 1 (Phase 3)**: Depends on Phase 2. T008 establishes the shared types; T009–T011 are test-first tasks and may be authored in parallel afterward. T008 and T013 may be developed in parallel; validation depends on T008; persistence depends on T008 and T013; service depends on validation and persistence; handler depends on shared HTTP errors and the service; route wiring follows the handler.
- **Polish (Phase 4)**: Depends on User Story 1 implementation; T019 and T020 can proceed in parallel, then complete T021.

### User Story Dependencies

- **User Story 1 (P1)**: No dependency on another story; it is the only story and is independently testable after foundational work.

### Within User Story 1

- Implement T008 before writing the story tests. Write T009–T011 before validation, persistence, and endpoint logic, and confirm they fail for missing behavior.
- Implement shared request/domain types before validation and service use.
- Create the SQL schema before database integration execution and repository completion.
- Complete validation and transactional persistence before service and HTTP wiring.
- Complete route wiring before end-to-end quickstart validation.

## Parallel Example: User Story 1

```text
# After T008, author these independent test suites together:
T009 validation tests in internal/loans/validation_test.go
T010 HTTP contract tests in internal/httpapi/loan_applications_test.go
T011 PostgreSQL integration tests in tests/integration/loan_applications_test.go

# Once test contracts are agreed, these different-file tasks can proceed in parallel:
T013 initial SQL schema in migrations/00001_create_loan_applications.sql
```

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 setup.
2. Complete Phase 2 shared foundation.
3. Author and run the User Story 1 tests.
4. Implement the data model, validation, migration, transactional persistence, service, and endpoint.
5. **Stop and validate** the married and unmarried submission flows independently using PostgreSQL.
6. Complete cross-cutting privacy/error review and quickstart validation.

### Incremental Delivery

There is one user story. Deliver the setup and foundation first, then the complete P1
submission path as the MVP. Do not split conditional spouse persistence or validation
from that path because both are required for a valid submission.

## Notes

- Every task uses the required checkbox, sequential task ID, optional parallel marker, required user-story label for story tasks, and exact file path(s).
- `[P]` tasks are limited to independent work in different files.
- All database operations use plain parameterized SQL through pgx; no ORM is permitted.
- Do not add ID-number/license-plate uniqueness, country-specific formats, or unstated age rules.
