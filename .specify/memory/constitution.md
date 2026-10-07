<!--
Sync Impact Report
Version change: 1.1.0 -> 1.2.0
Modified principles: I. Contract-First REST API (expanded to require a timestamp format)
Added sections: none
Removed sections: none
Follow-up TODOs: Confirm the original ratification date and replace TODO(RATIFICATION_DATE).
-->

# Lending Backend Constitution

## Core Principles

### I. Contract-First REST API
HTTP endpoints MUST use resource-oriented REST conventions, explicit request and response
schemas, and appropriate HTTP methods and status codes. Changes to externally consumed
contracts MUST be deliberate, documented, and backward-compatible unless an approved
breaking change includes a migration path. Invalid requests MUST receive consistent,
actionable error responses without exposing internal implementation details. Every
timestamp accepted or returned by the API MUST use the format
`YYYY-MM-DDTHH:MM:SS±HH:MM`, including seconds and an explicit timezone offset, for
example `2026-05-05T14:30:00+07:00`. This keeps client behavior predictable and makes
API changes reviewable.

### II. Idiomatic Go
Implementation MUST use supported Go language and standard-library conventions, explicit
dependency boundaries, and context propagation for request-scoped work. Errors MUST be
handled explicitly and wrapped with useful context at system boundaries; panics MUST NOT
be used for ordinary control flow. Packages MUST have focused responsibilities, and
abstractions MUST be introduced only where they clarify a real domain or integration
boundary. This favors maintainable code and makes failures diagnosable.

### III. PostgreSQL Integrity
PostgreSQL is the system of record for persistent application data. Database constraints
MUST enforce applicable invariants, and operations that must succeed or fail together
MUST use transactions. Schema changes MUST be delivered as versioned, reviewable
migrations with a safe deployment and rollback or recovery strategy. Queries MUST use
parameter binding rather than interpolating untrusted values. Every database table MUST
include the audit and soft-delete columns `created_at`, `created_by`, `updated_at`,
`updated_by`, `deleted` (boolean), and `deleted_at`. These rules protect financial and
customer data from inconsistent writes, support accountability, and enable consistent
soft deletion.

### IV. Automated Testing
Behavioral changes MUST include automated tests at the appropriate level. Unit tests
MUST cover domain logic and edge cases; API and persistence changes MUST include
integration tests that verify HTTP behavior, SQL behavior, and relevant migrations
against PostgreSQL. Tests MUST be deterministic and runnable through documented project
commands. The tests provide evidence that contracts and persisted invariants remain
correct as the service evolves.

### V. Secure and Observable Operations
All externally supplied data MUST be validated at system boundaries, and access to
endpoints and data MUST be authorized according to the application's established
security policy. Secrets MUST be supplied through approved runtime configuration and
MUST NOT be committed or written to logs. Logs MUST be structured and useful for
diagnosis while excluding credentials and sensitive customer or financial data.
Security-relevant failures and operational errors MUST remain visible to operators.

## Technology and Data Constraints

The service MUST be implemented in Go and expose its application interface as a REST
API. PostgreSQL MUST remain the authoritative persistent database. Database access MUST
preserve transaction boundaries and use least-privilege credentials in deployed
environments. Schema migrations MUST be compatible with the deployment sequence and
MUST apply the required audit and soft-delete columns to every table; any unavoidable
downtime or data transformation risk MUST be documented and approved before release.
Data retention and handling MUST follow applicable project and legal requirements; this
constitution does not substitute for those requirements.

## Development Workflow

Every change MUST be reviewed for API compatibility, data integrity, security, tests,
and operational impact before merge. Continuous integration MUST run formatting,
static analysis, and the relevant automated tests; a change MUST NOT be merged when a
required check fails. Changes to API contracts, database schemas, or runtime
configuration MUST include the corresponding documentation and rollout considerations.
Exceptions to a principle require an explicit rationale, documented risk, and approval
from the designated project maintainers before implementation.

## Governance

This constitution governs implementation and review decisions for the lending backend.
Amendments MUST be proposed as reviewed changes that explain the motivation, affected
principles, compatibility impact, and any required migration or follow-up work. Project
maintainers MUST approve amendments, and each amendment MUST update the version,
ratification date when first adopted, and last-amended date.

Constitution versions use semantic versioning. A MAJOR increment is required when a
principle is removed or changed incompatibly; a MINOR increment is required when a
principle or governance requirement is added or materially expanded; a PATCH increment
is required for clarifications and non-semantic wording changes. Reviews and release
readiness checks MUST verify compliance with these principles. Any approved exception
MUST identify the violated rule, scope, owner, mitigation, and expiration or review
date.

**Version**: 1.2.0 | **Ratified**: TODO(RATIFICATION_DATE): original adoption date unknown | **Last Amended**: 2026-10-07
