# Data Model: Loan Submission

## Loan Application

Represents one accepted loan proposal.

| Field | Type / constraint | Notes |
|---|---|---|
| `id` | UUID primary key | Returned to identify the submitted application. |
| `amount` | Exact decimal, greater than zero, divisible by 100 | In the application's configured currency. |
| `tenure_months` | Integer, 3–48 inclusive | Whole months. |
| `status` | Application status | Set to `IN_PROGRESS` at creation. |
| `created_at`, `created_by` | Audit fields | Required by project data governance. |
| `updated_at`, `updated_by` | Audit fields | Required by project data governance. |
| `deleted`, `deleted_at` | Soft-delete fields | `deleted` is boolean. |

## Customer

One applicant record belongs to exactly one loan application in this feature.

| Field | Type / constraint | Notes |
|---|---|---|
| `application_id` | UUID primary key and foreign key | Enforces one customer per application. |
| `full_name` | Required non-empty text | Applicant's full name. |
| `id_number` | Required non-empty text | No country-specific format or uniqueness rule is specified. |
| `address_street` | Required non-empty text | Street portion of the address. |
| `address_city` | Required non-empty text | City portion of the address. |
| `address_postal_code` | Required non-empty text | ZIP/postal code. |
| `birth_date` | Calendar date | Must be a real date in `YYYY-MM-DD` request format. No age threshold is specified. |
| `phone_number` | Required non-empty text | Detailed format follows the supported customer locale. |
| `email` | Required text with syntactically valid email | Ownership/deliverability verification is out of scope. |
| `monthly_income` | Exact decimal, non-negative | In the application's configured currency. |
| `marital_status` | Required status distinguishing married from unmarried | Spouse presence and required fields depend on this value. |
| `created_at`, `created_by`, `updated_at`, `updated_by`, `deleted`, `deleted_at` | Required audit and soft-delete fields | As required for every persisted table. |

## Spouse

An optional one-to-one record that exists only for a married customer. It is stored
within the same transaction as the application.

| Field | Type / constraint | Notes |
|---|---|---|
| `application_id` | UUID primary key and foreign key | Enforces at most one spouse per application. |
| `full_name` | Required non-empty text | Spouse's full name. |
| `id_card_number` | Required non-empty text | No country-specific format or uniqueness rule is specified. |
| `phone_number` | Required non-empty text | Detailed format follows the supported customer locale. |
| `monthly_income` | Exact decimal, non-negative | In the application's configured currency. |
| Audit/soft-delete fields | Same fields as other persisted tables | Required by project data governance. |

## Car Collateral

One collateral record belongs to exactly one application.

| Field | Type / constraint | Notes |
|---|---|---|
| `application_id` | UUID primary key and foreign key | Enforces one car collateral record per application. |
| `brand` | Required non-empty text | No vehicle-catalog lookup. |
| `model` | Required non-empty text | No vehicle-catalog lookup. |
| `manufacturing_year` | Four-digit year, no later than current calendar year | Request validation enforces the current-year comparison. |
| `license_plate` | Required non-empty text | No country-specific format validation. |
| Audit/soft-delete fields | Same fields as other persisted tables | Required by project data governance. |

## Relationships and Invariants

- A loan application has exactly one customer and one car collateral record.
- A loan application has a spouse record if and only if the customer is married.
- Unmarried submissions omit spouse data; married submissions require every spouse
  field.
- All records for one submission are created atomically. A failed validation or database
  operation leaves no partial application records.
- All four tables carry `created_at`, `created_by`, `updated_at`, `updated_by`,
  `deleted`, and `deleted_at`.
- No uniqueness/deduplication rule is specified for customer ID numbers or license
  plates; do not add one without a product requirement.

## State

New applications are persisted with `IN_PROGRESS`. Later status transitions are outside
this feature.
