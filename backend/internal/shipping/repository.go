package shipping

import (
	"context"
	"strings"

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

func (r Repository) OwnFleetQuote(ctx context.Context, req QuoteRequest) ([]QuoteOption, error) {
	zones, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	out := []QuoteOption{}
	for _, zone := range zones {
		if !matchesProvince(zone.ProvinceCodes, req.ProvinceCode) {
			continue
		}
		kg := int64((req.WeightGrams + 999) / 1000)
		if kg < 1 {
			kg = 1
		}
		out = append(out, QuoteOption{
			ID:               "own-fleet-" + strings.ToLower(strings.ReplaceAll(zone.ZoneName, " ", "-")),
			CarrierName:      "ELIXIR",
			ServiceName:      "Envío a domicilio",
			PriceCents:       zone.BaseCostCents + zone.PerKGCents*kg,
			EstimatedDaysMin: zone.EstimatedDaysMin,
			EstimatedDaysMax: zone.EstimatedDaysMax,
		})
	}
	return out, nil
}

func matchesProvince(codes []string, province string) bool {
	if len(codes) == 0 {
		return true
	}
	normalized := provinceAliases(strings.ToUpper(strings.TrimSpace(province)))
	for _, code := range codes {
		for _, alias := range provinceAliases(strings.ToUpper(strings.TrimSpace(code))) {
			for _, p := range normalized {
				if alias == "AR" && p != "C" && p != "CF" && p != "CABA" && p != "B" && p != "BA" && p != "BUENOS_AIRES" {
					return true
				}
				if alias == p {
					return true
				}
			}
		}
	}
	return false
}

func provinceAliases(code string) []string {
	switch code {
	case "C", "CF", "CABA":
		return []string{"C", "CF", "CABA"}
	case "B", "BA", "BUENOS_AIRES":
		return []string{"B", "BA", "BUENOS_AIRES"}
	case "AR":
		return []string{"AR"}
	default:
		return []string{code}
	}
}
