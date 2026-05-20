package discount

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Pool *pgxpool.Pool
}

func (r Repository) ByCode(ctx context.Context, code string) (*Code, error) {
	if r.Pool == nil {
		return nil, pgx.ErrNoRows
	}
	row := r.Pool.QueryRow(ctx, `SELECT id, code, discount_type, discount_value, min_order_cents, max_uses, uses, active, expires_at, created_at FROM discount_codes WHERE code=$1`, strings.ToUpper(strings.TrimSpace(code)))
	var c Code
	return &c, row.Scan(&c.ID, &c.Code, &c.DiscountType, &c.DiscountValue, &c.MinOrderCents, &c.MaxUses, &c.Uses, &c.Active, &c.ExpiresAt, &c.CreatedAt)
}

func (r Repository) IncrementUse(ctx context.Context, code string) error {
	if r.Pool == nil || strings.TrimSpace(code) == "" {
		return nil
	}
	_, err := r.Pool.Exec(ctx, `UPDATE discount_codes SET uses = uses + 1 WHERE code=$1`, strings.ToUpper(strings.TrimSpace(code)))
	return err
}
