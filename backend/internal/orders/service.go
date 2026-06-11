package orders

import (
	"context"
	"errors"
	"net/mail"
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
	if len(req.Items) == 0 {
		return CartValidationResult{Valid: false, Items: []ValidatedItem{}, Errors: []string{"El carrito está vacío"}}, nil
	}
	if len(req.Items) > 50 {
		return CartValidationResult{Valid: false, Items: []ValidatedItem{}, Errors: []string{"El carrito tiene demasiados productos"}}, nil
	}
	quantities := map[string]int{}
	order := []CartItemRequest{}
	for _, in := range req.Items {
		in.VariantID = strings.TrimSpace(in.VariantID)
		if in.VariantID == "" || in.Quantity <= 0 {
			result.Valid = false
			result.Errors = append(result.Errors, "La cantidad debe ser mayor a cero")
			continue
		}
		if in.Quantity > 99 {
			result.Valid = false
			result.Errors = append(result.Errors, "La cantidad solicitada es demasiado alta")
			continue
		}
		if _, exists := quantities[in.VariantID]; !exists {
			order = append(order, in)
		}
		quantities[in.VariantID] += in.Quantity
	}
	for _, in := range order {
		in.Quantity = quantities[in.VariantID]
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
	if err := normalizeCreateOrderRequest(&req); err != nil {
		return nil, err
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
	req.DiscountCode = discountResult.Code
	return s.Repo.Create(ctx, req, validation, discountResult.DiscountARSCents)
}

func CalculateTotal(subtotal, shipping, discountCents int64) int64 {
	total := subtotal + shipping - discountCents
	if total < 0 {
		return 0
	}
	return total
}

func normalizeCreateOrderRequest(req *CreateOrderRequest) error {
	req.CustomerName = strings.TrimSpace(req.CustomerName)
	req.CustomerEmail = strings.ToLower(strings.TrimSpace(req.CustomerEmail))
	req.CustomerPhone = strings.TrimSpace(req.CustomerPhone)
	req.DiscountCode = strings.ToUpper(strings.TrimSpace(req.DiscountCode))
	req.Notes = strings.TrimSpace(req.Notes)
	if len(req.CustomerName) < 2 || len(req.CustomerName) > 120 {
		return errors.New("faltan datos del cliente")
	}
	if len(req.CustomerEmail) > 254 {
		return errors.New("email inválido")
	}
	if _, err := mail.ParseAddress(req.CustomerEmail); err != nil {
		return errors.New("email inválido")
	}
	if len(req.CustomerPhone) > 40 {
		return errors.New("teléfono inválido")
	}
	if req.ShippingCostARSCents < 0 || req.ShippingCostARSCents > 100_000_000 {
		return errors.New("costo de envío inválido")
	}
	if len(req.Notes) > 1000 {
		return errors.New("las notas no pueden superar 1000 caracteres")
	}
	if len(req.ShippingAddress) == 0 {
		return errors.New("dirección de envío requerida")
	}
	return nil
}
