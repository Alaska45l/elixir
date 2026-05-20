package discount

import "time"

type Code struct {
	ID            string     `json:"id"`
	Code          string     `json:"code"`
	DiscountType  string     `json:"discount_type"`
	DiscountValue int64      `json:"discount_value"`
	MinOrderCents int64      `json:"min_order_cents"`
	MaxUses       *int       `json:"max_uses,omitempty"`
	Uses          int        `json:"uses"`
	Active        bool       `json:"active"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type ValidationResult struct {
	Valid            bool   `json:"valid"`
	Code             string `json:"code,omitempty"`
	DiscountARSCents int64  `json:"discount_ars_cents"`
	Message          string `json:"message,omitempty"`
}
