package discount

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	Repo Repository
	Now  func() time.Time
}

func (s Service) Validate(ctx context.Context, rawCode string, subtotal int64) ValidationResult {
	code := strings.ToUpper(strings.TrimSpace(rawCode))
	if code == "" {
		return ValidationResult{Valid: true}
	}
	c, err := s.Repo.ByCode(ctx, code)
	if errors.Is(err, pgx.ErrNoRows) {
		return ValidationResult{Message: "Código inválido"}
	}
	if err != nil {
		return ValidationResult{Message: "No se pudo validar el código"}
	}
	return ValidateCode(*c, subtotal, s.now())
}

func ValidateCode(c Code, subtotal int64, now time.Time) ValidationResult {
	if !c.Active {
		return ValidationResult{Code: c.Code, Message: "Código inactivo"}
	}
	if c.ExpiresAt != nil && c.ExpiresAt.Before(now) {
		return ValidationResult{Code: c.Code, Message: "Código vencido"}
	}
	if c.MaxUses != nil && c.Uses >= *c.MaxUses {
		return ValidationResult{Code: c.Code, Message: "Código sin usos disponibles"}
	}
	if subtotal < c.MinOrderCents {
		return ValidationResult{Code: c.Code, Message: "El mínimo de compra no fue alcanzado"}
	}
	discount := int64(0)
	switch c.DiscountType {
	case "percent":
		if c.DiscountValue < 0 || c.DiscountValue > 100 {
			return ValidationResult{Code: c.Code, Message: "Código inválido"}
		}
		discount = subtotal * c.DiscountValue / 100
	case "fixed":
		discount = c.DiscountValue
	default:
		return ValidationResult{Code: c.Code, Message: "Código inválido"}
	}
	if discount > subtotal {
		discount = subtotal
	}
	return ValidationResult{Valid: true, Code: c.Code, DiscountARSCents: discount}
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}
