CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS loan_applications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    amount NUMERIC(18,2) NOT NULL CHECK (amount > 0 AND amount % 100 = 0),
    tenure_months INTEGER NOT NULL CHECK (tenure_months BETWEEN 3 AND 48),
    status VARCHAR(32) NOT NULL DEFAULT 'IN_PROGRESS' CHECK (status IN ('IN_PROGRESS', 'APPROVED', 'REJECTED', 'CANCELLED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL DEFAULT 'system',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL DEFAULT 'system',
    deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS customers (
    application_id UUID PRIMARY KEY REFERENCES loan_applications(id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    id_number VARCHAR(255) NOT NULL,
    address_street VARCHAR(255) NOT NULL,
    address_city VARCHAR(255) NOT NULL,
    address_postal_code VARCHAR(64) NOT NULL,
    birth_date DATE NOT NULL,
    phone_number VARCHAR(64) NOT NULL,
    email VARCHAR(255) NOT NULL,
    monthly_income NUMERIC(18,2) NOT NULL CHECK (monthly_income >= 0),
    marital_status VARCHAR(16) NOT NULL CHECK (marital_status IN ('MARRIED', 'UNMARRIED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL DEFAULT 'system',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL DEFAULT 'system',
    deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS spouses (
    application_id UUID PRIMARY KEY REFERENCES loan_applications(id) ON DELETE CASCADE,
    full_name VARCHAR(255) NOT NULL,
    id_card_number VARCHAR(255) NOT NULL,
    phone_number VARCHAR(64) NOT NULL,
    monthly_income NUMERIC(18,2) NOT NULL CHECK (monthly_income >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL DEFAULT 'system',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL DEFAULT 'system',
    deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ NULL
);

CREATE TABLE IF NOT EXISTS car_collaterals (
    application_id UUID PRIMARY KEY REFERENCES loan_applications(id) ON DELETE CASCADE,
    brand VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    manufacturing_year INTEGER NOT NULL CHECK (manufacturing_year BETWEEN 1000 AND 9999),
    license_plate VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by VARCHAR(255) NOT NULL DEFAULT 'system',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_by VARCHAR(255) NOT NULL DEFAULT 'system',
    deleted BOOLEAN NOT NULL DEFAULT false,
    deleted_at TIMESTAMPTZ NULL
);
