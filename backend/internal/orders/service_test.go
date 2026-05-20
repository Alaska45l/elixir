package orders

import (
	"context"
	"errors"
	"testing"
)

type fakeOrderRepo struct {
	item ValidatedItem
}

func (f fakeOrderRepo) FetchVariant(context.Context, string) (ValidatedItem, error) {
	return f.item, nil
}
func (f fakeOrderRepo) Create(context.Context, CreateOrderRequest, CartValidationResult, int64) (*Order, error) {
	return &Order{}, nil
}
func (f fakeOrderRepo) ByExternalReference(context.Context, string) (*Order, error) {
	return nil, errors.New("not implemented")
}

func TestValidateCartCorrectsMismatchedPrice(t *testing.T) {
	s := Service{Repo: fakeOrderRepo{item: ValidatedItem{ProductName: "Nocturno Oud", UnitPriceARSCents: 10000, AvailableStock: 10}}}
	res, err := s.ValidateCart(context.Background(), CartValidationRequest{Items: []CartItemRequest{{VariantID: "v1", Quantity: 1, ClientUnitPriceCents: 9000}}})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Items[0].CorrectedPrice || res.SubtotalARSCents != 10000 {
		t.Fatalf("expected corrected server price, got %+v", res)
	}
}

func TestValidateCartRejectsOutOfStock(t *testing.T) {
	s := Service{Repo: fakeOrderRepo{item: ValidatedItem{ProductName: "Nocturno Oud", UnitPriceARSCents: 10000, AvailableStock: 1}}}
	res, _ := s.ValidateCart(context.Background(), CartValidationRequest{Items: []CartItemRequest{{VariantID: "v1", Quantity: 2}}})
	if res.Valid {
		t.Fatal("expected invalid cart for out-of-stock quantity")
	}
}

func TestValidateCartRejectsZeroQuantity(t *testing.T) {
	s := Service{Repo: fakeOrderRepo{}}
	res, _ := s.ValidateCart(context.Background(), CartValidationRequest{Items: []CartItemRequest{{VariantID: "v1", Quantity: 0}}})
	if res.Valid {
		t.Fatal("expected invalid cart for zero quantity")
	}
}

func TestCalculateTotal(t *testing.T) {
	if got := CalculateTotal(10000, 2500, 1500); got != 11000 {
		t.Fatalf("expected 11000, got %d", got)
	}
}
