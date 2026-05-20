package shipping

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Zone struct {
	ID               string   `json:"id"`
	ZoneName         string   `json:"zone_name"`
	ProvinceCodes    []string `json:"province_codes"`
	BaseCostCents    int64    `json:"base_cost_cents"`
	PerKGCents       int64    `json:"per_kg_cents"`
	EstimatedDaysMin int      `json:"estimated_days_min"`
	EstimatedDaysMax int      `json:"estimated_days_max"`
	Active           bool     `json:"active"`
}

type Repository struct {
	Pool *pgxpool.Pool
}

func (r Repository) List(ctx context.Context) ([]Zone, error) {
	if r.Pool == nil {
		return []Zone{}, nil
	}
	rows, err := r.Pool.Query(ctx, `SELECT id, zone_name, COALESCE(province_codes,'{}'), base_cost_cents, per_kg_cents, estimated_days_min, estimated_days_max, active FROM shipping_zones WHERE active=true ORDER BY base_cost_cents ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Zone{}
	for rows.Next() {
		var z Zone
		if err := rows.Scan(&z.ID, &z.ZoneName, &z.ProvinceCodes, &z.BaseCostCents, &z.PerKGCents, &z.EstimatedDaysMin, &z.EstimatedDaysMax, &z.Active); err != nil {
			return nil, err
		}
		out = append(out, z)
	}
	return out, rows.Err()
}
