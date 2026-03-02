package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/finset/app/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

func Connect() (*Pool, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse DATABASE_URL: %w", err)
	}

	cfg.MaxConns = 10
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.HealthCheckPeriod = 1 * time.Minute

	var pool *pgxpool.Pool
	for attempt := 1; attempt <= 10; attempt++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
		if err != nil {
			log.Printf("DB connect attempt %d/10 failed: %v — retrying in 3s", attempt, err)
			time.Sleep(3 * time.Second)
			continue
		}
		if err = pool.Ping(context.Background()); err != nil {
			pool.Close()
			log.Printf("DB ping attempt %d/10 failed: %v — retrying in 3s", attempt, err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}
	if err != nil {
		return nil, fmt.Errorf("could not connect to database after 10 attempts: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return &Pool{pool}, nil
}

func (p *Pool) Migrate() error {
	ctx := context.Background()

	_, err := p.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS transactions (
			id          TEXT          PRIMARY KEY,
			type        TEXT          NOT NULL CHECK (type IN ('income','expense')),
			amount      NUMERIC(12,2) NOT NULL CHECK (amount > 0),
			category    TEXT          NOT NULL DEFAULT '',
			method      TEXT          NOT NULL DEFAULT 'Cash',
			date        DATE          NOT NULL,
			note        TEXT          NOT NULL DEFAULT '',
			created_at  TIMESTAMPTZ   NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create transactions table: %w", err)
	}

	addCols := []string{
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS category   TEXT        NOT NULL DEFAULT ''`,
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS method     TEXT        NOT NULL DEFAULT 'Cash'`,
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS note       TEXT        NOT NULL DEFAULT ''`,
		`ALTER TABLE transactions ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`,
	}
	for _, sql := range addCols {
		if _, err := p.Exec(ctx, sql); err != nil {
			return fmt.Errorf("alter table: %w", err)
		}
	}

	if _, err = p.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_transactions_date ON transactions (date DESC)`); err != nil {
		return fmt.Errorf("create index date: %w", err)
	}
	if _, err = p.Exec(ctx, `CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions (type)`); err != nil {
		return fmt.Errorf("create index type: %w", err)
	}

	log.Println("Database migration complete")
	return nil
}

func (p *Pool) ListTransactions(ctx context.Context) ([]models.Transaction, error) {
	rows, err := p.Query(ctx, `
		SELECT id, type, amount, category, method, TO_CHAR(date,'YYYY-MM-DD'), note, created_at
		FROM transactions
		ORDER BY date DESC, created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.Type, &t.Amount, &t.Category, &t.Method, &t.Date, &t.Note, &t.CreatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	if txs == nil {
		txs = []models.Transaction{}
	}
	return txs, rows.Err()
}

func (p *Pool) GetTransaction(ctx context.Context, id string) (*models.Transaction, error) {
	var t models.Transaction
	err := p.QueryRow(ctx, `
		SELECT id, type, amount, category, method, TO_CHAR(date,'YYYY-MM-DD'), note, created_at
		FROM transactions WHERE id = $1
	`, id).Scan(&t.ID, &t.Type, &t.Amount, &t.Category, &t.Method, &t.Date, &t.Note, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (p *Pool) CreateTransaction(ctx context.Context, id string, r models.CreateTransactionRequest) (*models.Transaction, error) {
	var t models.Transaction
	err := p.QueryRow(ctx, `
		INSERT INTO transactions (id, type, amount, category, method, date, note)
		VALUES ($1, $2, $3, $4, $5, $6::DATE, $7)
		RETURNING id, type, amount, category, method, TO_CHAR(date,'YYYY-MM-DD'), note, created_at
	`, id, r.Type, r.Amount, r.Category, r.Method, r.Date, r.Note).
		Scan(&t.ID, &t.Type, &t.Amount, &t.Category, &t.Method, &t.Date, &t.Note, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (p *Pool) DeleteTransaction(ctx context.Context, id string) (bool, error) {
	tag, err := p.Exec(ctx, `DELETE FROM transactions WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (p *Pool) BulkInsert(ctx context.Context, txs []models.Transaction) (int, error) {
	tx, err := p.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var inserted int
	for _, t := range txs {
		tag, err := tx.Exec(ctx, `
			INSERT INTO transactions (id, type, amount, category, method, date, note, created_at)
			VALUES ($1, $2, $3, $4, $5, $6::DATE, $7, $8)
			ON CONFLICT (id) DO NOTHING
		`, t.ID, t.Type, t.Amount, t.Category, t.Method, t.Date, t.Note, t.CreatedAt)
		if err != nil {
			return 0, err
		}
		inserted += int(tag.RowsAffected())
	}
	return inserted, tx.Commit(ctx)
}

type Stats struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
	Count        int     `json:"count"`
}

func (p *Pool) GetStats(ctx context.Context) (Stats, error) {
	var s Stats
	err := p.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN type='income'  THEN amount ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END), 0),
			COUNT(*)
		FROM transactions
	`).Scan(&s.TotalIncome, &s.TotalExpense, &s.Count)
	s.Balance = s.TotalIncome - s.TotalExpense
	return s, err
}

type MonthlyRow struct {
	Month   string  `json:"month"`
	Year    int     `json:"year"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

// GetMonthlyFlow returns per-month totals for the last N months.
// Uses fmt.Sprintf to inline the months value — avoids pgx parameter
// type casting issues with INTERVAL arithmetic.
func (p *Pool) GetMonthlyFlow(ctx context.Context, months int) ([]MonthlyRow, error) {
	if months < 1 {
		months = 7
	}
	if months > 24 {
		months = 24
	}

	// Build the cutoff date in Go — no $1 in INTERVAL expression
	cutoff := time.Now().AddDate(0, -(months - 1), 0)
	cutoffStr := fmt.Sprintf("%d-%02d-01", cutoff.Year(), cutoff.Month())

	rows, err := p.Query(ctx, `
		SELECT
			TO_CHAR(DATE_TRUNC('month', date), 'Mon')  AS month,
			EXTRACT(YEAR FROM date)::INT               AS year,
			COALESCE(SUM(CASE WHEN type='income'  THEN amount ELSE 0 END), 0) AS income,
			COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END), 0) AS expense
		FROM transactions
		WHERE date >= $1::DATE
		GROUP BY DATE_TRUNC('month', date)
		ORDER BY DATE_TRUNC('month', date) ASC
	`, cutoffStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []MonthlyRow
	for rows.Next() {
		var r MonthlyRow
		if err := rows.Scan(&r.Month, &r.Year, &r.Income, &r.Expense); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	if result == nil {
		result = []MonthlyRow{}
	}
	return result, rows.Err()
}

type CategoryRow struct {
	Category string  `json:"category"`
	Total    float64 `json:"total"`
}

func (p *Pool) GetCategoryBreakdown(ctx context.Context) ([]CategoryRow, error) {
	rows, err := p.Query(ctx, `
		SELECT category, SUM(amount) AS total
		FROM transactions
		WHERE type = 'expense'
		GROUP BY category
		ORDER BY total DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []CategoryRow
	for rows.Next() {
		var r CategoryRow
		if err := rows.Scan(&r.Category, &r.Total); err != nil {
			return nil, err
		}
		result = append(result, r)
	}
	if result == nil {
		result = []CategoryRow{}
	}
	return result, rows.Err()
}
