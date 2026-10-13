# Research: Loan Submission

## Go PostgreSQL Access

- **Decision**: Use `github.com/jackc/pgx/v5` and `pgxpool` with explicit SQL statements
  and positional parameters. Do not use an ORM.
- **Rationale**: Native pgx provides PostgreSQL-specific functionality and a concurrency-
  safe connection pool. Parameter binding separates untrusted request values from SQL.
  Pin a compatible released version in Go module files when bootstrapping the service.
- **Alternatives considered**: `pgx/v5/stdlib` for `database/sql` interoperability; an
  ORM (not allowed by the stated project constraint).
- **References**: [pgx v5](https://pkg.go.dev/github.com/jackc/pgx/v5),
  [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool),
  [Go dependency management](https://go.dev/doc/modules/managing-dependencies),
  [Go SQL injection guidance](https://go.dev/doc/database/sql-injection).

## Atomic Persistence

- **Decision**: Insert the application and its customer, optional spouse, and collateral
  records within one pgx transaction. Use `RETURNING` to obtain generated identifiers;
  commit only after every required insert succeeds, and roll back on every failure path.
- **Rationale**: A submission is an all-or-nothing operation; separate autocommitted
  inserts could leave partial financial/customer records.
- **Alternatives considered**: A single compound SQL statement may work for simpler
  relationships but is less direct for the multiple related records. Separate inserts
  without a transaction violate the spec's atomicity requirement.
- **References**: [PostgreSQL transactions](https://www.postgresql.org/docs/current/tutorial-transactions.html),
  [pgxpool transaction API](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#Pool.Begin),
  [Go transaction guidance](https://go.dev/doc/database/execute-transactions).

## Request Validation and Error Contract

- **Decision**: Validate request shape and business rules at the HTTP boundary, then
  enforce durable invariants again with PostgreSQL constraints. Document errors using
  `application/problem+json` with stable field-level validation details.
- **Rationale**: Boundary checks give clear feedback; database constraints remain the
  concurrency-safe integrity backstop. A stable error shape is easier for clients to
  consume and avoids leaking database internals.
- **Alternatives considered**: A project-specific JSON error envelope is also valid if
  one becomes established; the current repository has no existing API contract.
- **References**: [RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html),
  [PostgreSQL constraints](https://www.postgresql.org/docs/current/ddl-constraints.html),
  [Go JSON decoder](https://pkg.go.dev/encoding/json#Decoder.DisallowUnknownFields).

## Request Context and Database Errors

- **Decision**: Propagate `r.Context()` through the service and every pgx call, use
  bounded request/database timeouts, and map only known SQLSTATE/constraint errors to
  client-safe responses. Treat unexpected database failures as server errors.
- **Rationale**: Cancellation bounds wasted database work. Typed error inspection avoids
  brittle message matching and prevents internal SQL details from reaching clients.
- **Alternatives considered**: Unbounded database calls or string-matching error text
  are less reliable and harder to operate.
- **References**: [Go cancellation guidance](https://go.dev/doc/database/cancel-operations),
  [PostgreSQL SQLSTATE codes](https://www.postgresql.org/docs/current/errcodes-appendix.html),
  [pgx PostgreSQL errors](https://pkg.go.dev/github.com/jackc/pgx/v5/pgconn#PgError).

## Migrations and Testing

- **Decision**: Maintain ordered, versioned plain-SQL migrations and use Goose as the
  migration runner. Test validation and HTTP mapping without a database; run migration,
  constraint, and transaction-rollback integration tests against PostgreSQL.
- **Rationale**: SQL migrations stay reviewable and meet the constitution's requirement
  for versioned schema changes. A real PostgreSQL instance is necessary to verify
  database semantics and transaction behavior.
- **Alternatives considered**: golang-migrate is another established SQL migration
  runner. A shared PostgreSQL test database is possible in an existing CI environment;
  a disposable container can improve isolation where container support is available.
- **References**: [Goose](https://github.com/pressly/goose),
  [Go testing](https://pkg.go.dev/testing),
  [Testcontainers PostgreSQL module](https://golang.testcontainers.org/modules/postgres/).
