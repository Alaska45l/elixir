package orders

import (
	"context"
	"errors"
	"strings"

	"elixir/backend/internal/discount"
)

type Service struct {
	Repo     Repository
	Discount discount.Service
}

type Repository interface {
	FetchVariant(ctx context.Context, variantID string) (ValidatedItem, error)
	Create(ctx context.Context, req CreateOrderRequest, validation CartValidationResult, discountCents int64) (*Order, error)
	ByExternalReference(ctx context.Context, ref string) (*Order, error)
}

func (s Service) ValidateCart(ctx context.Context, req CartValidationRequest) (CartValidationResult, error) {
	result := CartValidationResult{Valid: true, Items: []ValidatedItem{}}
	for _, in := range req.Items {
		if in.Quantity <= 0 {
			result.Valid = false
			result.Errors = append(result.Errors, "La cantidad debe ser mayor a cero")
			continue
		}
		item, err := s.Repo.FetchVariant(ctx, in.VariantID)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, "Una fragancia ya no está disponible")
			continue
		}
		item.Quantity = in.Quantity
		item.CorrectedPrice = in.ClientUnitPriceCents > 0 && in.ClientUnitPriceCents != item.UnitPriceARSCents
		if item.AvailableStock < in.Quantity {
			result.Valid = false
			result.Errors = append(result.Errors, item.ProductName+" no tiene stock suficiente")
		}
		item.SubtotalARSCents = item.UnitPriceARSCents * int64(in.Quantity)
		result.SubtotalARSCents += item.SubtotalARSCents
		result.Items = append(result.Items, item)
	}
	return result, nil
}

func (s Service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
	if strings.TrimSpace(req.CustomerName) == "" || strings.TrimSpace(req.CustomerEmail) == "" {
		return nil, errors.New("faltan datos del cliente")
	}
	validation, err := s.ValidateCart(ctx, CartValidationRequest{Items: req.Items})
	if err != nil {
		return nil, err
	}
	if !validation.Valid || len(validation.Items) == 0 {
		return nil, errors.New("carrito inválido")
	}
	discountResult := s.Discount.Validate(ctx, req.DiscountCode, validation.SubtotalARSCents)
	if req.DiscountCode != "" && !discountResult.Valid {
		return nil, errors.New(discountResult.Message)
	}
	return s.Repo.Create(ctx, req, validation, discountResult.DiscountARSCents)
}

func CalculateTotal(subtotal, shipping, discountCents int64) int64 {
	total := subtotal + shipping - discountCents
	if total < 0 {
		return 0
	}
	return total
}
