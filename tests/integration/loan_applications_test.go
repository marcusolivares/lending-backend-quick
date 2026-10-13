package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/marcusolivares/lending-backend-quick/internal/postgres"
	"github.com/pressly/goose/v3"
)

func TestLoanApplicationLifecycle(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL not configured")
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, dsn)
	if err != nil {
		t.Skipf("database unavailable: %v", err)
	}
	defer func() {
		if err := postgres.ClosePool(context.Background(), pool); err != nil {
			t.Logf("close pool: %v", err)
		}
	}()

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open database for goose: %v", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("set goose dialect: %v", err)
	}
	if err := goose.Up(db, "../../migrations"); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	if _, err := pool.Exec(ctx, `TRUNCATE TABLE car_collaterals, spouses, customers, loan_applications RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("reset database: %v", err)
	}

	if _, err := pool.Exec(ctx, `INSERT INTO loan_applications (id, amount, tenure_months, status, created_at, created_by, updated_at, updated_by, deleted, deleted_at)
		VALUES ('00000000-0000-0000-0000-000000000001', 150000, 24, 'IN_PROGRESS', NOW(), 'system', NOW(), 'system', false, NULL)`); err != nil {
		t.Fatalf("seed loan application: %v", err)
	}

	var count int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM loan_applications`).Scan(&count); err != nil {
		t.Fatalf("count loan applications: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 loan application, got %d", count)
	}
}
